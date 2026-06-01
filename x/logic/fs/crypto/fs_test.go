package cryptofs

import (
	"bytes"
	"crypto/sha256"
	"crypto/sha512"
	"errors"
	"io"
	"io/fs"
	"os"
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

	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/devfile"
	fsiface "github.com/axone-protocol/axoned/v15/x/logic/fs/internal/iface"
)

func TestCryptoDeviceFSOpen(t *testing.T) {
	vfs := NewFS(newSDKContext())

	Convey("Given a crypto device filesystem", t, func() {
		_, err := vfs.Open("hash/sha256")
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
	})
}

func TestCryptoDeviceFSOpenFileValidation(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("crypto fs should implement OpenFileFS")
	}

	Convey("Given a crypto device filesystem", t, func() {
		Convey("when opening with unsupported mode", func() {
			_, err := ofs.OpenFile("hash/sha256", os.O_RDONLY, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})

		Convey("when opening an unknown path", func() {
			_, err := ofs.OpenFile("hash/unknown", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when opening an escaping path", func() {
			_, err := ofs.OpenFile("../hash/sha256", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func TestCryptoDeviceFSFunctional(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("crypto fs should implement OpenFileFS")
	}

	readAll := func(path string, request []byte) ([]byte, error) {
		t.Helper()

		file, err := ofs.OpenFile(path, os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if len(request) > 0 {
			writer, ok := file.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("crypto file should implement Write")
			}

			if _, err := writer.Write(request); err != nil {
				return nil, err
			}
		}

		return io.ReadAll(file)
	}

	Convey("Given a crypto device filesystem", t, func() {
		request := []byte("hello world")

		Convey("when hashing with sha256", func() {
			expected := sha256.Sum256(request)
			got, err := readAll("hash/sha256", request)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected[:])
		})

		Convey("when hashing with sha512", func() {
			expected := sha512.Sum512(request)
			got, err := readAll("hash/sha512", request)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected[:])
		})

		Convey("when hashing with md5", func() {
			expected := []byte{94, 182, 59, 187, 224, 30, 238, 208, 147, 203, 34, 187, 143, 90, 205, 195}
			got, err := readAll("hash/md5", request)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected)
		})

		Convey("when hashing an empty request", func() {
			expected := sha256.Sum256(nil)
			got, err := readAll("hash/sha256", nil)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected[:])
		})
	})
}

func TestCryptoDeviceFSErrors(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("crypto fs should implement OpenFileFS")
	}

	Convey("Given a crypto device filesystem", t, func() {
		Convey("when request payload exceeds limit", func() {
			file, err := ofs.OpenFile("hash/sha256", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			huge := bytes.Repeat([]byte("x"), maxRequestBytes+1)
			_, err = file.(interface{ Write([]byte) (int, error) }).Write(huge)
			So(errors.Is(err, devfile.ErrWriteLimit), ShouldBeTrue)
		})
	})
}

func newSDKContext() sdk.Context {
	stateStore := store.NewCommitMultiStore(dbm.NewMemDB(), log.NewNopLogger(), metrics.NewNoOpMetrics())
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return ctx.WithHeaderInfo(coreheader.Info{
		Height: 42,
		Time:   time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC),
	})
}
