package predicate

import (
	"github.com/ichiban/prolog/engine"
)

// AtomPair are terms with principal functor (-)/2.
// For example, the term -(A, B) denotes the pair of elements A and B.
var AtomPair = engine.NewAtom("-")

// AtomJSON are terms with principal functor json/1.
// It is used to represent json objects.
var AtomJSON = engine.NewAtom("json")

// AtomAt are terms with principal functor (@)/1.
// It is used to represent special values in json objects.
var AtomAt = engine.NewAtom("@")

// AtomTrue is the term true.
var AtomTrue = engine.NewAtom("true")

// AtomFalse is the term false.
var AtomFalse = engine.NewAtom("false")

// AtomNull is the term null.
var AtomNull = engine.NewAtom("null")

// MakeNull returns the compound term @(null).
// It is used to represent the null value in json objects.
func MakeNull() engine.Term {
	return AtomAt.Apply(AtomNull)
}

// MakeBool returns the compound term @(true) if b is true, otherwise @(false).
func MakeBool(b bool) engine.Term {
	if b {
		return AtomAt.Apply(AtomTrue)
	}

	return AtomAt.Apply(AtomFalse)
}
