package predicate

import (
	"encoding/json"
	"fmt"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
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
	json.Unmarshal([]byte(j), &values)

	return jsonToTerms(values)
}

func jsonToTerms(value any) (engine.Term, error) {
	switch v := value.(type) {
	case string:
		return util.StringToTerm(v), nil
	default:
		return nil, fmt.Errorf("could not convert %s (%T) to a prolog term", v, v)
	}
}
