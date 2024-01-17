package prolog

import (
	"github.com/ichiban/prolog/engine"
)

// JSONNull returns the compound term @(null).
// It is used to represent the null value in json objects.
func JSONNull() engine.Term {
	return AtomAt.Apply(AtomNull)
}

// JSONBool returns the compound term @(true) if b is true, otherwise @(false).
func JSONBool(b bool) engine.Term {
	if b {
		return AtomAt.Apply(AtomTrue)
	}

	return AtomAt.Apply(AtomFalse)
}

// JSONEmptyArray returns is the compound term @([]).
// It is used to represent the empty array in json objects.
func JSONEmptyArray() engine.Term {
	return AtomAt.Apply(AtomEmptyArray)
}

// ExtractJSONTerm is a utility function that would extract all attribute of a JSON object
// that is represented in prolog with the `json` atom.
//
// This function will ensure the json atom follow our json object representation in prolog.
//
// A JSON object is represented like this :
//
// ```
// json([foo-bar])
// ```
//
// That give a JSON object: `{"foo": "bar"}`
// Returns the map of all attributes with its term value.
func ExtractJSONTerm(term engine.Compound, env *engine.Env) (map[string]engine.Term, error) {
	if term.Functor() != AtomJSON || term.Arity() != 1 {
		return nil, engine.TypeError(AtomTypeJSON, term, env)
	}

	iter, err := ListIterator(term.Arg(0), env)
	if err != nil {
		return nil, err
	}
	terms := make(map[string]engine.Term, 0)
	for iter.Next() {
		current := iter.Current()
		pair, ok := current.(engine.Compound)
		if !ok || pair.Functor() != AtomPair || pair.Arity() != 2 {
			return nil, engine.TypeError(AtomTypePair, current, env)
		}

		key, ok := pair.Arg(0).(engine.Atom)
		if !ok {
			return nil, engine.TypeError(AtomTypeAtom, pair.Arg(0), env)
		}
		terms[key.String()] = pair.Arg(1)
	}
	return terms, nil
}
