package prolog

import "github.com/ichiban/prolog/engine"

// Common atoms.
var (
	AtomAs         = engine.NewAtom("as")       // AtomAs is the term used to indicate the as encoding type option.
	AtomAt         = engine.NewAtom("@")        // AtomAt are terms with principal functor (@)/1. It is used to represent special values in json objects.
	AtomAtom       = engine.NewAtom("atom")     // AtomAtom is the term used to indicate the atom atom,
	AtomCharset    = engine.NewAtom("charset")  // AtomCharset is the term used to indicate the charset encoding type option.
	AtomCompound   = engine.NewAtom("compound") // AtomCompound is the term used to indicate the atom compound,
	AtomDot        = engine.NewAtom(".")        // AtomDot is the term used to represent the dot in a list.
	AtomEmpty      = engine.NewAtom("")         // AtomEmpty is the term used to represent empty.
	AtomEmptyArray = engine.NewAtom("[]")       // AtomEmptyArray is the term [].
	AtomEmptyList  = engine.NewAtom("[]")       // AtomEmptyList is the term used to represent an empty list.
	AtomEncoding   = engine.NewAtom("encoding") // AtomEncoding is the term used to indicate the encoding type option.
	AtomError      = engine.NewAtom("error")    // AtomError is the term used to indicate the error.
	AtomFalse      = engine.NewAtom("false")    // AtomFalse is the term false.
	AtomHex        = engine.NewAtom("hex")      // AtomHex is the term used to indicate the hexadecimal encoding type option.
	AtomJSON       = engine.NewAtom("json")     // AtomJSON are terms with principal functor json/1. // It is used to represent json objects.
	AtomList       = engine.NewAtom("list")     // AtomList is the term used to indicate the atom list,
	AtomNull       = engine.NewAtom("null")     // AtomNull is the term null.
	AtomOctet      = engine.NewAtom("octet")    // AtomOctet is the term used to indicate the byte encoding type option.
	AtomPadding    = engine.NewAtom("padding")  // AtomPadding is the term used to indicate the padding encoding type option.
	AtomPair       = engine.NewAtom("-")        // AtomPair are terms with principal functor (-)/2. For example, the term -(A, B) denotes the pair of elements A and B.
	AtomTrue       = engine.NewAtom("true")     // AtomTrue is the term true.
	AtomUtf8       = engine.NewAtom("utf8")     // AtomUtf8 is the term used to indicate the UTF-8 encoding type option.
)
