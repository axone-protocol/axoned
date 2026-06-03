package codec

import (
	"bytes"
	"encoding/json"
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

func TestAll(t *testing.T) {
	Convey("Given the codec registry", t, func() {
		So(slices.Contains(All(), "bech32"), ShouldBeTrue)
		So(slices.Contains(All(), "json"), ShouldBeTrue)
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
			payload:    []byte("json([foo=bar])."),
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
			name:           "decode with empty bech32 payload",
			codecName:      "bech32",
			request:        []byte("decode "),
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

		// JSON codec - decode tests
		{
			name:           "json decode object with whitespace",
			codecName:      "json",
			request:        []byte("decode\n{\n  \"foo\": \"bar\",\n  \"ok\": true\n}"),
			expectedOutput: "ok(json([=(foo,bar),=(ok,@(true))])).\n",
		},
		{
			name:           "json decode space-separated payload is not tokenized",
			codecName:      "json",
			request:        []byte("decode {\"foo\":\"bar baz\"}"),
			expectedOutput: "ok(json([=(foo,'bar baz')])).\n",
		},
		{
			name:           "json decode empty payload as null",
			codecName:      "json",
			request:        []byte("decode\n"),
			expectedOutput: "ok(@(null)).\n",
		},
		{
			name:           "json decode malformed payload",
			codecName:      "json",
			request:        []byte("decode\n{&"),
			expectedOutput: "error(syntax_error(json(malformed_json(1)))).\n",
		},
		{
			name:           "json decode rejects trailing document",
			codecName:      "json",
			request:        []byte("decode\n{\"foo\":\"bar\"}{\"foo\":\"bar\"}"),
			expectedOutput: "error(syntax_error(json(malformed_json(14)))).\n",
		},

		// JSON codec - encode tests
		{
			name:           "json encode object",
			codecName:      "json",
			request:        []byte("encode\njson([foo=bar,ok= @(true)])."),
			expectedOutput: "ok('{\"foo\":\"bar\",\"ok\":true}').\n",
		},
		{
			name:           "json encode malformed prolog term",
			codecName:      "json",
			request:        []byte("encode\njson([foo=bar]). trailing"),
			expectedOutput: "error(syntax_error(prolog(malformed_term))).\n",
		},
		{
			name:           "json encode variable fails fast",
			codecName:      "json",
			request:        []byte("encode\nJson."),
			expectedOutput: "error(instantiation_error).\n",
		},
		{
			name:           "json encode invalid JSON term",
			codecName:      "json",
			request:        []byte("encode\nfoo([a=b])."),
			expectedOutput: "error(type_error(json,foo([=(a,b)]))).\n",
		},
		{
			name:           "json encode invalid number",
			codecName:      "json",
			request:        []byte("encode\n1.8e308."),
			expectedOutput: "error(domain_error(json_number,1.8e+308)).\n",
		},
		{
			name:           "json unknown raw command",
			codecName:      "json",
			request:        []byte("unknown\n{}"),
			expectedOutput: "error(invalid_request).\n",
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

func TestJSONCodecDecode(t *testing.T) {
	testCases := []struct {
		name     string
		payload  []byte
		expected string
	}{
		{
			name:     "array with primitive values",
			payload:  []byte(`[1,"two",false,null]`),
			expected: "ok([1.0,two,@(false),@(null)]).\n",
		},
		{
			name:     "boolean true",
			payload:  []byte(`true`),
			expected: "ok(@(true)).\n",
		},
		{
			name:     "number",
			payload:  []byte(`1.5`),
			expected: "ok(1.5).\n",
		},
		{
			name:     "unexpected EOF",
			payload:  []byte(`{"foo":`),
			expected: "error(syntax_error(json(eof))).\n",
		},
	}

	Convey("Given a JSON codec decoder", t, func() {
		codec := &jsonCodec{}
		for _, tc := range testCases {
			Convey(tc.name, func() {
				So(renderTerm(t, codec.Decode(tc.payload)), ShouldEqual, tc.expected)
			})
		}
	})
}

func TestJSONCodecEncode(t *testing.T) {
	testCases := []struct {
		name     string
		payload  []byte
		expected string
	}{
		{
			name:     "atom",
			payload:  []byte("hello."),
			expected: "ok('\"hello\"').\n",
		},
		{
			name:     "empty list",
			payload:  []byte("[]."),
			expected: "ok([]).\n",
		},
		{
			name:     "integer",
			payload:  []byte("42."),
			expected: "ok('42').\n",
		},
		{
			name:     "float",
			payload:  []byte("1.5."),
			expected: "ok('1.5').\n",
		},
		{
			name:     "array",
			payload:  []byte("[foo,42,@(false),@(null),[]]."),
			expected: "ok('[\"foo\",42,false,null,[]]').\n",
		},
		{
			name:     "malformed term",
			payload:  []byte("."),
			expected: "error(syntax_error(prolog(malformed_term))).\n",
		},
	}

	Convey("Given a JSON codec encoder", t, func() {
		codec := &jsonCodec{}
		for _, tc := range testCases {
			Convey(tc.name, func() {
				So(renderTerm(t, codec.Encode(tc.payload)), ShouldEqual, tc.expected)
			})
		}
	})
}

func TestJSONCodecEncodeWriterErrors(t *testing.T) {
	testCases := []struct {
		name       string
		termSource string
		failAt     int
	}{
		{
			name:       "object opening",
			termSource: "json([foo=bar]).",
			failAt:     0,
		},
		{
			name:       "object key",
			termSource: "json([foo=bar]).",
			failAt:     1,
		},
		{
			name:       "object separator",
			termSource: "json([foo=bar]).",
			failAt:     2,
		},
		{
			name:       "object value",
			termSource: "json([foo=bar]).",
			failAt:     3,
		},
		{
			name:       "object comma",
			termSource: "json([foo=bar,baz=qux]).",
			failAt:     4,
		},
		{
			name:       "object closing",
			termSource: "json([foo=bar]).",
			failAt:     4,
		},
		{
			name:       "array opening",
			termSource: "[foo].",
			failAt:     0,
		},
		{
			name:       "array value",
			termSource: "[foo].",
			failAt:     1,
		},
		{
			name:       "array comma",
			termSource: "[foo,bar].",
			failAt:     2,
		},
		{
			name:       "array closing",
			termSource: "[foo].",
			failAt:     2,
		},
	}

	Convey("Given JSON encoder writer failures", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				term, err := parseJSONTerm([]byte(tc.termSource))
				So(err, ShouldBeNil)

				err = encodeTermToJSON(term, &failAtWriter{failAt: tc.failAt, err: io.ErrClosedPipe}, engine.NewEnv())
				So(err, ShouldEqual, io.ErrClosedPipe)
			})
		}
	})
}

func TestJSONCodecErrorTerms(t *testing.T) {
	testCases := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "unmarshal type",
			err:      &json.UnmarshalTypeError{Offset: 7, Value: "number"},
			expected: "syntax_error(json(malformed_json(7,number))).\n",
		},
		{
			name:     "EOF",
			err:      io.EOF,
			expected: "syntax_error(json(eof)).\n",
		},
		{
			name:     "unknown",
			err:      errors.New("boom"),
			expected: "syntax_error(json(unknown)).\n",
		},
	}

	Convey("Given JSON error term mapping", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				So(renderTerm(t, jsonErrorTerm(tc.err)), ShouldEqual, tc.expected)
			})
		}

		Convey("when an encoder error is not a Prolog exception", func() {
			So(renderTerm(t, exceptionFormal(errors.New("boom"))), ShouldEqual, "system_error.\n")
		})
	})
}

func TestJSONMarshalToStreamErrors(t *testing.T) {
	Convey("Given JSON marshaling to a stream", t, func() {
		Convey("when Go JSON rejects the value", func() {
			var buf bytes.Buffer
			err := marshalToJSONStream(math.Inf(1), engine.NewAtom("inf"), &buf, engine.NewEnv())
			So(err, ShouldNotBeNil)
		})
	})
}

func TestJSONEncodeTermValidation(t *testing.T) {
	Convey("Given JSON term validation", t, func() {
		Convey("when an unsupported term implementation is encoded", func() {
			var buf bytes.Buffer
			err := encodeTermToJSON(unknownTerm{}, &buf, engine.NewEnv())
			So(renderTerm(t, exceptionFormal(err)), ShouldEqual, "type_error(json,unknown).\n")
		})
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
		file, err := ofs.OpenFile("bech32", os.O_RDWR, 0)
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

type failAtWriter struct {
	writes int
	failAt int
	err    error
}

func (w *failAtWriter) Write(p []byte) (int, error) {
	if w.writes == w.failAt {
		return 0, w.err
	}

	w.writes++
	return len(p), nil
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

type unknownTerm struct{}

func (t unknownTerm) WriteTerm(w io.Writer, _ *engine.WriteOptions, _ *engine.Env) error {
	_, err := w.Write([]byte("unknown"))
	return err
}

func (t unknownTerm) Compare(_ engine.Term, _ *engine.Env) int {
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
