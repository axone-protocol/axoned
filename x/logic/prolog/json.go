package prolog

import (
	"github.com/axone-protocol/prolog/v2/engine"
)

var (
	nullTerm  = AtomAt.Apply(AtomNull)
	trueTerm  = AtomAt.Apply(AtomTrue)
	falseTerm = AtomAt.Apply(AtomFalse)
)

// JSONNull returns the compound term @(null).
// It is used to represent the null value in json objects.
func JSONNull() engine.Term {
	return nullTerm
}

// JSONBool returns the compound term @(true) if b is true, otherwise @(false).
func JSONBool(b bool) engine.Term {
	if b {
		return trueTerm
	}

	return falseTerm
}

// AssertJSON resolves a term as a JSON object and returns it as engine.Compound.
// If conversion fails, the function returns nil and the error.
func AssertJSON(term engine.Term, env *engine.Env) (engine.Compound, error) {
	if compound, ok := env.Resolve(term).(engine.Compound); ok {
		if compound.Functor() == AtomJSON && compound.Arity() == 1 {
			return compound, nil
		}
	}

	return nil, engine.TypeError(AtomTypeJSON, term, env)
}
