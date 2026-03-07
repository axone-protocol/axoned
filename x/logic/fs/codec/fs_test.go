package codec

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

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fsiface "github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
)

func TestCodecDeviceFSOpenFileValidation(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("codec fs should implement OpenFileFS")
	}

	Convey("Given a codec device filesystem", t, func() {
		Convey("when opening with unsupported mode", func() {
			_, err := ofs.OpenFile(devicePath, os.O_RDONLY, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})

		Convey("when opening an unknown path", func() {
			_, err := ofs.OpenFile("unknown", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})
}

func TestCodecDeviceFSTransactionLifecycle(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs := vfs.(fsiface.OpenFileFS)

	Convey("Given an opened bech32 codec device", t, func() {
		file, err := ofs.OpenFile(devicePath, os.O_RDWR, 0)
		So(err, ShouldBeNil)

		rw, ok := file.(interface {
			Read([]byte) (int, error)
			Write([]byte) (int, error)
			Close() error
		})
		So(ok, ShouldBeTrue)

		Convey("when writing a decode request and reading the response", func() {
			request := []byte("  decode   axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4  \r\n")

			n, err := rw.Write(request[:8])
			So(err, ShouldBeNil)
			So(n, ShouldEqual, 8)

			n, err = rw.Write(request[8:])
			So(err, ShouldBeNil)
			So(n, ShouldEqual, len(request)-8)

			got := bytes.NewBuffer(nil)
			buf := make([]byte, 16)
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

			So(got.String(), ShouldEqual, "ok(-(axone,[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30])).\n")
		})

		Convey("when writing an encode request and reading the response", func() {
			request := []byte("encode axone a3a717f4a2af31a2aa0fb58d44868da81238f71e\n")

			_, err := rw.Write(request)
			So(err, ShouldBeNil)

			got := bytes.NewBuffer(nil)
			buf := make([]byte, 16)
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

			So(got.String(), ShouldEqual, "ok(axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4).\n")
		})
	})
}

func TestCodecDeviceFSResponses(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs := vfs.(fsiface.OpenFileFS)

	readAll := func(request []byte) ([]byte, error) {
		file, err := ofs.OpenFile(devicePath, os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}

		if len(request) > 0 {
			if _, err := file.(interface{ Write([]byte) (int, error) }).Write(request); err != nil {
				return nil, err
			}
		}

		defer file.Close()
		return io.ReadAll(file)
	}

	Convey("Given a codec device filesystem", t, func() {
		Convey("when no request bytes are written before the first read", func() {
			response, err := readAll(nil)
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_request).\n")
		})

		Convey("when the request opcode is unknown", func() {
			response, err := readAll([]byte("unknown"))
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_request).\n")
		})

		Convey("when the request contains malformed UTF-8", func() {
			response, err := readAll([]byte{0xff})
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_request).\n")
		})

		Convey("when the bech32 payload is invalid", func() {
			response, err := readAll([]byte("decode bad"))
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_bech32).\n")
		})

		Convey("when the encode request is missing the hex payload", func() {
			response, err := readAll([]byte("encode axone"))
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_request).\n")
		})

		Convey("when the request uses tabs as separators", func() {
			response, err := readAll([]byte("encode\taxone\t00"))
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_request).\n")
		})

		Convey("when the encode bytes are not valid hex", func() {
			response, err := readAll([]byte("encode axone 0011zz"))
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_bytes).\n")
		})

		Convey("when the encode bytes have odd-length hex", func() {
			response, err := readAll([]byte("encode axone 123"))
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_bytes).\n")
		})

		Convey("when the encode hrp exceeds the protocol limit", func() {
			request := append([]byte("encode "), bytes.Repeat([]byte{'a'}, 256)...)
			request = append(request, []byte(" 00")...)

			response, err := readAll(request)
			So(err, ShouldBeNil)
			So(string(response), ShouldEqual, "error(invalid_hrp).\n")
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
