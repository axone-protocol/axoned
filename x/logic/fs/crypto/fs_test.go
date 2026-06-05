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
	"github.com/axone-protocol/axoned/v15/x/logic/util"
)

const (
	ed25519PubKeyHex = "53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5a"
	ed25519SigHex    = "889bcfd331e8e43b5ebf430301dffb6ac9e2fce69f6227b43552fe3dc8cc1ee0" +
		"0c1cc53452a8712e9d5f80086dff8cf4999c1b93ed6c6e403c09334cb61ddd0b"

	secp256r1PubKeyHex = "0213c8426be471e55506f7ce4f7df557a42e310df09f92eb732ca3085e797cef9b"
	secp256r1SigHex    = "30450220099e6f9dd218e0e304efa7a4224b0058a8e3aec73367ec239bee4ed8e" +
		"d7d85db022100b504d3d0d2e879b04705c0e5a2b40b0521a5ab647ea207bd81134e1a4eb79e47"

	secp256k1PubKeyHex = "026b5450187ee9c63ba9e42cb6018d8469c903aca116178e223de76e49fe63b71c"
	secp256k1SigHex    = "304402201448201bb4408549b0997f4b9ad9ed36f3cf8bb9c433fc7f3ba48c6b6e" +
		"39476e022053f7d056f7ffeab9a79f3a36bc2ba969ddd530a3a1495d1ed7bba00039820223"

	expectedOKTrue         = "ok(true).\n"
	expectedInvalidRequest = "error(invalid_request).\n"
)

func TestCryptoDeviceFSOpen(t *testing.T) {
	vfs := NewFS(newSDKContext())

	Convey("Given a crypto device filesystem", t, func() {
		_, err := vfs.Open("sha256")
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
			_, err := ofs.OpenFile("sha256", os.O_RDONLY, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})

		Convey("when opening an unknown path", func() {
			_, err := ofs.OpenFile("unknown", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when opening an escaping path", func() {
			_, err := ofs.OpenFile("../sha256", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})

		Convey("when opening a path with extra segments", func() {
			_, err := ofs.OpenFile("sha256/extra", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})
}

func TestCryptoDeviceFSFunctional(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("crypto fs should implement OpenFileFS")
	}

	Convey("Given a crypto device filesystem", t, func() {
		request := []byte("hello world")

		Convey("when hashing with sha256", func() {
			expected := sha256.Sum256(request)
			got, err := readAllFromCryptoDevice(t, ofs, "sha256", request)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected[:])
		})

		Convey("when hashing with sha512", func() {
			expected := sha512.Sum512(request)
			got, err := readAllFromCryptoDevice(t, ofs, "sha512", request)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected[:])
		})

		Convey("when hashing with md5", func() {
			expected := []byte{94, 182, 59, 187, 224, 30, 238, 208, 147, 203, 34, 187, 143, 90, 205, 195}
			got, err := readAllFromCryptoDevice(t, ofs, "md5", request)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected)
		})

		Convey("when hashing an empty request", func() {
			expected := sha256.Sum256(nil)
			got, err := readAllFromCryptoDevice(t, ofs, "sha256", nil)
			So(err, ShouldBeNil)
			So(got, ShouldResemble, expected[:])
		})
	})
}

func TestCryptoDeviceFSSignatureResponses(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("crypto fs should implement OpenFileFS")
	}

	cases := []struct {
		name     string
		path     string
		request  []byte
		expected string
	}{
		{
			name: "valid ed25519 signature",
			path: keyAlgEd25519,
			request: signatureRequest(
				ed25519PubKeyHex,
				"9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d",
				ed25519SigHex,
			),
			expected: expectedOKTrue,
		},
		{
			name: "invalid ed25519 signature",
			path: keyAlgEd25519,
			request: signatureRequest(
				ed25519PubKeyHex,
				"9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9e",
				ed25519SigHex,
			),
			expected: "ok(false).\n",
		},
		{
			name: "valid secp256r1 signature",
			path: "secp256r1",
			request: signatureRequest(
				secp256r1PubKeyHex,
				"e50c26e89f734b2ee12041ff27874c901891f74a0f0cf470333312a3034ce3be",
				secp256r1SigHex,
			),
			expected: expectedOKTrue,
		},
		{
			name: "valid secp256k1 signature",
			path: "secp256k1",
			request: signatureRequest(
				secp256k1PubKeyHex,
				"dece063885d3648078f903b6a3e8989f649dc3368cd9c8d69755ed9dcb6a0995",
				secp256k1SigHex,
			),
			expected: expectedOKTrue,
		},
		{
			name:     "malformed signature request",
			path:     keyAlgEd25519,
			request:  []byte("verify bad\n"),
			expected: expectedInvalidRequest,
		},
		{
			name:     "unsupported signature operation",
			path:     keyAlgEd25519,
			request:  []byte("sign 00 00 00\n"),
			expected: "error(unsupported_operation).\n",
		},
		{
			name: "invalid signature public key",
			path: keyAlgEd25519,
			request: signatureRequest(
				"53167ac3fc4b720daa45b04fc73fe752578fa23a10048422d6904b7f4f7bba5b5b",
				"9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d",
				ed25519SigHex,
			),
			expected: "error(invalid_key).\n",
		},
	}

	Convey("Given a crypto device filesystem", t, func() {
		for _, tc := range cases {
			Convey("when handling "+tc.name, func() {
				got, err := readAllFromCryptoDevice(t, ofs, tc.path, tc.request)
				So(err, ShouldBeNil)
				So(string(got), ShouldEqual, tc.expected)
			})
		}
	})
}

func TestCryptoDeviceFSSignatureRequestValidation(t *testing.T) {
	Convey("Given signature request validation", t, func() {
		cases := []struct {
			name     string
			request  []byte
			expected string
		}{
			{
				name:     "empty request",
				request:  []byte{},
				expected: expectedInvalidRequest,
			},
			{
				name:     "blank request",
				request:  []byte("   \n"),
				expected: expectedInvalidRequest,
			},
			{
				name:     "tab separator",
				request:  []byte("verify\t00 00 00\n"),
				expected: expectedInvalidRequest,
			},
			{
				name:     "control character",
				request:  []byte{'v', 'e', 'r', 'i', 'f', 'y', ' ', 0x7f},
				expected: expectedInvalidRequest,
			},
			{
				name:     "invalid utf8",
				request:  []byte{0xff},
				expected: expectedInvalidRequest,
			},
			{
				name:     "missing signature token",
				request:  []byte("verify 00 00\n"),
				expected: expectedInvalidRequest,
			},
			{
				name:     "invalid public key hex",
				request:  []byte("verify gg 00 00\n"),
				expected: expectedInvalidRequest,
			},
			{
				name:     "invalid message hex",
				request:  []byte("verify 00 gg 00\n"),
				expected: expectedInvalidRequest,
			},
			{
				name:     "invalid signature hex",
				request:  []byte("verify 00 00 gg\n"),
				expected: expectedInvalidRequest,
			},
			{
				name:     "odd-length signature hex",
				request:  []byte("verify 00 00 0\n"),
				expected: expectedInvalidRequest,
			},
		}

		for _, tc := range cases {
			Convey("when handling "+tc.name, func() {
				got, err := handleSignatureRequest(util.KeyAlgEd25519, tc.request)
				So(err, ShouldBeNil)
				So(string(got), ShouldEqual, tc.expected)
			})
		}
	})
}

func TestCryptoDeviceCommitErrors(t *testing.T) {
	Convey("Given a crypto device commit function", t, func() {
		Convey("when request reading fails", func() {
			commit := device{hashAlg: util.HashAlgSha256, kind: deviceKindHash}.makeCommitFunc()
			err := commit(failingReader{}, bytes.NewBuffer(nil))
			So(errors.Is(err, errRead), ShouldBeTrue)
		})

		Convey("when response writing fails", func() {
			commit := device{hashAlg: util.HashAlgSha256, kind: deviceKindHash}.makeCommitFunc()
			err := commit(bytes.NewReader([]byte("hello")), failingWriter{})
			So(errors.Is(err, errWrite), ShouldBeTrue)
		})

		Convey("when device kind is unknown", func() {
			commit := device{}.makeCommitFunc()
			err := commit(bytes.NewReader([]byte("hello")), bytes.NewBuffer(nil))
			So(errors.Is(err, fs.ErrInvalid), ShouldBeTrue)
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
			file, err := ofs.OpenFile("sha256", os.O_RDWR, 0)
			So(err, ShouldBeNil)

			huge := bytes.Repeat([]byte("x"), maxRequestBytes+1)
			_, err = file.(interface{ Write([]byte) (int, error) }).Write(huge)
			So(errors.Is(err, devfile.ErrWriteLimit), ShouldBeTrue)
		})
	})
}

func signatureRequest(pubKey, msg, sig string) []byte {
	return []byte("verify " + pubKey + " " + msg + " " + sig + "\n")
}

var (
	errRead  = errors.New("read failed")
	errWrite = errors.New("write failed")
)

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) {
	return 0, errRead
}

type failingWriter struct{}

func (failingWriter) Write([]byte) (int, error) {
	return 0, errWrite
}

func readAllFromCryptoDevice(t *testing.T, ofs fsiface.OpenFileFS, path string, request []byte) ([]byte, error) {
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

func newSDKContext() sdk.Context {
	stateStore := store.NewCommitMultiStore(dbm.NewMemDB(), log.NewNopLogger(), metrics.NewNoOpMetrics())
	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
	return ctx.WithHeaderInfo(coreheader.Info{
		Height: 42,
		Time:   time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC),
	})
}
