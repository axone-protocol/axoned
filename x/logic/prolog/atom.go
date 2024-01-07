package prolog

import "github.com/ichiban/prolog/engine"

var (
	// AtomAs is the term used to indicate the as encoding type option.
	AtomAs = engine.NewAtom("as")
	// AtomAt are terms with principal functor (@)/1 used to represent special values in json objects.
	AtomAt = engine.NewAtom("@")
	// AtomAtom is the term used to indicate the atom atom.
	AtomAtom = engine.NewAtom("atom")
	// AtomCharacter is the term used to indicate the character type.
	AtomCharacter = engine.NewAtom("character")
	// AtomCharacterCode is the term used to indicate the character code type.
	AtomCharacterCode = engine.NewAtom("character_code")
	// AtomCharset is the term used to indicate the charset encoding type option.
	AtomCharset = engine.NewAtom("charset")
	// AtomCompound is the term used to indicate the atom compound.
	AtomCompound = engine.NewAtom("compound")
	// AtomDot is the term used to represent the dot in a list.
	AtomDot = engine.NewAtom(".")
	// AtomEmpty is the term used to represent empty.
	AtomEmpty = engine.NewAtom("")
	// AtomEmptyArray is the term [].
	AtomEmptyArray = engine.NewAtom("[]")
	// AtomEmptyList is the term used to represent an empty list.
	AtomEmptyList = engine.NewAtom("[]")
	// AtomEncoding is the term used to indicate the encoding type option.
	AtomEncoding = engine.NewAtom("encoding")
	// AtomError is the term used to indicate the error.
	AtomError = engine.NewAtom("error")
	// AtomFalse is the term false.
	AtomFalse = engine.NewAtom("false")
	// AtomHex is the term used to indicate the hexadecimal encoding type option.
	AtomHex = engine.NewAtom("hex")
	// AtomJSON are terms with principal functor json/1 used to represent json objects.
	AtomJSON = engine.NewAtom("json")
	// AtomList is the term used to indicate the atom list.
	AtomList = engine.NewAtom("list")
	// AtomNull is the term null.
	AtomNull = engine.NewAtom("null")
	// AtomOctet is the term used to indicate the byte encoding type option.
	AtomOctet = engine.NewAtom("octet")
	// AtomPadding is the term used to indicate the padding encoding type option.
	AtomPadding = engine.NewAtom("padding")
	// AtomPair are terms with principal functor (-)/2.
	// For example, the term -(A, B) denotes the pair of elements A and B.
	AtomPair = engine.NewAtom("-")
	// AtomText is the term used to indicate the atom text.
	AtomText = engine.NewAtom("text")
	// AtomTrue is the term true.
	AtomTrue = engine.NewAtom("true")
	// AtomUtf8 is the term used to indicate the UTF-8 encoding type option.
	AtomUtf8 = engine.NewAtom("utf8")
)
