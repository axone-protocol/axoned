package wasm

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"sync"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/devfile"
	fsiface "github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v14/x/logic/testutil"
)

const testContractAddress = "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk"

func TestWasmDeviceFSOpenFileValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := newSDKContext()
	keeper := testutil.NewMockWasmKeeper(ctrl)
	vfs := NewFS(ctx, keeper)
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("wasm fs should implement OpenFileFS")
	}

	Convey("Given a wasm device filesystem", t, func() {
		Convey("when opening with unsupported mode", func() {
			_, err := ofs.OpenFile(testContractAddress+"/query", os.O_RDONLY, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})

		Convey("when opening an unknown path", func() {
			_, err := ofs.OpenFile(testContractAddress+"/info", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when opening with invalid contract address", func() {
			_, err := ofs.OpenFile("not-bech32/query", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})
}

func TestWasmDeviceFSTransactionLifecycle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := newSDKContext()
	contractAddr := sdk.MustAccAddressFromBech32(testContractAddress)
	request := []byte(`{"ping":"pong"}`)
	response := []byte(`{"ok":true}`)

	keeper := testutil.NewMockWasmKeeper(ctrl)
	keeper.EXPECT().
		QuerySmart(ctx, contractAddr, request).
		Return(response, nil).
		Times(1)

	vfs := NewFS(ctx, keeper)
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("wasm fs should implement OpenFileFS")
	}

	Convey("Given an opened wasm query device", t, func() {
		file, err := ofs.OpenFile(testContractAddress+"/query", os.O_RDWR, 0)
		So(err, ShouldBeNil)

		rw, ok := file.(interface {
			Read([]byte) (int, error)
			Write([]byte) (int, error)
			Close() error
		})
		So(ok, ShouldBeTrue)

		Convey("when writing request bytes and reading response", func() {
			n, err := rw.Write(request[:6])
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 6)

			n, err = rw.Write(request[6:])
			So(err, ShouldBeNil)
			So(n, ShouldEqual, len(request)-6)

			got := bytes.NewBuffer(nil)
			buf := make([]byte, 4)
			for {
				readN, readErr := rw.Read(buf)
				if readN > 0 {
					_, _ = got.Write(buf[:readN])
				}
				if errors.Is(readErr, io.EOF) {
					break
				}
				So(readErr, ShouldBeNil)
			}
			So(got.Bytes(), ShouldResemble, response)

			_, err = rw.Write([]byte("!"))
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func TestWasmDeviceFSErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := newSDKContext()
	contractAddr := sdk.MustAccAddressFromBech32(testContractAddress)

	Convey("Given a wasm device filesystem", t, func() {
		Convey("when first read happens without request payload", func() {
			keeper := testutil.NewMockWasmKeeper(ctrl)
			vfs := NewFS(ctx, keeper)
			ofs := vfs.(fsiface.OpenFileFS)

			file, err := ofs.OpenFile(testContractAddress+"/query", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			buf := make([]byte, 8)
			_, err = file.Read(buf)
			So(errors.Is(err, devfile.ErrInvalidRequest), ShouldBeTrue)
		})

		Convey("when request payload exceeds limit", func() {
			keeper := testutil.NewMockWasmKeeper(ctrl)
			vfs := NewFS(ctx, keeper)
			ofs := vfs.(fsiface.OpenFileFS)

			file, err := ofs.OpenFile(testContractAddress+"/query", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			huge := make([]byte, maxRequestBytes+1)
			_, err = file.(interface{ Write([]byte) (int, error) }).Write(huge)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})

		Convey("when wasm query fails at commit", func() {
			keeper := testutil.NewMockWasmKeeper(ctrl)
			keeper.EXPECT().
				QuerySmart(ctx, contractAddr, []byte(`{"ping":"pong"}`)).
				Return(nil, errors.New("boom")).
				Times(1)

			vfs := NewFS(ctx, keeper)
			ofs := vfs.(fsiface.OpenFileFS)
			file, err := ofs.OpenFile(testContractAddress+"/query", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			_, err = file.(interface{ Write([]byte) (int, error) }).Write([]byte(`{"ping":"pong"}`))
			So(err, ShouldBeNil)

			_, err = file.Read(make([]byte, 8))
			So(errors.Is(err, errWasmQueryFailed), ShouldBeTrue)

			var pathErr *fs.PathError
			So(errors.As(err, &pathErr), ShouldBeTrue)
			So(pathErr.Path, ShouldEqual, testContractAddress+"/query")
		})

		Convey("when wasm response exceeds the configured maximum size", func() {
			keeper := testutil.NewMockWasmKeeper(ctrl)
			keeper.EXPECT().
				QuerySmart(ctx, contractAddr, []byte(`{"ping":"pong"}`)).
				Return(make([]byte, maxResponseBytes+1), nil).
				Times(1)

			vfs := NewFS(ctx, keeper)
			ofs := vfs.(fsiface.OpenFileFS)
			file, err := ofs.OpenFile(testContractAddress+"/query", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			_, err = file.(interface{ Write([]byte) (int, error) }).Write([]byte(`{"ping":"pong"}`))
			So(err, ShouldBeNil)

			_, err = file.Read(make([]byte, 8))
			So(errors.Is(err, devfile.ErrResponseTooLarge), ShouldBeTrue)
		})

		Convey("when closed before first read", func() {
			keeper := testutil.NewMockWasmKeeper(ctrl)
			vfs := NewFS(ctx, keeper)
			ofs := vfs.(fsiface.OpenFileFS)
			file, err := ofs.OpenFile(testContractAddress+"/query", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			_, err = file.(interface{ Write([]byte) (int, error) }).Write([]byte(`{"ping":"pong"}`))
			So(err, ShouldBeNil)

			So(file.Close(), ShouldBeNil)
		})
	})
}

func newSDKContext() sdk.Context {
	bech32ConfigOnce.Do(func() {
		sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")
	})
	stateStore := store.NewCommitMultiStore(dbm.NewMemDB(), log.NewNopLogger(), metrics.NewNoOpMetrics())
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return ctx.WithHeaderInfo(coreheader.Info{
		Height: 42,
		Time:   time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC),
	})
}

var bech32ConfigOnce sync.Once
