package predicate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/axone-protocol/prolog/engine"
	"github.com/samber/lo"

	"github.com/axone-protocol/axoned/v10/x/logic/prolog"
)

var (
	// AtomSyntaxErrorJSON represents a syntax error related to JSON.
	AtomSyntaxErrorJSON = engine.NewAtom("json")

	// AtomMalformedJSON represents a specific type of JSON syntax error where the JSON is malformed.
	AtomMalformedJSON = engine.NewAtom("malformed_json")

	// AtomEOF represents a specific type of JSON syntax error where an unexpected end-of-file occurs.
	AtomEOF = engine.NewAtom("eof")

	// AtomUnknown represents an unknown or unspecified syntax error.
	AtomUnknown = engine.NewAtom("unknown")

	// AtomValidJSONNumber is the atom denoting a valid JSON number.
	AtomValidJSONNumber = engine.NewAtom("json_number")
)

var (
	errWrongStreamType = errors.New("wrong stream type")
	errWrongIOMode     = errors.New("wrong i/o mode")
	errPastEndOfStream = errors.New("past end of stream")
	errInvalidUTF8     = errors.New("invalid UTF-8")
)

var (
	operationInput  = engine.NewAtom("input")
	operationOutput = engine.NewAtom("output")
)

var (
	permissionTypeStream          = engine.NewAtom("stream")
	permissionTypeTextStream      = engine.NewAtom("text_stream")
	permissionTypePastEndOfStream = engine.NewAtom("past_end_of_stream")
)

// JSONRead is a predicate that reads a JSON from a stream and unifies it with a Prolog term.
//
// See json_prolog/2 for the canonical representation of the JSON term.
//
// The signature is as follows:
//
//	json_read(+Stream, ?Term) is det
//
// Where:
//   - Stream is the input stream from which the JSON is read.
//   - Term is the Prolog term that represents the JSON structure.
func JSONRead(vm *engine.VM, stream, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	is, err := prolog.AssertStream(stream, env)
	if err != nil {
		return engine.Error(err)
	}

	decoder := newTextStreamDecoder(is)
	decoded, err := decodeJSONToTerm(decoder, env)
	if err != nil {
		return engine.Error(err)
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return engine.Error(
			engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomMalformedJSON.Apply(engine.Integer(decoder.InputOffset()))), env))
	}

	return engine.Unify(vm, term, decoded, cont, env)
}

// JSONWrite is a predicate that writes a Prolog term as a JSON to a stream.
//
// The JSON object is of the same format as produced by json_read/2.
//
// The signature is as follows:
//
//	json_write(+Stream, +Term) is det
//
// Where:
//   - Stream is the output stream to which the JSON is written.
//   - Term is the Prolog term that represents the JSON structure.
func JSONWrite(_ *engine.VM, stream, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	os, err := prolog.AssertStream(stream, env)
	if err != nil {
		return engine.Error(err)
	}

	buf := newTextStreamWriter(os)
	if err := encodeTermToJSON(term, buf, env); err != nil {
		return engine.Error(err)
	}

	return cont(env)
}

// JSONProlog is a predicate that unifies a JSON into a prolog term and vice versa.
//
// The signature is as follows:
//
//	json_prolog(?Json, ?Term) is det
//
// Where:
//   - Json is the textual representation of the JSON, as either an atom, a list of character codes, or a list of characters.
//   - Term is the Prolog term that represents the JSON structure.
//
// # JSON canonical representation
//
// The canonical representation for Term is:
//   - A JSON object is mapped to a Prolog term json(NameValueList), where NameValueList is a list of Name=Value key values.
//     Name is an atom created from the JSON string.
//   - A JSON array is mapped to a Prolog list of JSON values.
//   - A JSON string is mapped to a Prolog atom.
//   - A JSON number is mapped to a Prolog number.
//   - The JSON constants true and false are mapped to @(true) and @(false).
//   - The JSON constant null is mapped to the Prolog term @(null).
//
// # Examples:
//
//	# JSON conversion to Prolog.
//	- json_prolog('{"foo": "bar"}', json([foo=bar])).
func JSONProlog(vm *engine.VM, j, p engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	switch {
	case prolog.IsGround(j, env):
		payload, err := prolog.TextTermToString(j, env)
		if err != nil {
			return engine.Error(err)
		}
		is := engine.NewInputTextStream(strings.NewReader(payload))
		defer is.Close()

		return JSONRead(vm, is, p, cont, env)
	default:
		var buf bytes.Buffer
		os := engine.NewOutputTextStream(&buf)
		defer os.Close()

		return JSONWrite(vm, os, p, func(env *engine.Env) *engine.Promise {
			return engine.Unify(vm, j, prolog.StringToAtom(buf.String()), cont, env)
		}, env)
	}
}

func encodeTermToJSON(term engine.Term, writer *textStreamWriter, env *engine.Env) (err error) {
	switch t := term.(type) {
	case engine.Atom:
		if term == prolog.AtomEmptyList {
			if err := writeToStream(writer, []byte("[]"), env); err != nil {
				return err
			}
		} else {
			return marshalToStream(t.String(), term, writer, env)
		}
	case engine.Integer:
		return marshalToStream(t, term, writer, env)
	case engine.Float:
		float, err := strconv.ParseFloat(t.String(), 64)
		if err != nil {
			return prologErrorToException(t, err, env)
		}
		return marshalToStream(float, term, writer, env)
	case engine.Compound:
		return encodeCompoundToJSON(t, writer, env)
	case engine.Variable:
		return engine.InstantiationError(env)
	default:
		return engine.TypeError(prolog.AtomTypeJSON, term, env)
	}

	return nil
}

func encodeCompoundToJSON(term engine.Compound, writer *textStreamWriter, env *engine.Env) error {
	switch {
	case term.Functor() == prolog.AtomDot:
		return encodeArrayToJSON(term, writer, env)
	case term.Functor() == prolog.AtomJSON:
		return encodeObjectToJSON(term, writer, env)
	case prolog.JSONBool(true).Compare(term, env) == 0:
		if err := writeToStream(writer, []byte("true"), env); err != nil {
			return err
		}
	case prolog.JSONBool(false).Compare(term, env) == 0:
		if err := writeToStream(writer, []byte("false"), env); err != nil {
			return err
		}
	case prolog.JSONNull().Compare(term, env) == 0:
		if err := writeToStream(writer, []byte("null"), env); err != nil {
			return err
		}
	default:
		return engine.TypeError(prolog.AtomTypeJSON, term, env)
	}

	return nil
}

func encodeObjectToJSON(term engine.Compound, writer *textStreamWriter, env *engine.Env) error {
	if _, err := prolog.AssertJSON(term, env); err != nil {
		return err
	}
	if err := writeToStream(writer, []byte("{"), env); err != nil {
		return err
	}
	if err := prolog.ForEach(term.Arg(0), env, func(t engine.Term, hasNext bool) error {
		k, v, err := prolog.AssertKeyValue(t, env)
		if err != nil {
			return err
		}
		if err := marshalToStream(k.String(), term, writer, env); err != nil {
			return err
		}
		if err := writeToStream(writer, []byte(":"), env); err != nil {
			return err
		}
		if err := encodeTermToJSON(v, writer, env); err != nil {
			return err
		}

		if hasNext {
			if err := writeToStream(writer, []byte(","), env); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	if err := writeToStream(writer, []byte("}"), env); err != nil {
		return err
	}
	return nil
}

func encodeArrayToJSON(term engine.Compound, writer *textStreamWriter, env *engine.Env) error {
	if err := writeToStream(writer, []byte("["), env); err != nil {
		return err
	}
	if err := prolog.ForEach(term, env, func(t engine.Term, hasNext bool) error {
		err := encodeTermToJSON(t, writer, env)
		if err != nil {
			return err
		}

		if hasNext {
			if err := writeToStream(writer, []byte(","), env); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}
	if err := writeToStream(writer, []byte("]"), env); err != nil {
		return err
	}

	return nil
}

func marshalToStream(data any, term engine.Term, writer *textStreamWriter, env *engine.Env) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return prologErrorToException(term, err, env)
	}
	if err := writeToStream(writer, bs, env); err != nil {
		return err
	}
	return nil
}

func writeToStream(writer *textStreamWriter, data []byte, env *engine.Env) error {
	if _, err := writer.Write(data); err != nil {
		return prologErrorToException(writer.stream, err, env)
	}
	return nil
}

func prologErrorToException(culprit engine.Term, err error, env *engine.Env) engine.Exception {
	if _, ok := lo.ErrorsAs[*strconv.NumError](err); ok {
		return engine.DomainError(AtomValidJSONNumber, culprit, env)
	}

	switch {
	case errors.Is(err, io.EOF):
		return engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomEOF), env)
	case err.Error() == errWrongIOMode.Error():
		return engine.PermissionError(operationOutput, permissionTypeStream, culprit, env)
	case err.Error() == errWrongStreamType.Error():
		return engine.PermissionError(operationOutput, permissionTypeTextStream, culprit, env)
	case err.Error() == errPastEndOfStream.Error():
		return engine.PermissionError(operationOutput, permissionTypePastEndOfStream, culprit, env)
	}

	return prolog.WithError(
		engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomUnknown), env), err, env)
}

func decodeJSONToTerm(decoder *textStreamDecoder, env *engine.Env) (engine.Term, error) {
	t, err := nextToken(decoder, env)
	if errors.Is(err, io.EOF) {
		return prolog.JSONNull(), nil
	}
	if err != nil {
		return nil, err
	}

	switch t := t.(type) {
	case json.Delim:
		switch t.String() {
		case "{":
			term, err := decodeJSONObjectToTerm(decoder, env)
			if err != nil {
				return nil, err
			}
			if _, err = decoder.Token(); err != nil {
				return nil, err
			}
			return term, nil
		case "[":
			term, err := decodeJSONArrayToTerm(decoder, env)
			if err != nil {
				return nil, err
			}
			if _, err = decoder.Token(); err != nil {
				return nil, err
			}
			return term, nil
		}
	case string:
		return prolog.StringToAtom(t), nil
	case float64:
		return engine.NewFloatFromString(strconv.FormatFloat(t, 'f', -1, 64))
	case bool:
		return prolog.JSONBool(t), nil
	case nil:
		return prolog.JSONNull(), nil
	}

	return nil, jsonErrorToException(decoder.stream, fmt.Errorf("unexpected token: %v", t), env)
}

func decodeJSONArrayToTerm(decoder *textStreamDecoder, env *engine.Env) (engine.Term, error) {
	var terms []engine.Term
	for decoder.More() {
		value, err := decodeJSONToTerm(decoder, env)
		if err != nil {
			return nil, err
		}
		terms = append(terms, value)
	}

	return engine.List(terms...), nil
}

func decodeJSONObjectToTerm(decoder *textStreamDecoder, env *engine.Env) (engine.Term, error) {
	var terms []engine.Term
	for decoder.More() {
		keyToken, err := nextToken(decoder, env)
		if err != nil {
			return nil, err
		}
		key := keyToken.(string)
		value, err := decodeJSONToTerm(decoder, env)
		if err != nil {
			return nil, err
		}
		terms = append(terms, prolog.AtomKeyValue.Apply(prolog.StringToAtom(key), value))
	}

	return prolog.AtomJSON.Apply(engine.List(terms...)), nil
}

func jsonErrorToException(culprit engine.Term, err error, env *engine.Env) engine.Exception {
	if err, ok := lo.ErrorsAs[*json.SyntaxError](err); ok {
		return engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomMalformedJSON.Apply(engine.Integer(err.Offset))), env)
	}

	switch {
	case errors.Is(err, io.EOF):
		return engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomEOF), env)
	case err.Error() == errWrongIOMode.Error():
		return engine.PermissionError(operationInput, permissionTypeStream, culprit, env)
	case err.Error() == errWrongStreamType.Error():
		return engine.PermissionError(operationInput, permissionTypeTextStream, culprit, env)
	case err.Error() == errPastEndOfStream.Error():
		return engine.PermissionError(operationInput, permissionTypePastEndOfStream, culprit, env)
	}

	if err, ok := lo.ErrorsAs[*json.UnmarshalTypeError](err); ok {
		return engine.SyntaxError(
			AtomSyntaxErrorJSON.Apply(AtomMalformedJSON.Apply(engine.Integer(err.Offset), prolog.StringToAtom(err.Value))), env)
	}

	return prolog.WithError(
		engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomUnknown), env), err, env)
}

func nextToken(decoder *textStreamDecoder, env *engine.Env) (json.Token, error) {
	t, err := decoder.Token()
	if err != nil {
		return nil, jsonErrorToException(decoder.stream, err, env)
	}
	return t, nil
}

type textStreamDecoder struct {
	stream *engine.Stream
	*json.Decoder
}

func newTextStreamDecoder(stream *engine.Stream) *textStreamDecoder {
	decoder := &textStreamDecoder{
		stream: stream,
	}
	decoder.Decoder = json.NewDecoder(decoder)

	return decoder
}

func (s *textStreamDecoder) Read(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}

	r, size, err := s.stream.ReadRune()
	if err != nil {
		return 0, err
	}

	n = utf8.EncodeRune(p, r)
	if n < size {
		return n, io.ErrShortBuffer
	}

	return n, nil
}

type textStreamWriter struct {
	stream *engine.Stream
}

func newTextStreamWriter(stream *engine.Stream) *textStreamWriter {
	return &textStreamWriter{
		stream: stream,
	}
}

func (s *textStreamWriter) Write(p []byte) (n int, err error) {
	for len(p) > 0 {
		r, size := utf8.DecodeRune(p)
		if r == utf8.RuneError && size == 1 {
			return n, errInvalidUTF8
		}
		if _, err := s.stream.WriteRune(r); err != nil {
			return n, err
		}
		p = p[size:]
		n += size
	}
	return n, nil
}
