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

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v10/x/logic/prolog"
)

// JSONProlog is a predicate that unifies a JSON into a prolog term and vice versa.
//
// The signature is as follows:
//
//	json_prolog(?Json, ?Term) is det
//
// Where:
//   - Json is the textual representation of the json, as an atom, a list of character codes, or list of characters.
//   - Term is a term that represents the JSON in the prolog world.
//
// # Examples:
//
//	# JSON conversion to Prolog.
//	- json_prolog('{"foo": "bar"}', json([foo-bar])).
func JSONProlog(vm *engine.VM, j, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	var result engine.Term

	switch t1 := env.Resolve(j).(type) {
	case engine.Variable:
	default:
		terms, err := decodeJSONToTerm(t1, env)
		if err != nil {
			return engine.Error(err)
		}
		result = terms
	}

	switch t2 := env.Resolve(term).(type) {
	case engine.Variable:
		if result == nil {
			return engine.Error(engine.InstantiationError(env))
		}
		return engine.Unify(vm, term, result, cont, env)
	default:
		b, err := encodeTermToJSON(t2, env)
		if err != nil {
			return engine.Error(err)
		}

		b, err = sdk.SortJSON(b)
		if err != nil {
			return engine.Error(
				prolog.WithError(
					engine.DomainError(prolog.ValidEncoding("json"), term, env), err, env))
		}
		var r engine.Term = prolog.BytesToAtom(b)
		return engine.Unify(vm, j, r, cont, env)
	}
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
