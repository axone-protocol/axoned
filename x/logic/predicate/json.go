package predicate

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
	"github.com/samber/lo"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AtomJSON is a term which represents a json as a compound term `json([Pair])`.
var AtomJSON = engine.NewAtom("json")

// JSONProlog is a predicate that will unify a JSON string into prolog terms and vice versa.
//
// json_prolog(?Json, ?Term) is det
//
// Where
//   - `Json` is the string representation of the json
//   - `Term` is an Atom that would be unified by the JSON representation as Prolog terms.
//
// In addition, when passing Json and Term, this predicate return true if both result match.
//
// Example:
//
// # JSON conversion to Prolog.
// - json_prolog('{"foo": "bar"}', json([foo-bar])).
func JSONProlog(vm *engine.VM, j, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var result engine.Term

		switch t1 := env.Resolve(j).(type) {
		case engine.Variable:
		case engine.Atom:
			terms, err := jsonStringToTerms(t1.String())
			if err != nil {
				return engine.Error(fmt.Errorf("json_prolog/2: %w", err))
			}
			result = terms
		default:
			return engine.Error(fmt.Errorf("json_prolog/2: cannot unify json with %T", t1))
		}

		switch t2 := env.Resolve(term).(type) {
		case engine.Variable:
			if result == nil {
				return engine.Error(fmt.Errorf("json_prolog/2: could not unify two variable"))
			}
			return engine.Unify(vm, term, result, cont, env)
		default:
			b, err := termsToJSON(t2, env)
			if err != nil {
				return engine.Error(fmt.Errorf("json_prolog/2: %w", err))
			}

			b, err = sdk.SortJSON(b)
			if err != nil {
				return engine.Error(fmt.Errorf("json_prolog/2: %w", err))
			}
			return engine.Unify(vm, j, util.StringToTerm(string(b)), cont, env)
		}
	})
}

func jsonStringToTerms(j string) (engine.Term, error) {
	var values any
	decoder := json.NewDecoder(strings.NewReader(j))
	decoder.UseNumber() // unmarshal a number into an interface{} as a Number instead of as a float64

	if err := decoder.Decode(&values); err != nil {
		return nil, err
	}

	return jsonToTerms(values)
}

func termsToJSON(term engine.Term, env *engine.Env) ([]byte, error) {
	switch t := term.(type) {
	case engine.Atom:
		return json.Marshal(t.String())
	case engine.Integer:
		return json.Marshal(t)
	case engine.Compound:
		switch t.Functor().String() {
		case ".": // Represent an engine.List
			if t.Arity() != 2 {
				return nil, fmt.Errorf("wrong term arity for array, give %d, expected %d", t.Arity(), 2)
			}

			iter := engine.ListIterator{List: t, Env: env}

			elements := make([]json.RawMessage, 0)
			for iter.Next() {
				element, err := termsToJSON(env.Resolve(iter.Current()), env)
				if err != nil {
					return nil, err
				}
				elements = append(elements, element)
			}
			return json.Marshal(elements)
		case AtomJSON.String():
			terms, err := ExtractJSONTerm(t, env)
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
			return json.Marshal(attributes)
		}

		switch {
		case AtomBool(true).Compare(t, env) == 0:
			return json.Marshal(true)
		case AtomBool(false).Compare(t, env) == 0:
			return json.Marshal(false)
		case AtomNull.Compare(t, env) == 0:
			return json.Marshal(nil)
		}

		return nil, fmt.Errorf("invalid functor %s", t.Functor())
	default:
		return nil, fmt.Errorf("could not convert %s {%T} to json", t, t)
	}
}

func jsonToTerms(value any) (engine.Term, error) {
	switch v := value.(type) {
	case string:
		return util.StringToTerm(v), nil
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
		return AtomBool(v), nil
	case nil:
		return AtomNull, nil
	case map[string]any:
		keys := lo.Keys(v)
		sort.Strings(keys)

		attributes := make([]engine.Term, 0, len(v))
		for _, key := range keys {
			attributeValue, err := jsonToTerms(v[key])
			if err != nil {
				return nil, err
			}
			attributes = append(attributes, AtomPair.Apply(engine.NewAtom(key), attributeValue))
		}
		return AtomJSON.Apply(engine.List(attributes...)), nil
	case []any:
		elements := make([]engine.Term, 0, len(v))
		for _, element := range v {
			term, err := jsonToTerms(element)
			if err != nil {
				return nil, err
			}
			elements = append(elements, term)
		}
		return engine.List(elements...), nil
	default:
		return nil, fmt.Errorf("could not convert %s (%T) to a prolog term", v, v)
	}
}
