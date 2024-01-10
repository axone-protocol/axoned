package prolog

import "github.com/ichiban/prolog/engine"

var (
	AtomDomainError        = engine.NewAtom("domain_error")         // AtomDomainError is the atom domain_error.
	AtomValidByte          = engine.NewAtom("valid_byte")           // AtomValidByte is the atom valid_byte.
	AtomValidCharset       = engine.NewAtom("valid_charset")        // AtomValidCharset is the atom valid_charset.
	AtomValidCharacterCode = engine.NewAtom("valid_character_code") // AtomValidCharacterCode is the atom valid_character_code.
	AtomValidEncoding      = engine.NewAtom("valid_encoding")       // AtomValidEncoding is the atom valid_encoding.
	AtomValidHexDigit      = engine.NewAtom("valid_hex_digit")      // AtomValidHexDigit is the atom valid_hex_digit.
)

func ValidCharset() engine.Term {
	return AtomValidCharset
}

func ValidEncoding(encoding string, cause error) engine.Term {
	return AtomValidEncoding.Apply(engine.NewAtom(encoding), StringToStringTerm(cause.Error()))
}

func ValidByte(v int64) engine.Term {
	return AtomValidByte.Apply(engine.Integer(v))
}

func ValidCharacterCode(c string) engine.Term {
	return AtomValidCharacterCode.Apply(engine.NewAtom(c))
}

func ValidHexDigit(d string) engine.Term {
	return AtomValidHexDigit.Apply(engine.NewAtom(d))
}
