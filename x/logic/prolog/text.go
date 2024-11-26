package prolog

import (
	"strings"
	"unicode/utf8"

	"github.com/axone-protocol/prolog/v2/engine"
)

// AtomToString try to convert a given atom to a string.
func AtomToString(term engine.Term, env *engine.Env) (string, error) {
	atom, err := AssertAtom(term, env)
	if err != nil {
		return "", err
	}
	return atom.String(), nil
}

// listTermToString try to convert a given list to a string using the provided
// converter function.
// The converter function is called for each element of the list and is expected
// to return a rune.
func listTermToString(
	term engine.Term,
	converter func(engine.Term, *engine.Env) (rune, error),
	env *engine.Env,
) (string, error) {
	iter, err := ListIterator(term, env)
	if err != nil {
		return "", err
	}
	var sb strings.Builder

	for iter.Next() {
		r, err := converter(iter.Current(), env)
		if err != nil {
			return sb.String(), err
		}
		sb.WriteRune(r)
	}
	return sb.String(), nil
}

// CharacterListTermToString try to convert a given list of characters to a string.
// Characters is a list of atoms, each representing a single character.
func CharacterListTermToString(term engine.Term, env *engine.Env) (string, error) {
	return listTermToString(term, AssertCharacter, env)
}

// CharacterCodeListTermToString try to convert a given list of character codes to a string.
// The character codes must be between 0 and 0x10ffff (i.e. a Rune).
func CharacterCodeListTermToString(term engine.Term, env *engine.Env) (string, error) {
	return listTermToString(term, AssertCharacterCode, env)
}

// OctetListTermToString try to convert a given list of bytes to a string.
// It's the same as CharacterCodeListTermToString, but expects the list to contain bytes.
// It's equivalent to the prolog encoding 'octet'.
func OctetListTermToString(term engine.Term, env *engine.Env) (string, error) {
	return listTermToString(term, func(term engine.Term, env *engine.Env) (rune, error) {
		b, err := AssertByte(term, env)
		if err != nil {
			return utf8.RuneError, err
		}
		return rune(b), nil
	}, env)
}

// TextTermToString try to convert a given Text term to a string.
// Text is an instantiated term which represents text as: an atom, a list of character codes, or list of characters.
func TextTermToString(term engine.Term, env *engine.Env) (string, error) {
	switch v := env.Resolve(term).(type) {
	case engine.Atom:
		return AtomToString(v, env)
	case engine.Compound:
		if IsList(v, env) {
			head := ListHead(v, env)
			if head == nil {
				return "", nil
			}

			switch head.(type) {
			case engine.Atom:
				return CharacterListTermToString(v, env)
			case engine.Integer:
				return CharacterCodeListTermToString(v, env)
			default:
				return "", engine.TypeError(AtomTypeCharacterCode, v, env)
			}
		}
	}
	return "", engine.TypeError(AtomTypeText, term, env)
}

// StringToAtom converts a string to an atom.
func StringToAtom(s string) engine.Atom {
	return engine.NewAtom(s)
}

// StringToCharacterListTerm converts a string to a term representing a list of characters.
func StringToCharacterListTerm(s string) engine.Term {
	terms := make([]engine.Term, 0, utf8.RuneCountInString(s))
	for _, r := range s {
		terms = append(terms, engine.NewAtom(string(r)))
	}

	return engine.List(terms...)
}

// StringToCharacterCodeListTerm converts a string to a term representing a list of character codes.
func StringToCharacterCodeListTerm(s string) engine.Term {
	terms := make([]engine.Term, 0, utf8.RuneCountInString(s))
	for _, r := range s {
		terms = append(terms, engine.Integer(r))
	}

	return engine.List(terms...)
}

// StringToOctetListTerm converts a string (utf8) to a term representing a list of bytes.
// This is the same as StringToCharacterCodeListTerm, but it returns an error when a rune is greater than 0xff.
// This is equivalent to the prolog encoding 'octet'.
func StringToOctetListTerm(s string, env *engine.Env) (engine.Term, error) {
	terms := make([]engine.Term, 0, utf8.RuneCountInString(s))
	for _, r := range s {
		if r > 0xff {
			return nil, engine.TypeError(AtomTypeByte, engine.Integer(r), env)
		}
		terms = append(terms, engine.Integer(r))
	}

	return engine.List(terms...), nil
}

// StringToByteListTerm converts a string (utf8) to a term representing a list of bytes.
// This is equivalent to the prolog encoding 'text'.
func StringToByteListTerm(s string) engine.Term {
	return BytesToByteListTerm([]byte(s))
}
