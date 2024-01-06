package prolog

import (
	"fmt"

	"github.com/ichiban/prolog/engine"
)

// JsonNull returns the compound term @(null).
// It is used to represent the null value in json objects.
func JsonNull() engine.Term {
	return AtomAt.Apply(AtomNull)
}

// JsonBool returns the compound term @(true) if b is true, otherwise @(false).
func JsonBool(b bool) engine.Term {
	if b {
		return AtomAt.Apply(AtomTrue)
	}

	return AtomAt.Apply(AtomFalse)
}

// JsonEmptyArray returns is the compound term @([]).
// It is used to represent the empty array in json objects.
func JsonEmptyArray() engine.Term {
	return AtomAt.Apply(AtomEmptyArray)
}

// ExtractJSONTerm is an utility function that would extract all attribute of a JSON object
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
	if term.Functor() != AtomJSON {
		return nil, fmt.Errorf("invalid functor %s. Expected %s", term.Functor().String(), AtomJSON.String())
	} else if term.Arity() != 1 {
		return nil, fmt.Errorf("invalid compound arity : %d but expected %d", term.Arity(), 1)
	}

	list := term.Arg(0)
	switch l := env.Resolve(list).(type) {
	case engine.Compound:
		iter := engine.ListIterator{
			List: l,
			Env:  env,
		}
		terms := make(map[string]engine.Term, 0)
		for iter.Next() {
			pair, ok := env.Resolve(iter.Current()).(engine.Compound)
			if !ok || pair.Functor() != AtomPair || pair.Arity() != 2 {
				return nil, fmt.Errorf("json attributes should be a pair")
			}

			key, ok := env.Resolve(pair.Arg(0)).(engine.Atom)
			if !ok {
				return nil, fmt.Errorf("first pair arg should be an atom")
			}
			terms[key.String()] = pair.Arg(1)
		}
		return terms, nil
	default:
		return nil, fmt.Errorf("json compound should contains one list, give %T", l)
	}
}
