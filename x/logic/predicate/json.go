package predicate

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/axone-protocol/prolog/engine"
	"github.com/samber/lo"

	"github.com/axone-protocol/axoned/v10/x/logic/prolog"
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
func JSONProlog(_ *engine.VM, json, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	forwardConverter := func(json []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		term, err := decodeJSONToTerm(json[0], env)
		if err != nil {
			return nil, err
		}
		return []engine.Term{term}, nil
	}
	backwardConverter := func(term []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		b, err := encodeTermToJSON(term[0], env)
		if err != nil {
			return nil, err
		}
		return []engine.Term{prolog.BytesToAtom(b)}, nil
	}
	return prolog.UnifyFunctionalPredicate(
		[]engine.Term{json}, []engine.Term{term}, prolog.AtomEmpty, forwardConverter, backwardConverter, cont, env)
}

// decodeJSONToTerm decode a JSON, given as a prolog text, into a prolog term.
func decodeJSONToTerm(j engine.Term, env *engine.Env) (engine.Term, error) {
	payload, err := prolog.TextTermToString(j, env)
	if err != nil {
		return nil, err
	}

	var values any
	decoder := json.NewDecoder(strings.NewReader(payload))
	decoder.UseNumber() // unmarshal a number into an interface{} as a Number instead of as a float64

	if err := decoder.Decode(&values); err != nil {
		return nil, prolog.WithError(
			engine.DomainError(prolog.ValidEncoding("json"), j, env), err, env)
	}

	term, err := jsonToTerm(values)
	if err != nil {
		return nil, prolog.WithError(
			engine.DomainError(prolog.ValidEncoding("json"), j, env), err, env)
	}

	return term, nil
}

// encodeTermToJSON converts a Prolog term to a JSON byte array.
func encodeTermToJSON(term engine.Term, env *engine.Env) ([]byte, error) {
	bs, err := termToJSON(term, env)

	var exception engine.Exception
	if err != nil && !errors.As(err, &exception) {
		return nil, prolog.WithError(engine.DomainError(prolog.ValidEncoding("json"), term, env), err, env)
	}

	return bs, err
}

func termToJSON(term engine.Term, env *engine.Env) ([]byte, error) {
	switch t := term.(type) {
	case engine.Atom:
		return json.Marshal(t.String())
	case engine.Integer:
		return json.Marshal(t)
	case engine.Float:
		float, err := strconv.ParseFloat(t.String(), 64)
		if err != nil {
			return nil, err
		}

		return json.Marshal(float)
	case engine.Compound:
		return compoundToJSON(t, env)
	}

	return nil, engine.TypeError(prolog.AtomTypeJSON, term, env)
}

func compoundToJSON(term engine.Compound, env *engine.Env) ([]byte, error) {
	switch {
	case term.Functor() == prolog.AtomDot:
		iter, err := prolog.ListIterator(term, env)
		if err != nil {
			return nil, err
		}

		elements := make([]json.RawMessage, 0)
		for iter.Next() {
			element, err := termToJSON(iter.Current(), env)
			if err != nil {
				return nil, err
			}
			elements = append(elements, element)
		}
		return json.Marshal(elements)
	case term.Functor() == prolog.AtomJSON:
		terms, err := prolog.ExtractJSONTerm(term, env)
		if err != nil {
			return nil, err
		}

		attributes := make(map[string]json.RawMessage, len(terms))
		for key, term := range terms {
			raw, err := termToJSON(term, env)
			if err != nil {
				return nil, err
			}
			attributes[key] = raw
		}
		return json.Marshal(attributes)
	case prolog.JSONBool(true).Compare(term, env) == 0:
		return json.Marshal(true)
	case prolog.JSONBool(false).Compare(term, env) == 0:
		return json.Marshal(false)
	case prolog.JSONEmptyArray().Compare(term, env) == 0:
		return json.Marshal([]json.RawMessage{})
	case prolog.JSONNull().Compare(term, env) == 0:
		return json.Marshal(nil)
	}

	return nil, engine.TypeError(prolog.AtomTypeJSON, term, env)
}

func jsonToTerm(value any) (engine.Term, error) {
	switch v := value.(type) {
	case string:
		return prolog.StringToAtom(v), nil
	case json.Number:
		return engine.NewFloatFromString(v.String())
	case bool:
		return prolog.JSONBool(v), nil
	case nil:
		return prolog.JSONNull(), nil
	case map[string]any:
		keys := lo.Keys(v)
		sort.Strings(keys)

		attributes := make([]engine.Term, 0, len(v))
		for _, key := range keys {
			attributeValue, err := jsonToTerm(v[key])
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, prolog.AtomPair.Apply(prolog.StringToAtom(key), attributeValue))
		}
		return prolog.AtomJSON.Apply(engine.List(attributes...)), nil
	case []any:
		if len(v) == 0 {
			return prolog.JSONEmptyArray(), nil
		}
		elements := make([]engine.Term, 0, len(v))
		for _, element := range v {
			term, err := jsonToTerm(element)
			if err != nil {
				return nil, err
			}
			elements = append(elements, term)
		}

		return engine.List(elements...), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}
