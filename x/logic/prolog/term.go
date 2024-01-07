package prolog

import (
	"encoding/hex"
	"unicode/utf8"

	"github.com/ichiban/prolog/engine"
)

// Tuple is a predicate which unifies the given term with a tuple of the given arity.
func Tuple(args ...engine.Term) engine.Term {
	return engine.Atom(0).Apply(args...)
}

// ListOfIntegers converts a list of integers to a term.
func ListOfIntegers(args ...int) engine.Term {
	terms := make([]engine.Term, 0, len(args))
	for _, arg := range args {
		terms = append(terms, engine.Integer(arg))
	}
	return engine.List(terms...)
}

// StringToTerm converts a string to a term.
// TODO: this function should be removed.
func StringToTerm(s string) engine.Term {
	return engine.NewAtom(s)
}

// StringToStringTerm converts a string to a term representing a list of characters.
func StringToStringTerm(s string) engine.Term {
	terms := make([]engine.Term, 0, utf8.RuneCountInString(s))
	for _, c := range s {
		terms = append(terms, engine.NewAtom(string(c)))
	}

	return engine.List(terms...)
}

// BytesToCodepointListTerm try to convert a given golang []byte into a list of codepoints.
func BytesToCodepointListTerm(in []byte, encoding engine.Atom, env *engine.Env) (engine.Term, error) {
	out, err := Decode(in, encoding, env)
	if err != nil {
		return nil, err
	}

	terms := make([]engine.Term, 0, len(out))
	for _, b := range out {
		terms = append(terms, engine.Integer(b))
	}
	return engine.List(terms...), nil
}

// BytesToCodepointListTermWithDefault is like the BytesToCodepointListTerm function but with a default encoding.
// This function panics if the conversion fails, which can't happen with the default encoding.
func BytesToCodepointListTermWithDefault(in []byte, env *engine.Env) engine.Term {
	term, err := BytesToCodepointListTerm(in, AtomEmpty, env)
	if err != nil {
		panic(err)
	}
	return term
}

// BytesToAtomListTerm try to convert a given golang []byte into a list of atoms, one for each character.
func BytesToAtomListTerm(in []byte, encoding engine.Atom, env *engine.Env) (engine.Term, error) {
	out, err := Decode(in, encoding, env)
	if err != nil {
		return nil, err
	}
	str := string(out)
	terms := make([]engine.Term, 0, len(str))
	for _, c := range str {
		terms = append(terms, engine.NewAtom(string(c)))
	}
	return engine.List(terms...), nil
}

// StringTermToBytes try to convert a given string into native golang []byte.
// String is an instantiated term which represents text as an atom, string, list of character codes or list or characters.
// Encoding is the supported encoding type:
//   - empty encoding or 'text', return the original bytes without modification.
//   - 'octet', decode the bytes as unicode code points and return the resulting bytes. If a code point is greater than
//     0xff, an error is returned.
//   - any other encoding label, convert the bytes to the specified encoding.
//
// The mapping from encoding labels to encodings is defined at https://encoding.spec.whatwg.org/.
func StringTermToBytes(str engine.Term, encoding engine.Atom, env *engine.Env) (bs []byte, err error) {
	v := env.Resolve(str)
	switch v := v.(type) {
	case engine.Atom:
		if bs, err = Encode([]byte(v.String()), encoding, env); err != nil {
			return nil, err
		}
		return bs, nil
	case engine.Compound:
		if IsList(v) {
			head := ListHead(v, env)
			if head == nil {
				return make([]byte, 0), nil
			}

			switch head.(type) {
			case engine.Atom:
				if bs, err = characterListToBytes(v, env); err != nil {
					return bs, err
				}
			case engine.Integer:
				if bs, err = characterCodeListToBytes(v, env); err != nil {
					return bs, err
				}
			default:
				return nil, engine.TypeError(AtomCharacterCode, v, env)
			}
			return Encode(bs, encoding, env)
		}
		return nil, engine.TypeError(AtomCharacterCode, str, env)
	default:
		return nil, engine.TypeError(AtomText, str, env)
	}
}

func characterListToBytes(str engine.Compound, env *engine.Env) ([]byte, error) {
	iter := engine.ListIterator{List: str, Env: env}
	bs := make([]byte, 0)

	for iter.Next() {
		e, err := AssertCharacter(env, iter.Current())
		if err != nil {
			return bs, err
		}
		rs := []rune(e.String())
		if len(rs) != 1 {
			return bs, engine.DomainError(ValidCharacterCode(e.String()), str, env)
		}

		bs = append(bs, []byte(e.String())...)
	}
	return bs, nil
}

func characterCodeListToBytes(str engine.Compound, env *engine.Env) ([]byte, error) {
	iter := engine.ListIterator{List: str, Env: env}
	bs := make([]byte, 0)

	for iter.Next() {
		e, err := AssertCharacterCode(env, iter.Current())
		if err != nil {
			return nil, err
		}
		if e < 0 || e > 255 {
			return nil, engine.DomainError(ValidByte(int64(e)), str, env)
		}

		bs = append(bs, byte(e))
	}
	return bs, nil
}

// TermHexToBytes try to convert an hexadecimal encoded atom to native golang []byte.
func TermHexToBytes(term engine.Term, env *engine.Env) ([]byte, error) {
	v, err := AssertAtom(env, term)
	if err != nil {
		return nil, err
	}

	src := []byte(v.String())
	result := make([]byte, hex.DecodedLen(len(src)))
	_, err = hex.Decode(result, src)
	if err != nil {
		err = engine.DomainError(ValidEncoding("hex", err), term, env)
	}
	return result, err
}
