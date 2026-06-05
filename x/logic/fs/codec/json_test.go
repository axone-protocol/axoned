package codec

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"math"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"

	. "github.com/smartystreets/goconvey/convey"
)

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
			expected: expectedMalformedProlog,
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
			termSource: testJSONFooBarTerm,
			failAt:     0,
		},
		{
			name:       "object key",
			termSource: testJSONFooBarTerm,
			failAt:     1,
		},
		{
			name:       "object separator",
			termSource: testJSONFooBarTerm,
			failAt:     2,
		},
		{
			name:       "object value",
			termSource: testJSONFooBarTerm,
			failAt:     3,
		},
		{
			name:       "object comma",
			termSource: "json([foo=bar,baz=qux]).",
			failAt:     4,
		},
		{
			name:       "object closing",
			termSource: testJSONFooBarTerm,
			failAt:     4,
		},
		{
			name:       "array opening",
			termSource: testListFooTerm,
			failAt:     0,
		},
		{
			name:       "array value",
			termSource: testListFooTerm,
			failAt:     1,
		},
		{
			name:       "array comma",
			termSource: "[foo,bar].",
			failAt:     2,
		},
		{
			name:       "array closing",
			termSource: testListFooTerm,
			failAt:     2,
		},
	}

	Convey("Given JSON encoder writer failures", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				term, err := parseCodecTerm([]byte(tc.termSource))
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

type unknownTerm struct{}

func (t unknownTerm) WriteTerm(w io.Writer, _ *engine.WriteOptions, _ *engine.Env) error {
	_, err := w.Write([]byte("unknown"))
	return err
}

func (t unknownTerm) Compare(_ engine.Term, _ *engine.Env) int {
	return 0
}
