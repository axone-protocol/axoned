package codec

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/axone-protocol/prolog/v3/engine"
	"github.com/samber/lo"

	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
)

var (
	atomSyntaxError     = engine.NewAtom("syntax_error")
	atomJSON            = engine.NewAtom(codecNameJSON)
	atomAt              = engine.NewAtom("@")
	atomNull            = engine.NewAtom("null")
	atomTrue            = engine.NewAtom("true")
	atomFalse           = engine.NewAtom("false")
	atomPrologSyntax    = engine.NewAtom("prolog")
	atomMalformedJSON   = engine.NewAtom("malformed_json")
	atomMalformedTerm   = engine.NewAtom("malformed_term")
	atomEOF             = engine.NewAtom("eof")
	atomUnknown         = engine.NewAtom("unknown")
	atomValidJSONNumber = engine.NewAtom("json_number")
	atomSystemError     = engine.NewAtom("system_error")
)

var (
	jsonNullTerm  = atomAt.Apply(atomNull)
	jsonTrueTerm  = atomAt.Apply(atomTrue)
	jsonFalseTerm = atomAt.Apply(atomFalse)
)

type jsonCodec struct{}

func init() {
	Register(&jsonCodec{})
}

func (c *jsonCodec) Name() string {
	return codecNameJSON
}

func (c *jsonCodec) Decode(payload []byte) engine.Term {
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoded, err := decodeJSONToTerm(decoder)
	if err != nil {
		return prolog.AtomError.Apply(jsonErrorTerm(err))
	}
	if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
		return prolog.AtomError.Apply(
			atomSyntaxError.Apply(atomJSON.Apply(atomMalformedJSON.Apply(engine.Integer(decoder.InputOffset())))),
		)
	}

	return atomOK.Apply(decoded)
}

func (c *jsonCodec) Encode(payload []byte) engine.Term {
	term, err := parseCodecTerm(payload)
	if err != nil {
		return prolog.AtomError.Apply(atomSyntaxError.Apply(atomPrologSyntax.Apply(atomMalformedTerm)))
	}

	var buf bytes.Buffer
	if err := encodeTermToJSON(term, &buf, engine.NewEnv()); err != nil {
		return prolog.AtomError.Apply(exceptionFormal(err))
	}

	return atomOK.Apply(engine.NewAtom(buf.String()))
}

func encodeTermToJSON(term engine.Term, writer io.Writer, env *engine.Env) error {
	switch t := env.Resolve(term).(type) {
	case engine.Atom:
		if t == prolog.AtomEmptyList {
			_, err := writer.Write([]byte("[]"))
			return err
		}
		return marshalToJSONStream(t.String(), term, writer, env)
	case engine.Integer:
		return marshalToJSONStream(t, term, writer, env)
	case engine.Float:
		float, err := strconv.ParseFloat(t.String(), 64)
		if err != nil {
			return engine.DomainError(atomValidJSONNumber, t, env)
		}
		return marshalToJSONStream(float, term, writer, env)
	case engine.Compound:
		return encodeCompoundToJSON(t, writer, env)
	case engine.Variable:
		return engine.InstantiationError(env)
	default:
		return engine.TypeError(atomJSON, term, env)
	}
}

func encodeCompoundToJSON(term engine.Compound, writer io.Writer, env *engine.Env) error {
	switch {
	case term.Functor() == prolog.AtomDot:
		return encodeArrayToJSON(term, writer, env)
	case term.Functor() == atomJSON:
		return encodeObjectToJSON(term, writer, env)
	case jsonBool(true).Compare(term, env) == 0:
		_, err := writer.Write([]byte("true"))
		return err
	case jsonBool(false).Compare(term, env) == 0:
		_, err := writer.Write([]byte("false"))
		return err
	case jsonNull().Compare(term, env) == 0:
		_, err := writer.Write([]byte("null"))
		return err
	default:
		return engine.TypeError(atomJSON, term, env)
	}
}

func encodeObjectToJSON(term engine.Compound, writer io.Writer, env *engine.Env) error {
	if _, err := assertJSON(term, env); err != nil {
		return err
	}
	if _, err := writer.Write([]byte("{")); err != nil {
		return err
	}
	if err := prolog.ForEach(term.Arg(0), env, func(t engine.Term, hasNext bool) error {
		k, v, err := prolog.AssertKeyValue(t, env)
		if err != nil {
			return err
		}
		if err := marshalToJSONStream(k.String(), term, writer, env); err != nil {
			return err
		}
		if _, err := writer.Write([]byte(":")); err != nil {
			return err
		}
		if err := encodeTermToJSON(v, writer, env); err != nil {
			return err
		}
		if hasNext {
			if _, err := writer.Write([]byte(",")); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	_, err := writer.Write([]byte("}"))
	return err
}

func encodeArrayToJSON(term engine.Compound, writer io.Writer, env *engine.Env) error {
	if _, err := writer.Write([]byte("[")); err != nil {
		return err
	}
	if err := prolog.ForEach(term, env, func(t engine.Term, hasNext bool) error {
		if err := encodeTermToJSON(t, writer, env); err != nil {
			return err
		}
		if hasNext {
			if _, err := writer.Write([]byte(",")); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	_, err := writer.Write([]byte("]"))
	return err
}

func marshalToJSONStream(data any, term engine.Term, writer io.Writer, env *engine.Env) error {
	bs, err := json.Marshal(data)
	if err != nil {
		if _, ok := lo.ErrorsAs[*strconv.NumError](err); ok {
			return engine.DomainError(atomValidJSONNumber, term, env)
		}
		return err
	}
	_, err = writer.Write(bs)
	return err
}

func decodeJSONToTerm(decoder *json.Decoder) (engine.Term, error) {
	token, err := decoder.Token()
	if errors.Is(err, io.EOF) {
		return jsonNull(), nil
	}
	if err != nil {
		return nil, err
	}

	switch token := token.(type) {
	case json.Delim:
		switch token.String() {
		case "{":
			term, err := decodeJSONObjectToTerm(decoder)
			if err != nil {
				return nil, err
			}
			if _, err = decoder.Token(); err != nil {
				return nil, err
			}
			return term, nil
		case "[":
			term, err := decodeJSONArrayToTerm(decoder)
			if err != nil {
				return nil, err
			}
			if _, err = decoder.Token(); err != nil {
				return nil, err
			}
			return term, nil
		}
	case string:
		return prolog.StringToAtom(token), nil
	case float64:
		return engine.NewFloatFromString(strconv.FormatFloat(token, 'f', -1, 64))
	case bool:
		return jsonBool(token), nil
	case nil:
		return jsonNull(), nil
	}

	return nil, fmt.Errorf("unexpected token: %v", token)
}

func decodeJSONArrayToTerm(decoder *json.Decoder) (engine.Term, error) {
	var terms []engine.Term
	for decoder.More() {
		value, err := decodeJSONToTerm(decoder)
		if err != nil {
			return nil, err
		}
		terms = append(terms, value)
	}

	return engine.List(terms...), nil
}

func decodeJSONObjectToTerm(decoder *json.Decoder) (engine.Term, error) {
	var terms []engine.Term
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		key, ok := keyToken.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected object key token: %v", keyToken)
		}
		value, err := decodeJSONToTerm(decoder)
		if err != nil {
			return nil, err
		}
		terms = append(terms, prolog.AtomKeyValue.Apply(prolog.StringToAtom(key), value))
	}

	return atomJSON.Apply(engine.List(terms...)), nil
}

func jsonErrorTerm(err error) engine.Term {
	if err, ok := lo.ErrorsAs[*json.SyntaxError](err); ok {
		return atomSyntaxError.Apply(atomJSON.Apply(atomMalformedJSON.Apply(engine.Integer(err.Offset))))
	}
	if err, ok := lo.ErrorsAs[*json.UnmarshalTypeError](err); ok {
		return atomSyntaxError.Apply(
			atomJSON.Apply(atomMalformedJSON.Apply(engine.Integer(err.Offset), prolog.StringToAtom(err.Value))),
		)
	}
	if errors.Is(err, io.EOF) {
		return atomSyntaxError.Apply(atomJSON.Apply(atomEOF))
	}
	return atomSyntaxError.Apply(atomJSON.Apply(atomUnknown))
}

func jsonNull() engine.Term {
	return jsonNullTerm
}

func jsonBool(b bool) engine.Term {
	if b {
		return jsonTrueTerm
	}
	return jsonFalseTerm
}

func assertJSON(term engine.Term, env *engine.Env) (engine.Compound, error) {
	if compound, ok := env.Resolve(term).(engine.Compound); ok {
		if compound.Functor() == atomJSON && compound.Arity() == 1 {
			return compound, nil
		}
	}

	return nil, engine.TypeError(atomJSON, term, env)
}

func exceptionFormal(err error) engine.Term {
	var exception engine.Exception
	if errors.As(err, &exception) {
		if compound, ok := exception.Term().(engine.Compound); ok &&
			compound.Functor() == prolog.AtomError && compound.Arity() >= 1 {
			return compound.Arg(0)
		}
	}
	return atomSystemError
}
