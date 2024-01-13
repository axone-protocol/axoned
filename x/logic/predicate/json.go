package predicate

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/logic/prolog"
)

// JSONProlog is a predicate that will unify a JSON string into prolog terms and vice versa.
//
// The signature is as follows:
//
//	json_prolog(?Json, ?Term) is det
//
// Where:
//   - Json is the string representation of the json
//   - Term is an Atom that would be unified by the JSON representation as Prolog terms.
//
// In addition, when passing Json and Term, this predicate return true if both result match.
//
// Examples:
//
//	# JSON conversion to Prolog.
//	- json_prolog('{"foo": "bar"}', json([foo-bar])).
func JSONProlog(vm *engine.VM, j, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var result engine.Term

		switch t1 := env.Resolve(j).(type) {
		case engine.Variable:
		case engine.Atom:
			terms, err := jsonStringToTerms(t1, env)
			if err != nil {
				return engine.Error(err)
			}
			result = terms
		default:
			return engine.Error(engine.TypeError(prolog.AtomTypeAtom, j, env))
		}

		switch t2 := env.Resolve(term).(type) {
		case engine.Variable:
			if result == nil {
				return engine.Error(engine.InstantiationError(env))
			}
			return engine.Unify(vm, term, result, cont, env)
		default:
			b, err := termsToJSON(t2, env)
			if err != nil {
				return engine.Error(err)
			}

			b, err = sdk.SortJSON(b)
			if err != nil {
				return engine.Error(
					prolog.WithError(
						engine.DomainError(prolog.ValidEncoding("json"), term, env), err, env))
			}
			var r engine.Term = engine.NewAtom(string(b))
			return engine.Unify(vm, j, r, cont, env)
		}
	})
}

func jsonStringToTerms(j engine.Atom, env *engine.Env) (engine.Term, error) {
	var values any
	decoder := json.NewDecoder(strings.NewReader(j.String()))
	decoder.UseNumber() // unmarshal a number into an interface{} as a Number instead of as a float64

	if err := decoder.Decode(&values); err != nil {
		return nil, prolog.WithError(
			engine.DomainError(prolog.ValidEncoding("json"), j, env), err, env)
	}

	term, err := jsonToTerms(values)
	if err != nil {
		return nil, prolog.WithError(
			engine.DomainError(prolog.ValidEncoding("json"), j, env), err, env)
	}

	return term, nil
}

func termsToJSON(term engine.Term, env *engine.Env) ([]byte, error) {
	asDomainError := func(bs []byte, err error) ([]byte, error) {
		if err != nil {
			return bs, prolog.WithError(
				engine.DomainError(prolog.ValidEncoding("json"), term, env), err, env)
		}
		return bs, err
	}
	switch t := term.(type) {
	case engine.Atom:
		return asDomainError(json.Marshal(t.String()))
	case engine.Integer:
		return asDomainError(json.Marshal(t))
	case engine.Compound:
		switch {
		case t.Functor() == prolog.AtomDot:
			iter, err := prolog.ListIterator(t, env)
			if err != nil {
				return nil, err
			}

			elements := make([]json.RawMessage, 0)
			for iter.Next() {
				element, err := termsToJSON(env.Resolve(iter.Current()), env)
				if err != nil {
					return nil, err
				}
				elements = append(elements, element)
			}
			return asDomainError(json.Marshal(elements))
		case t.Functor() == prolog.AtomJSON:
			terms, err := prolog.ExtractJSONTerm(t, env)
			if err != nil {
				return nil, err
			}

			attributes := make(map[string]json.RawMessage, len(terms))
			for key, term := range terms {
				raw, err := termsToJSON(env.Resolve(term), env)
				if err != nil {
					return nil, err
				}
				attributes[key] = raw
			}
			return asDomainError(json.Marshal(attributes))
		case prolog.JSONBool(true).Compare(t, env) == 0:
			return asDomainError(json.Marshal(true))
		case prolog.JSONBool(false).Compare(t, env) == 0:
			return asDomainError(json.Marshal(false))
		case prolog.JSONEmptyArray().Compare(t, env) == 0:
			return asDomainError(json.Marshal([]json.RawMessage{}))
		case prolog.JSONNull().Compare(t, env) == 0:
			return asDomainError(json.Marshal(nil))
		default:
			// no-op
		}
	default:
		// no-op
	}

	return nil, engine.TypeError(prolog.AtomTypeJSON, term, env)
}

func jsonToTerms(value any) (engine.Term, error) {
	switch v := value.(type) {
	case string:
		var r engine.Term = engine.NewAtom(v)
		return r, nil
	case json.Number:
		r, ok := math.NewIntFromString(string(v))
		if !ok {
			return nil, fmt.Errorf("could not convert number '%s' into integer term, decimal number is not handled yet", v)
		}
		if !r.IsInt64() {
			return nil, fmt.Errorf("could not convert number '%s' into integer term, overflow", v)
		}
		return engine.Integer(r.Int64()), nil
	case bool:
		return prolog.JSONBool(v), nil
	case nil:
		return prolog.JSONNull(), nil
	case map[string]any:
		keys := lo.Keys(v)
		sort.Strings(keys)

		attributes := make([]engine.Term, 0, len(v))
		for _, key := range keys {
			attributeValue, err := jsonToTerms(v[key])
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, prolog.AtomPair.Apply(engine.NewAtom(key), attributeValue))
		}
		return prolog.AtomJSON.Apply(engine.List(attributes...)), nil
	case []any:
		elements := make([]engine.Term, 0, len(v))
		if len(v) == 0 {
			return prolog.JSONEmptyArray(), nil
		}

		for _, element := range v {
			term, err := jsonToTerms(element)
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
