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

	// AtomUtf8 is the term used to indicate the UTF-8 encoding type option.
	AtomUtf8 = engine.NewAtom("utf8")

	// AtomHex is the term used to indicate the hexadecimal encoding type option.
	AtomHex = engine.NewAtom("hex")

	// AtomOctet is the term used to indicate the byte encoding type option.
	AtomOctet = engine.NewAtom("octet")

	// AtomCharset is the term used to indicate the charset encoding type option.
	AtomCharset = engine.NewAtom("charset")

	// AtomPadding is the term used to indicate the padding encoding type option.
	AtomPadding = engine.NewAtom("padding")

	// AtomAs is the term used to indicate the as encoding type option.
	AtomAs = engine.NewAtom("as")
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
