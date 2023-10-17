package predicate

import (
	"github.com/ichiban/prolog/engine"
)

var (
	// AtomPair are terms with principal functor (-)/2.
	// For example, the term -(A, B) denotes the pair of elements A and B.
	AtomPair = engine.NewAtom("-")

	// AtomJSON are terms with principal functor json/1.
	// It is used to represent json objects.
	AtomJSON = engine.NewAtom("json")

	// AtomAt are terms with principal functor (@)/1.
	// It is used to represent special values in json objects.
	AtomAt = engine.NewAtom("@")

	// AtomTrue is the term true.
	AtomTrue = engine.NewAtom("true")

	// AtomFalse is the term false.
	AtomFalse = engine.NewAtom("false")

	// AtomEmptyArray is the term [].
	AtomEmptyArray = engine.NewAtom("[]")

	// AtomNull is the term null.
	AtomNull = engine.NewAtom("null")

	// AtomEncoding is the term used to indicate the encoding type option.
	AtomEncoding = engine.NewAtom("encoding")
)

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

// MakeEmptyArray returns is the compound term @([]).
// It is used to represent the empty array in json objects.
func MakeEmptyArray() engine.Term {
	return AtomAt.Apply(AtomEmptyArray)
}
