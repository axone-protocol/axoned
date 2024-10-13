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
)

var (
	operationInput                = engine.NewAtom("input")
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
//	json_read(+Stream, -Term) is det
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
	case prolog.IsGround(p, env):
		var buf bytes.Buffer
		err := encodeTermToJSON(p, &buf, env)
		if err != nil {
			return engine.Error(err)
		}
		return engine.Unify(vm, prolog.BytesToAtom(buf.Bytes()), j, cont, env)
	default:
		return engine.Error(engine.InstantiationError(env))
	}
}

func encodeTermToJSON(term engine.Term, buf *bytes.Buffer, env *engine.Env) (err error) {
	switch t := term.(type) {
	case engine.Atom:
		if term == prolog.AtomEmptyList {
			buf.Write([]byte("[]"))
		} else {
			return marshalToBuffer(t.String(), term, buf, env)
		}
	case engine.Integer:
		return marshalToBuffer(t, term, buf, env)
	case engine.Float:
		float, err := strconv.ParseFloat(t.String(), 64)
		if err != nil {
			return prologErrorToException(t, err, env)
		}
		return marshalToBuffer(float, term, buf, env)
	case engine.Compound:
		return encodeCompoundToJSON(t, buf, env)
	default:
		return engine.TypeError(prolog.AtomTypeJSON, term, env)
	}

	return nil
}

func marshalToBuffer(data any, term engine.Term, buf *bytes.Buffer, env *engine.Env) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return prologErrorToException(term, err, env)
	}
	buf.Write(bs)

	return nil
}

func encodeCompoundToJSON(term engine.Compound, buf *bytes.Buffer, env *engine.Env) error {
	switch {
	case term.Functor() == prolog.AtomDot:
		return encodeArrayToJSON(term, buf, env)
	case term.Functor() == prolog.AtomJSON:
		return encodeObjectToJSON(term, buf, env)
	case prolog.JSONBool(true).Compare(term, env) == 0:
		buf.Write([]byte("true"))
	case prolog.JSONBool(false).Compare(term, env) == 0:
		buf.Write([]byte("false"))
	case prolog.JSONNull().Compare(term, env) == 0:
		buf.Write([]byte("null"))
	default:
		return engine.TypeError(prolog.AtomTypeJSON, term, env)
	}

	return nil
}

func encodeObjectToJSON(term engine.Compound, buf *bytes.Buffer, env *engine.Env) error {
	if _, err := prolog.AssertJSON(term, env); err != nil {
		return err
	}
	buf.Write([]byte("{"))
	if err := prolog.ForEach(term.Arg(0), env, func(t engine.Term, hasNext bool) error {
		k, v, err := prolog.AssertKeyValue(t, env)
		if err != nil {
			return err
		}
		if err := marshalToBuffer(k.String(), term, buf, env); err != nil {
			return err
		}
		buf.Write([]byte(":"))
		if err := encodeTermToJSON(v, buf, env); err != nil {
			return err
		}

		if hasNext {
			buf.Write([]byte(","))
		}
		return nil
	}); err != nil {
		return err
	}
	buf.Write([]byte("}"))
	return nil
}

func encodeArrayToJSON(term engine.Compound, buf *bytes.Buffer, env *engine.Env) error {
	buf.Write([]byte("["))
	if err := prolog.ForEach(term, env, func(t engine.Term, hasNext bool) error {
		err := encodeTermToJSON(t, buf, env)
		if err != nil {
			return err
		}

		if hasNext {
			buf.Write([]byte(","))
		}

		return nil
	}); err != nil {
		return err
	}
	buf.Write([]byte("]"))

	return nil
}

func jsonErrorToException(stream engine.Term, err error, env *engine.Env) engine.Exception {
	if err, ok := lo.ErrorsAs[*json.SyntaxError](err); ok {
		return engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomMalformedJSON.Apply(engine.Integer(err.Offset))), env)
	}

	switch {
	case errors.Is(err, io.EOF):
		return engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomEOF), env)
	case err.Error() == errWrongIOMode.Error():
		return engine.PermissionError(operationInput, permissionTypeStream, stream, env)
	case err.Error() == errWrongStreamType.Error():
		return engine.PermissionError(operationInput, permissionTypeTextStream, stream, env)
	case err.Error() == errPastEndOfStream.Error():
		return engine.PermissionError(operationInput, permissionTypePastEndOfStream, stream, env)
	}

	if err, ok := lo.ErrorsAs[*json.UnmarshalTypeError](err); ok {
		return engine.SyntaxError(
			AtomSyntaxErrorJSON.Apply(AtomMalformedJSON.Apply(engine.Integer(err.Offset), prolog.StringToAtom(err.Value))), env)
	}

	return prolog.WithError(
		engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomUnknown), env), err, env)
}

func prologErrorToException(culprit engine.Term, err error, env *engine.Env) engine.Exception {
	if _, ok := lo.ErrorsAs[*strconv.NumError](err); ok {
		return engine.DomainError(AtomValidJSONNumber, culprit, env)
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
