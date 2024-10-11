package predicate

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

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
//   - A JSON object is mapped to a Prolog term json(NameValueList), where NameValueList is a list of Name-Value pairs.
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
//	- json_prolog('{"foo": "bar"}', json([foo-bar])).
func JSONProlog(_ *engine.VM, j, p engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	forwardConverter := func(in []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		payload, err := prolog.TextTermToString(in[0], env)
		if err != nil {
			return nil, err
		}

		decoder := json.NewDecoder(strings.NewReader(payload))
		term, err := decodeJSONToTerm(decoder, env)
		if err != nil {
			return nil, err
		}
		if _, err := decoder.Token(); !errors.Is(err, io.EOF) {
			return nil, engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomMalformedJSON.Apply(engine.Integer(decoder.InputOffset()))), env)
		}

		return []engine.Term{term}, nil
	}
	backwardConverter := func(in []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		var buf bytes.Buffer
		err := encodeTermToJSON(in[0], &buf, env)
		if err != nil {
			return nil, err
		}

		return []engine.Term{prolog.BytesToAtom(buf.Bytes())}, nil
	}
	return prolog.UnifyFunctionalPredicate(
		[]engine.Term{j}, []engine.Term{p}, prolog.AtomEmpty, forwardConverter, backwardConverter, cont, env)
}

func encodeTermToJSON(term engine.Term, buf *bytes.Buffer, env *engine.Env) (err error) {
	marshalToBuffer := func(data any) error {
		bs, err := json.Marshal(data)
		if err != nil {
			return prologErrorToException(term, err, env)
		}
		buf.Write(bs)

		return nil
	}

	switch t := term.(type) {
	case engine.Atom:
		if term == prolog.AtomEmptyList {
			buf.Write([]byte("[]"))
		} else {
			return marshalToBuffer(t.String())
		}
	case engine.Integer:
		return marshalToBuffer(t)
	case engine.Float:
		float, err := strconv.ParseFloat(t.String(), 64)
		if err != nil {
			return prologErrorToException(t, err, env)
		}
		return marshalToBuffer(float)
	case engine.Compound:
		return encodeCompoundToJSON(t, buf, env)
	default:
		return engine.TypeError(prolog.AtomTypeJSON, term, env)
	}

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
		k, v, err := prolog.AssertPair(t, env)
		if err != nil {
			return err
		}
		key, err := prolog.AssertAtom(k, env)
		if err != nil {
			return err
		}
		bs, err := json.Marshal(key.String())
		if err != nil {
			return prologErrorToException(t, err, env)
		}
		buf.Write(bs)
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

func jsonErrorToException(err error, env *engine.Env) engine.Exception {
	if err, ok := lo.ErrorsAs[*json.SyntaxError](err); ok {
		return engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomMalformedJSON.Apply(engine.Integer(err.Offset))), env)
	}

	if errors.Is(err, io.EOF) {
		return engine.SyntaxError(AtomSyntaxErrorJSON.Apply(AtomEOF), env)
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

func nextToken(decoder *json.Decoder, env *engine.Env) (json.Token, error) {
	t, err := decoder.Token()
	if err != nil {
		return nil, jsonErrorToException(err, env)
	}
	return t, nil
}

func decodeJSONToTerm(decoder *json.Decoder, env *engine.Env) (engine.Term, error) {
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

	return nil, jsonErrorToException(fmt.Errorf("unexpected token: %v", t), env)
}

func decodeJSONArrayToTerm(decoder *json.Decoder, env *engine.Env) (engine.Term, error) {
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

func decodeJSONObjectToTerm(decoder *json.Decoder, env *engine.Env) (engine.Term, error) {
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
		terms = append(terms, prolog.AtomPair.Apply(prolog.StringToAtom(key), value))
	}

	return prolog.AtomJSON.Apply(engine.List(terms...)), nil
}
