package codec

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"math"
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
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	fsiface "github.com/axone-protocol/axoned/v15/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/prologterm"
	logictypes "github.com/axone-protocol/axoned/v15/x/logic/types"
)

const (
	expectedInvalidRequest  = "error(invalid_request).\n"
	expectedMalformedProlog = "error(syntax_error(prolog(malformed_term))).\n"
	testJSONFooBarTerm      = "json([foo=bar])."
	testListFooTerm         = "[foo]."
)

func TestAll(t *testing.T) {
	Convey("Given the codec registry", t, func() {
		So(slices.Contains(All(), codecNameBech32), ShouldBeTrue)
		So(slices.Contains(All(), codecNameJSON), ShouldBeTrue)
		So(slices.Contains(All(), codecNameText), ShouldBeTrue)
	})
}

func TestCodecDeviceFSOpen(t *testing.T) {
	vfs := NewFS(newSDKContext())

	Convey("Given a codec device filesystem", t, func() {
		_, err := vfs.Open(codecNameBech32)
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
			_, err := ofs.OpenFile(codecNameBech32, os.O_RDONLY, 0)
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
			commit := makeCommitFunc(&stubCodec{})
			err := commit(errReader{err: io.ErrUnexpectedEOF}, &bytes.Buffer{})
			So(err, ShouldEqual, io.ErrUnexpectedEOF)
		})

		Convey("when rendering the response fails", func() {
			commit := makeCommitFunc(&stubCodec{decodeTerm: badTerm{err: errTestTerm}})
			err := commit(bytes.NewBufferString("decode anything"), &bytes.Buffer{})
			So(err, ShouldEqual, errTestTerm)
		})

		Convey("when writing the response fails", func() {
			commit := makeCommitFunc(&stubCodec{decodeTerm: atomOK.Apply(engine.NewAtom("ok"))})
			err := commit(bytes.NewBufferString("decode anything"), errWriter{err: io.ErrClosedPipe})
			So(err, ShouldEqual, io.ErrClosedPipe)
		})
	})
}

func TestHandleRequest(t *testing.T) {
	Convey("Given request dispatch", t, func() {
		codec := &stubCodec{decodeTerm: engine.NewAtom("decoded"), encodeTerm: engine.NewAtom("encoded")}

		Convey("when decode is requested", func() {
			term := handleRequest(codec, []byte("decode value"))
			So(term, ShouldEqual, codec.decodeTerm)
			So(codec.decodeInput, ShouldResemble, []byte("value"))
		})

		Convey("when encode is requested", func() {
			term := handleRequest(codec, []byte("encode hrp deadbeef"))
			So(term, ShouldEqual, codec.encodeTerm)
			So(codec.encodeInput, ShouldResemble, []byte("hrp deadbeef"))
		})

		Convey("when the payload contains spaces and newlines", func() {
			term := handleRequest(codec, []byte("decode\n{\"foo\":\"bar baz\"}\n"))
			So(term, ShouldEqual, codec.decodeTerm)
			So(codec.decodeInput, ShouldResemble, []byte("{\"foo\":\"bar baz\"}\n"))
		})
	})
}

func TestSplitRequestCommand(t *testing.T) {
	testCases := []struct {
		name       string
		request    []byte
		command    []byte
		payload    []byte
		expectedOK bool
	}{
		{
			name:       "empty request",
			request:    nil,
			expectedOK: false,
		},
		{
			name:       "empty after trimming leading spaces",
			request:    []byte("   "),
			expectedOK: false,
		},
		{
			name:       "empty command before line ending",
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
			name:       "delete control char in command is rejected",
			request:    []byte("de\x7fcode value"),
			expectedOK: false,
		},
		{
			name:       "space separator",
			request:    []byte("  decode value  \r\n"),
			command:    []byte("decode"),
			payload:    []byte("value  \r\n"),
			expectedOK: true,
		},
		{
			name:       "line separator",
			request:    []byte("encode\njson([foo=bar])."),
			command:    []byte("encode"),
			payload:    []byte(testJSONFooBarTerm),
			expectedOK: true,
		},
		{
			name:       "CRLF separator",
			request:    []byte("decode\r\n{\"foo\":\"bar\"}"),
			command:    []byte("decode"),
			payload:    []byte("{\"foo\":\"bar\"}"),
			expectedOK: true,
		},
		{
			name:       "CR without LF is rejected",
			request:    []byte("decode\r{\"foo\":\"bar\"}"),
			expectedOK: false,
		},
	}

	Convey("Given request command parsing", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				command, payload, ok := splitRequestCommand(tc.request)
				So(ok, ShouldEqual, tc.expectedOK)
				So(command, ShouldResemble, tc.command)
				So(payload, ShouldResemble, tc.payload)
			})
		}
	})
}

func TestSplitBech32Payload(t *testing.T) {
	testCases := []struct {
		name     string
		input    []byte
		expected [][]byte
		ok       bool
	}{
		{
			name:  "empty payload",
			input: nil,
			ok:    false,
		},
		{
			name:  "spaces only",
			input: []byte("   "),
			ok:    false,
		},
		{
			name:     "multiple tokens with repeated separators",
			input:    []byte("hrp   deadbeef"),
			expected: [][]byte{[]byte("hrp"), []byte("deadbeef")},
			ok:       true,
		},
		{
			name:     "trailing spaces are ignored",
			input:    []byte("value   "),
			expected: [][]byte{[]byte("value")},
			ok:       true,
		},
		{
			name:  "tabs are rejected",
			input: []byte("hrp\tdeadbeef"),
			ok:    false,
		},
		{
			name:  "invalid utf8 is rejected",
			input: []byte{0xff},
			ok:    false,
		},
		{
			name:  "control bytes are rejected",
			input: []byte("hrp \x7f"),
			ok:    false,
		},
	}

	Convey("Given Bech32 payload tokenization", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				tokens, ok := splitBech32Payload(tc.input)
				So(ok, ShouldEqual, tc.ok)
				So(tokens, ShouldResemble, tc.expected)
			})
		}
	})
}

func TestSplitSpaceTokens(t *testing.T) {
	Convey("Given Bech32 space tokenization", t, func() {
		So(splitSpaceTokens([]byte("   ")), ShouldResemble, [][]byte{})
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
			codecName:      codecNameBech32,
			request:        nil,
			expectedOutput: expectedInvalidRequest,
		},
		{
			name:           "unknown command",
			codecName:      codecNameBech32,
			request:        []byte("unknown"),
			expectedOutput: expectedInvalidRequest,
		},
		{
			name:           "malformed UTF-8",
			codecName:      codecNameBech32,
			request:        []byte{0xff},
			expectedOutput: expectedInvalidRequest,
		},
		{
			name:           "tab separator not allowed",
			codecName:      codecNameBech32,
			request:        []byte("encode\taxone\t00"),
			expectedOutput: expectedInvalidRequest,
		},
		{
			name:           "decode with insufficient arguments",
			codecName:      codecNameBech32,
			request:        []byte("decode"),
			expectedOutput: expectedInvalidRequest,
		},
		{
			name:           "decode with empty bech32 payload",
			codecName:      codecNameBech32,
			request:        []byte("decode "),
			expectedOutput: expectedInvalidRequest,
		},
		{
			name:           "encode with insufficient arguments",
			codecName:      codecNameBech32,
			request:        []byte("encode axone"),
			expectedOutput: expectedInvalidRequest,
		},
		{
			name:           "encode with too many arguments",
			codecName:      codecNameBech32,
			request:        []byte("encode hrp hex extra"),
			expectedOutput: expectedInvalidRequest,
		},

		// Bech32 codec - decode tests
		{
			name:           "bech32 decode valid address",
			codecName:      codecNameBech32,
			request:        []byte("decode axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4"),
			expectedOutput: "ok(-(axone,[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30])).\n",
		},
		{
			name:           "bech32 decode with whitespace normalization",
			codecName:      codecNameBech32,
			request:        []byte("  decode   axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4  \r\n"),
			expectedOutput: "ok(-(axone,[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30])).\n",
		},
		{
			name:           "bech32 decode invalid payload",
			codecName:      codecNameBech32,
			request:        []byte("decode bad"),
			expectedOutput: "error(invalid_bech32).\n",
		},

		// Bech32 codec - encode tests
		{
			name:           "bech32 encode valid bytes",
			codecName:      codecNameBech32,
			request:        []byte("encode axone a3a717f4a2af31a2aa0fb58d44868da81238f71e"),
			expectedOutput: "ok(axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4).\n",
		},
		{
			name:           "bech32 encode with uppercase hex",
			codecName:      codecNameBech32,
			request:        []byte("encode axone A3A717F4A2AF31A2AA0FB58D44868DA81238F71E"),
			expectedOutput: "ok(axone15wn30a9z4uc692s0kkx5fp5d4qfr3ac77gvjg4).\n",
		},
		{
			name:           "bech32 encode invalid hex characters",
			codecName:      codecNameBech32,
			request:        []byte("encode axone 0011zz"),
			expectedOutput: "error(invalid_bytes).\n",
		},
		{
			name:           "bech32 encode odd-length hex",
			codecName:      codecNameBech32,
			request:        []byte("encode axone 123"),
			expectedOutput: "error(invalid_bytes).\n",
		},

		// JSON codec - decode tests
		{
			name:           "json decode object with whitespace",
			codecName:      codecNameJSON,
			request:        []byte("decode\n{\n  \"foo\": \"bar\",\n  \"ok\": true\n}"),
			expectedOutput: "ok(json([=(foo,bar),=(ok,@(true))])).\n",
		},
		{
			name:           "json decode space-separated payload is not tokenized",
			codecName:      codecNameJSON,
			request:        []byte("decode {\"foo\":\"bar baz\"}"),
			expectedOutput: "ok(json([=(foo,'bar baz')])).\n",
		},
		{
			name:           "json decode empty payload as null",
			codecName:      codecNameJSON,
			request:        []byte("decode\n"),
			expectedOutput: "ok(@(null)).\n",
		},
		{
			name:           "json decode malformed payload",
			codecName:      codecNameJSON,
			request:        []byte("decode\n{&"),
			expectedOutput: "error(syntax_error(json(malformed_json(1)))).\n",
		},
		{
			name:           "json decode rejects trailing document",
			codecName:      codecNameJSON,
			request:        []byte("decode\n{\"foo\":\"bar\"}{\"foo\":\"bar\"}"),
			expectedOutput: "error(syntax_error(json(malformed_json(14)))).\n",
		},

		// JSON codec - encode tests
		{
			name:           "json encode object",
			codecName:      codecNameJSON,
			request:        []byte("encode\njson([foo=bar,ok= @(true)])."),
			expectedOutput: "ok('{\"foo\":\"bar\",\"ok\":true}').\n",
		},
		{
			name:           "json encode malformed prolog term",
			codecName:      codecNameJSON,
			request:        []byte("encode\njson([foo=bar]). trailing"),
			expectedOutput: expectedMalformedProlog,
		},
		{
			name:           "json encode variable fails fast",
			codecName:      codecNameJSON,
			request:        []byte("encode\nJson."),
			expectedOutput: "error(instantiation_error).\n",
		},
		{
			name:           "json encode invalid JSON term",
			codecName:      codecNameJSON,
			request:        []byte("encode\nfoo([a=b])."),
			expectedOutput: "error(type_error(json,foo([=(a,b)]))).\n",
		},
		{
			name:           "json encode invalid number",
			codecName:      codecNameJSON,
			request:        []byte("encode\n1.8e308."),
			expectedOutput: "error(domain_error(json_number,1.8e+308)).\n",
		},
		{
			name:           "json unknown raw command",
			codecName:      codecNameJSON,
			request:        []byte("unknown\n{}"),
			expectedOutput: expectedInvalidRequest,
		},

		// Text codec tests
		{
			name:           "text encode utf8",
			codecName:      codecNameText,
			request:        []byte("encode\ntext(utf8, aap)."),
			expectedOutput: "ok([97,97,112]).\n",
		},
		{
			name:           "text decode utf8",
			codecName:      codecNameText,
			request:        []byte("decode\nbytes(utf8, [97,97,112])."),
			expectedOutput: "ok([a,a,p]).\n",
		},
		{
			name:           "text encode text",
			codecName:      codecNameText,
			request:        []byte("encode\ntext(text, 'ù')."),
			expectedOutput: "ok([195,185]).\n",
		},
		{
			name:           "text decode text",
			codecName:      codecNameText,
			request:        []byte("decode\nbytes(text, [195,185])."),
			expectedOutput: "ok([ù]).\n",
		},
		{
			name:           "text encode octet",
			codecName:      codecNameText,
			request:        []byte("encode\ntext(octet, 'ù')."),
			expectedOutput: "ok([249]).\n",
		},
		{
			name:           "text decode octet",
			codecName:      codecNameText,
			request:        []byte("decode\nbytes(octet, [249])."),
			expectedOutput: "ok([ù]).\n",
		},
		{
			name:           "text encode utf-16le",
			codecName:      codecNameText,
			request:        []byte("encode\ntext('utf-16le', '今日は')."),
			expectedOutput: "ok([202,78,229,101,111,48]).\n",
		},
		{
			name:           "text decode utf-16be",
			codecName:      codecNameText,
			request:        []byte("decode\nbytes('utf-16be', [0,97,0,97,0,112])."),
			expectedOutput: "ok([a,a,p]).\n",
		},
		{
			name:           "text encode invalid charset",
			codecName:      codecNameText,
			request:        []byte("encode\ntext(foo, aap)."),
			expectedOutput: "error(type_error(charset,foo)).\n",
		},
		{
			name:           "text decode invalid byte",
			codecName:      codecNameText,
			request:        []byte("decode\nbytes(latin2, [400])."),
			expectedOutput: "error(type_error(byte,400)).\n",
		},
		{
			name:           "text encode malformed term",
			codecName:      codecNameText,
			request:        []byte("encode\ntext(utf8, aap). trailing"),
			expectedOutput: expectedMalformedProlog,
		},
		{
			name:           "text encode invalid request term",
			codecName:      codecNameText,
			request:        []byte("encode\nfoo(utf8, aap)."),
			expectedOutput: expectedInvalidRequest,
		},
		{
			name:           "text decode invalid bytes term",
			codecName:      codecNameText,
			request:        []byte("decode\nbytes(utf8, foo(bar))."),
			expectedOutput: "error(type_error(list,foo(bar))).\n",
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

func TestCodecDeviceFSGasMetering(t *testing.T) {
	request := []byte("decode bad")
	expectedOutput := "error(invalid_bech32).\n"

	ctx := newSDKContext().
		WithGasMeter(storetypes.NewGasMeter(1_000)).
		WithValue(logictypes.IOCoeffContextKey, uint64(2))

	vfs := NewFS(ctx)
	ofs, ok := vfs.(fsiface.OpenFileFS)
	if !ok {
		t.Fatal("codec fs should implement OpenFileFS")
	}

	Convey("Given a codec device filesystem with I/O gas metering", t, func() {
		file, err := ofs.OpenFile(codecNameBech32, os.O_RDWR, 0)
		So(err, ShouldBeNil)
		defer file.Close()

		writer, ok := file.(interface{ Write([]byte) (int, error) })
		if !ok {
			t.Fatal("codec file should implement Write")
		}

		_, err = writer.Write(request)
		So(err, ShouldBeNil)

		response, err := io.ReadAll(file)
		So(err, ShouldBeNil)
		So(string(response), ShouldEqual, expectedOutput)
		So(ctx.GasMeter().GasConsumed(), ShouldEqual, uint64(len(request)+len(expectedOutput))*2)
	})
}

func TestCodecDeviceFSIOGasHelpers(t *testing.T) {
	Convey("Given codec I/O gas helper functions", t, func() {
		Convey("when transferred bytes are non-positive", func() {
			gasMeter := storetypes.NewInfiniteGasMeter()

			consumeTransferredIOGas(gasMeter, 0, 2)
			consumeTransferredIOGas(gasMeter, -5, 2)

			So(gasMeter.GasConsumed(), ShouldEqual, uint64(0))
		})

		Convey("when transferred bytes are positive", func() {
			gasMeter := storetypes.NewInfiniteGasMeter()

			consumeTransferredIOGas(gasMeter, 4, 2)

			So(gasMeter.GasConsumed(), ShouldEqual, uint64(8))
		})

		Convey("when units are zero", func() {
			gasMeter := storetypes.NewInfiniteGasMeter()

			consumeIOGas(gasMeter, 0, 3)

			So(gasMeter.GasConsumed(), ShouldEqual, uint64(0))
		})

		Convey("when coefficient is zero", func() {
			gasMeter := storetypes.NewInfiniteGasMeter()

			consumeIOGas(gasMeter, 7, 0)

			So(gasMeter.GasConsumed(), ShouldEqual, uint64(7))
		})

		Convey("when multiplication overflows", func() {
			gasMeter := storetypes.NewInfiniteGasMeter()

			consumeIOGas(gasMeter, 2, math.MaxUint64)

			So(gasMeter.GasConsumed(), ShouldEqual, uint64(math.MaxUint64))
		})
	})
}

func TestMultiplyUint64Overflow(t *testing.T) {
	testCases := []struct {
		name             string
		a                uint64
		b                uint64
		expected         uint64
		expectedOverflow bool
	}{
		{
			name:             "zero operand does not overflow",
			a:                0,
			b:                42,
			expected:         0,
			expectedOverflow: false,
		},
		{
			name:             "regular multiplication",
			a:                6,
			b:                7,
			expected:         42,
			expectedOverflow: false,
		},
		{
			name:             "overflow",
			a:                2,
			b:                math.MaxUint64,
			expected:         0,
			expectedOverflow: true,
		},
	}

	Convey("Given uint64 multiplication overflow detection", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				got, overflow := multiplyUint64Overflow(tc.a, tc.b)
				So(got, ShouldEqual, tc.expected)
				So(overflow, ShouldEqual, tc.expectedOverflow)
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
	decodeTerm  engine.Term
	encodeTerm  engine.Term
	decodeInput []byte
	encodeInput []byte
}

func (c stubCodec) Name() string {
	return "stub"
}

func (c *stubCodec) Decode(input []byte) engine.Term {
	c.decodeInput = append([]byte(nil), input...)
	if c.decodeTerm != nil {
		return c.decodeTerm
	}

	return engine.NewAtom("decoded")
}

func (c *stubCodec) Encode(input []byte) engine.Term {
	c.encodeInput = append([]byte(nil), input...)
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

func renderTerm(t *testing.T, term engine.Term) string {
	t.Helper()

	bs, err := prologterm.Render(term, true)
	if err != nil {
		t.Fatal(err)
	}
	return string(bs)
}
