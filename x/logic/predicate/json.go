package predicate

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"cosmossdk.io/math"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
	"github.com/samber/lo"
)

// AtomJSON is a term which represents a json as a compound term `json([Pair])`.
var AtomJSON = engine.NewAtom("json")

// JsonProlog is a predicate that will unify a JSON string into prolog terms and vice versa.
//
// json_prolog(?Json, ?Term) is det
// TODO:
func JsonProlog(vm *engine.VM, j, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	switch t1 := env.Resolve(j).(type) {
	case engine.Variable:
	case engine.Atom:
		terms, err := jsonStringToTerms(t1.String())
		if err != nil {
			return engine.Error(fmt.Errorf("json_prolog/2: %w", err))
		}

		return engine.Unify(vm, term, terms, cont, env)
	default:
		return engine.Error(fmt.Errorf("did_components/2: cannot unify json with %T", t1))
	}

	switch env.Resolve(term).(type) {
	default:
		return engine.Error(fmt.Errorf("json_prolog/2: not implemented"))
	}
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
