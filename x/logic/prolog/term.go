package prolog

import (
	"encoding/hex"
	"fmt"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
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
func StringToTerm(s string) engine.Term {
	return engine.NewAtom(s)
}

// BytesToCodepointListTerm try to convert a given golang []byte into a list of codepoints.
func BytesToCodepointListTerm(in []byte, encoding string) (engine.Term, error) {
	out, err := util.Decode(in, encoding)
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
func BytesToCodepointListTermWithDefault(in []byte) engine.Term {
	term, err := BytesToCodepointListTerm(in, "")
	if err != nil {
		panic(err)
	}
	return term
}

// BytesToAtomListTerm try to convert a given golang []byte into a list of atoms, one for each character.
func BytesToAtomListTerm(in []byte, encoding string) (engine.Term, error) {
	out, err := util.Decode(in, encoding)
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
func StringTermToBytes(str engine.Term, encoding string, env *engine.Env) (bs []byte, err error) {
	v := env.Resolve(str)
	switch v := v.(type) {
	case engine.Atom:
		if bs, err = util.Encode([]byte(v.String()), encoding); err != nil {
			return nil, EncodingError(encoding, err, env)
		}
		return
	case engine.Compound:
		if IsList(v) {
			iter := engine.ListIterator{List: v, Env: env}
			bs := make([]byte, 0)
			index := 0

			for iter.Next() {
				term := env.Resolve(iter.Current())
				index++

				switch t := term.(type) {
				case engine.Integer:
					if t >= 0 && t <= 255 {
						bs = append(bs, byte(t))
					} else {
						return nil, fmt.Errorf("invalid integer value '%d' in list at position %d: out of byte range (0-255)", t, index)
					}
				case engine.Atom:
					rs := []rune(t.String())
					if len(rs) != 1 {
						return nil, fmt.Errorf("invalid character_code '%s' value in list at position %d: should be a single character",
							t.String(), index)
					}

					bs = append(bs, []byte(t.String())...)
				default:
					return nil, fmt.Errorf("invalid term type in list at position %d: %T, only character_code or integer allowed", index, term)
				}
			}
			return util.Encode(bs, encoding)
		}
		return nil, fmt.Errorf("invalid compound term: expected a list of character_code or integer")
	default:
		return nil, fmt.Errorf("term should be a List, given %T", str)
	}
}

// TermHexToBytes try to convert an hexadecimal encoded atom to native golang []byte.
func TermHexToBytes(term engine.Term, env *engine.Env) ([]byte, error) {
	v := env.Resolve(term)
	switch v := v.(type) {
	case engine.Atom:
		src := []byte(v.String())
		result := make([]byte, hex.DecodedLen(len(src)))
		_, err := hex.Decode(result, src)
		return result, err
	default:
		return nil, fmt.Errorf("invalid term: expected a hexadecimal encoded atom, given %T", term)
	}
}
