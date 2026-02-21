package prolog

import "github.com/axone-protocol/prolog/v3/engine"

var (
	// AtomAs is the term used to indicate the as encoding type option.
	AtomAs = engine.NewAtom("as")
	// AtomAt are terms with principal functor (@)/1 used to represent special values in json objects.
	AtomAt = engine.NewAtom("@")
	// AtomDIDComponents is a term which represents a DID as a compound term `did_components(Method, ID, Path, Query, Fragment)`.
	AtomDIDComponents = engine.NewAtom("did_components")
	// AtomDot is the term used to represent the dot in a list.
	AtomDot = engine.NewAtom(".")
	// AtomEmpty is the term used to represent empty.
	AtomEmpty = engine.NewAtom("")
	// AtomEmptyList is the term used to represent an empty list.
	AtomEmptyList = engine.NewAtom("[]")
	// AtomEncoding is the term used to indicate the encoding type option.
	AtomEncoding = engine.NewAtom("encoding")
	// AtomError is the term used to indicate the error.
	AtomError = engine.NewAtom("error")
	// AtomFalse is the term false.
	AtomFalse = engine.NewAtom("false")
	// AtomFragment is the term used to indicate the fragment component.
	AtomFragment = engine.NewAtom("fragment")
	// AtomHex is the term used to indicate the hexadecimal encoding type option.
	AtomHex = engine.NewAtom("hex")
	// AtomJSON are terms with principal functor json/1 used to represent json objects.
	AtomJSON = engine.NewAtom("json")
	// AtomNull is the term null.
	AtomNull = engine.NewAtom("null")
	// AtomOctet is the term used to indicate the byte encoding type option.
	AtomOctet = engine.NewAtom("octet")
	// AtomPadding is the term used to indicate the padding encoding type option.
	AtomPadding = engine.NewAtom("padding")
	// AtomCharset is the term used to indicate the charset encoding type option.
	AtomCharset = engine.NewAtom("charset")
	// AtomPair are terms with principal functor (-)/2.
	// For example, the term -(A, B) denotes the pair of elements A and B.
	AtomPair = engine.NewAtom("-")
	// AtomKeyValue are terms with principal functor (=)/2.
	// For example, the term =(A, B) denotes the mapping of key A with value B.
	AtomKeyValue = engine.NewAtom("=")
	// AtomPath is the term used to indicate the path component.
	AtomPath = engine.NewAtom("path")
	// AtomQuoted is the term used to indicate the quoted write_term/3 option.
	AtomQuoted = engine.NewAtom("quoted")
	// AtomQueryValue is the term used to indicate the query value component.
	AtomQueryValue = engine.NewAtom("query_value")
	// AtomSegment is the term used to indicate the segment component.
	AtomSegment = engine.NewAtom("segment")
	// AtomText is the term used to indicate the atom text.
	AtomText = engine.NewAtom("text")
	// AtomBoolean is the term used to indicate the atom boolean.
	AtomBoolean = engine.NewAtom("boolean")
	// AtomTrue is the term true.
	AtomTrue = engine.NewAtom("true")
	// AtomUtf8 is the term used to indicate the UTF-8 encoding type option.
	AtomUtf8 = engine.NewAtom("utf8")
)
