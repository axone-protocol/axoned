package codec

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/axone-protocol/prolog/v3/engine"
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

func TestAll(t *testing.T) {
	Convey("Given the codec registry", t, func() {
		So(slices.Contains(All(), "bech32"), ShouldBeTrue)
	})
}

func TestCodecDeviceFSOpen(t *testing.T) {
	vfs := NewFS(newSDKContext())

	Convey("Given a codec device filesystem", t, func() {
		_, err := vfs.Open("bech32")
		So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
	})
}

func TestCodecDeviceFSOpenFileValidation(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("codec fs should implement OpenFileFS")
	}

	Convey("Given a codec device filesystem", t, func() {
		Convey("when opening with unsupported mode", func() {
			_, err := ofs.OpenFile("bech32", os.O_RDONLY, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})

		Convey("when opening an unknown path", func() {
			_, err := ofs.OpenFile("unknown", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when opening an escaping path", func() {
			_, err := ofs.OpenFile("../bech32", os.O_RDWR, 0)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func TestMakeCommitFunc(t *testing.T) {
	Convey("Given a codec commit function", t, func() {
		Convey("when reading the request fails", func() {
			commit := makeCommitFunc(stubCodec{})
			err := commit(errReader{err: io.ErrUnexpectedEOF}, &bytes.Buffer{})
			So(err, ShouldEqual, io.ErrUnexpectedEOF)
		})

		Convey("when rendering the response fails", func() {
			commit := makeCommitFunc(stubCodec{decodeTerm: badTerm{err: errTestTerm}})
			err := commit(bytes.NewBufferString("decode anything"), &bytes.Buffer{})
			So(err, ShouldEqual, errTestTerm)
		})

		Convey("when writing the response fails", func() {
			commit := makeCommitFunc(stubCodec{decodeTerm: atomOK.Apply(engine.NewAtom("ok"))})
			err := commit(bytes.NewBufferString("decode anything"), errWriter{err: io.ErrClosedPipe})
			So(err, ShouldEqual, io.ErrClosedPipe)
		})
	})
}

func TestHandleRequest(t *testing.T) {
	codec := stubCodec{decodeTerm: engine.NewAtom("decoded"), encodeTerm: engine.NewAtom("encoded")}

	Convey("Given request dispatch", t, func() {
		Convey("when decode is requested with one argument", func() {
			term := handleRequest(codec, []byte("decode value"))
			So(term, ShouldEqual, codec.decodeTerm)
		})

		Convey("when encode is requested with two arguments", func() {
			term := handleRequest(codec, []byte("encode hrp deadbeef"))
			So(term, ShouldEqual, codec.encodeTerm)
		})
	})
}

func TestNormalizeRequestLine(t *testing.T) {
	testCases := []struct {
		name       string
		request    []byte
		expected   []byte
		expectedOK bool
	}{
		{
			name:       "empty request",
			request:    nil,
			expectedOK: false,
		},
		{
			name:       "empty after trimming line ending and spaces",
			request:    []byte("  \r\n"),
			expectedOK: false,
		},
		{
			name:       "invalid utf8",
			request:    []byte{0xff},
			expectedOK: false,
		},
		{
			name:       "tab is rejected",
			request:    []byte("decode\tvalue"),
			expectedOK: false,
		},
		{
			name:       "delete control char is rejected",
			request:    append([]byte("decode "), 0x7f),
			expectedOK: false,
		},
		{
			name:       "valid request is trimmed",
			request:    []byte("  decode value  \r\n"),
			expected:   []byte("decode value"),
			expectedOK: true,
		},
	}

	Convey("Given request normalization", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				line, ok := normalizeRequestLine(tc.request)
				So(ok, ShouldEqual, tc.expectedOK)
				So(line, ShouldResemble, tc.expected)
			})
		}
	})
}

func TestSplitRequestLine(t *testing.T) {
	testCases := []struct {
		name     string
		line     []byte
		expected [][]byte
	}{
		{
			name:     "empty line",
			line:     nil,
			expected: [][]byte{},
		},
		{
			name:     "spaces only",
			line:     []byte("   "),
			expected: [][]byte{},
		},
		{
			name:     "multiple tokens with repeated separators",
			line:     []byte("encode  hrp   deadbeef"),
			expected: [][]byte{[]byte("encode"), []byte("hrp"), []byte("deadbeef")},
		},
		{
			name:     "trailing spaces are ignored",
			line:     []byte("decode value   "),
			expected: [][]byte{[]byte("decode"), []byte("value")},
		},
	}

	Convey("Given request tokenization", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				So(splitRequestLine(tc.line), ShouldResemble, tc.expected)
			})
		}
	})
}

func TestCodecDeviceFSFunctional(t *testing.T) {
	vfs := NewFS(newSDKContext())
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("codec fs should implement OpenFileFS")
	}

	readAll := func(codecName string, request []byte) ([]byte, error) {
		t.Helper()

		file, err := ofs.OpenFile(codecName, os.O_RDWR, 0)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		if len(request) > 0 {
			writer, ok := file.(interface{ Write([]byte) (int, error) })
			if !ok {
				t.Fatal("codec file should implement Write")
			}

			if _, err := writer.Write(request); err != nil {
				return nil, err
			}
		}

		return io.ReadAll(file)
	}

	testCases := []struct {
		name           string
		codecName      string
		request        []byte
		expectedOutput string
	}{
		// Protocol-level validation (codec-agnostic)
		{
			name:           "empty request",
			codecName:      "bech32",
			request:        nil,
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "unknown command",
			codecName:      "bech32",
			request:        []byte("unknown"),
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "malformed UTF-8",
			codecName:      "bech32",
			request:        []byte{0xff},
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "tab separator not allowed",
			codecName:      "bech32",
			request:        []byte("encode\taxone\t00"),
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "decode with insufficient arguments",
			codecName:      "bech32",
			request:        []byte("decode"),
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "encode with insufficient arguments",
			codecName:      "bech32",
			request:        []byte("encode axone"),
			expectedOutput: "error(invalid_request).\n",
		},
		{
			name:           "encode with too many arguments",
			codecName:      "bech32",
			request:        []byte("encode hrp hex extra"),
			expectedOutput: "error(invalid_request).\n",
		},

		// Bech32 codec - decode tests
		{
			name:           "bech32 decode valid address",
			codecName:      "bech32",
			request:        []byte("decode axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4"),
			expectedOutput: "ok(-(axone,[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30])).\n",
		},
		{
			name:           "bech32 decode with whitespace normalization",
			codecName:      "bech32",
			request:        []byte("  decode   axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4  \r\n"),
			expectedOutput: "ok(-(axone,[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30])).\n",
		},
		{
			name:           "bech32 decode invalid payload",
			codecName:      "bech32",
			request:        []byte("decode bad"),
			expectedOutput: "error(invalid_bech32).\n",
		},

		// Bech32 codec - encode tests
		{
			name:           "bech32 encode valid bytes",
			codecName:      "bech32",
			request:        []byte("encode axone a3a717f4a2af31a2aa0fb58d44868da81238f71e"),
			expectedOutput: "ok(axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4).\n",
		},
		{
			name:           "bech32 encode with uppercase hex",
			codecName:      "bech32",
			request:        []byte("encode axone A3A717F4A2AF31A2AA0FB58D44868DA81238F71E"),
			expectedOutput: "ok(axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4).\n",
		},
		{
			name:           "bech32 encode invalid hex characters",
			codecName:      "bech32",
			request:        []byte("encode axone 0011zz"),
			expectedOutput: "error(invalid_bytes).\n",
		},
		{
			name:           "bech32 encode odd-length hex",
			codecName:      "bech32",
			request:        []byte("encode axone 123"),
			expectedOutput: "error(invalid_bytes).\n",
		},
	}

	Convey("Given a codec device filesystem", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				response, err := readAll(tc.codecName, tc.request)
				So(err, ShouldBeNil)
				So(string(response), ShouldEqual, tc.expectedOutput)
			})
		}
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

var errTestTerm = errors.New("term write failure")

type stubCodec struct {
	decodeTerm engine.Term
	encodeTerm engine.Term
}

func (c stubCodec) Name() string {
	return "stub"
}

func (c stubCodec) Decode(_ [][]byte) engine.Term {
	if c.decodeTerm != nil {
		return c.decodeTerm
	}

	return engine.NewAtom("decoded")
}

func (c stubCodec) Encode(_ [][]byte) engine.Term {
	if c.encodeTerm != nil {
		return c.encodeTerm
	}

	return engine.NewAtom("encoded")
}

type errReader struct {
	err error
}

func (r errReader) Read(_ []byte) (int, error) {
	return 0, r.err
}

type errWriter struct {
	err error
}

func (w errWriter) Write(_ []byte) (int, error) {
	return 0, w.err
}

type badTerm struct {
	err error
}

func (t badTerm) WriteTerm(_ io.Writer, _ *engine.WriteOptions, _ *engine.Env) error {
	return t.err
}

func (t badTerm) Compare(_ engine.Term, _ *engine.Env) int {
	return 0
}
