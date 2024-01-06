package util

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"
	"golang.org/x/net/html/charset"
)

var (
	// AtomDot is the term used to represent the dot in a list.
	AtomDot = engine.NewAtom(".")

	// AtomEmpty is the term used to represent empty.
	AtomEmpty = engine.NewAtom("")

	// AtomEmptyList is the term used to represent an empty list.
	AtomEmptyList = engine.NewAtom("[]")
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
	out, err := decode(in, encoding)
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
	out, err := decode(in, encoding)
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
func StringTermToBytes(str engine.Term, encoding string, env *engine.Env) ([]byte, error) {
	v := env.Resolve(str)
	switch v := v.(type) {
	case engine.Atom:
		return encode([]byte(v.String()), encoding)
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
			return encode(bs, encoding)
		}
		return nil, fmt.Errorf("invalid compound term: expected a list of character_code or integer")
	default:
		return nil, fmt.Errorf("term should be a List, given %T", str)
	}
}

// encode converts a byte slice to a specified encoding.
//
// In case of:
//   - empty encoding label or 'text', return the original bytes without modification.
//   - 'octet', decode the bytes as unicode code points and return the resulting bytes. If a code point is greater than
//     0xff, an error is returned.
//   - any other encoding label, convert the bytes to the specified encoding.
func encode(bs []byte, label string) ([]byte, error) {
	switch label {
	case "", "text":
		return bs, nil
	case "octet":
		result := make([]byte, 0, len(bs))
		for i := 0; i < len(bs); {
			runeValue, width := utf8.DecodeRune(bs[i:])

			if runeValue > 0xff {
				return nil, fmt.Errorf("cannot convert character '%c' to %s", runeValue, label)
			}
			result = append(result, byte(runeValue))
			i += width
		}
		return result, nil
	default:
		encoding, _ := charset.Lookup(label)
		if encoding == nil {
			return nil, fmt.Errorf("invalid encoding: %s", label)
		}
		return encoding.NewEncoder().Bytes(bs)
	}
}

// decode converts a byte slice from a specified encoding.
// decode function is the reverse of encode function.
func decode(bs []byte, label string) ([]byte, error) {
	switch label {
	case "", "text":
		return bs, nil
	case "octet":
		var buffer bytes.Buffer
		for _, b := range bs {
			buffer.WriteRune(rune(b))
		}
		return buffer.Bytes(), nil
	default:
		encoding, _ := charset.Lookup(label)
		if encoding == nil {
			return nil, fmt.Errorf("invalid encoding: %s", label)
		}
		return encoding.NewDecoder().Bytes(bs)
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

// PredicateMatches returns a function that matches the given predicate against the given other predicate.
// If the other predicate contains a slash, it is matched as is. Otherwise, the other predicate is matched against the
// first part of the given predicate.
// For example:
//   - matchPredicate("foo/0")("foo/0") -> true
//   - matchPredicate("foo/0")("foo/1") -> false
//   - matchPredicate("foo/0")("foo") -> true
//   - matchPredicate("foo/0")("bar") -> false
//
// The function is curried, and is a binary relation that is reflexive, associative (but not commutative).
func PredicateMatches(this string) func(string) bool {
	return func(that string) bool {
		if strings.Contains(that, "/") {
			return this == that
		}
		return strings.Split(this, "/")[0] == that
	}
}

// IsList returns true if the given term is a list.
func IsList(term engine.Term) bool {
	switch v := term.(type) {
	case engine.Compound:
		return v.Functor() == AtomDot && v.Arity() == 2
	case engine.Atom:
		return v == AtomEmptyList
	}

	return false
}

// IsEmptyList returns true if the given term is an empty list.
func IsEmptyList(term engine.Term) bool {
	if v, ok := term.(engine.Atom); ok {
		return v == AtomEmptyList
	}
	return false
}

// IsVariable returns true if the given term is a variable.
func IsVariable(term engine.Term) bool {
	_, ok := term.(engine.Variable)
	return ok
}

// IsAtom returns true if the given term is an atom.
func IsAtom(term engine.Term) bool {
	_, ok := term.(engine.Atom)
	return ok
}

// IsCompound returns true if the given term is a compound.
func IsCompound(term engine.Term) bool {
	_, ok := term.(engine.Compound)
	return ok
}

// IsFullyInstantiated returns true if the given term is fully instantiated.
func IsFullyInstantiated(term engine.Term, env *engine.Env) bool {
	switch term := env.Resolve(term).(type) {
	case engine.Variable:
		return false
	case engine.Compound:
		for i := 0; i < term.Arity(); i++ {
			if !IsFullyInstantiated(term.Arg(i), env) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

func AreFullyInstantiated(terms []engine.Term, env *engine.Env) bool {
	_, ok := lo.Find(terms, func(t engine.Term) bool {
		return IsFullyInstantiated(t, env)
	})

	return ok
}

// AssertAtom resolves a term and attempts to convert it into an engine.Atom if possible.
// If conversion fails, the function returns the empty atom and the error.
func AssertAtom(env *engine.Env, t engine.Term) (engine.Atom, error) {
	switch t := env.Resolve(t).(type) {
	case engine.Atom:
		return t, nil
	default:
		return AtomEmpty,
			fmt.Errorf("invalid term '%s' - expected engine.Atom but got %T", t, t)
	}
}

// GetOption returns the value of the first option with the given name in the given options.
// An option is a compound with the given name as functor and one argument which is
// a term, for instance `opt(v)`.
// The options are either a list of options or an option.
// If no option is found nil is returned.
func GetOption(name engine.Atom, options engine.Term, env *engine.Env) (engine.Term, error) {
	extractOption := func(term engine.Term) (engine.Term, error) {
		switch v := term.(type) {
		case engine.Compound:
			if v.Functor() == name {
				if v.Arity() != 1 {
					return nil, fmt.Errorf("invalid arity for compound '%s': %d but expected 1", name, v.Arity())
				}

				return v.Arg(0), nil
			}
			return nil, nil
		case nil:
			return nil, nil
		default:
			return nil, fmt.Errorf("invalid term '%s' - expected engine.Compound but got %T", term, v)
		}
	}

	resolvedTerm := env.Resolve(options)

	if IsEmptyList(resolvedTerm) {
		return nil, nil
	}

	if IsList(resolvedTerm) {
		iter := engine.ListIterator{List: resolvedTerm, Env: env}

		for iter.Next() {
			opt := env.Resolve(iter.Current())

			term, err := extractOption(opt)
			if err != nil {
				return nil, err
			}

			if term != nil {
				return term, nil
			}
		}
		return nil, nil
	}

	return extractOption(resolvedTerm)
}

// GetOptionWithDefault returns the value of the first option with the given name in the given options, or the given
// default value if no option is found.
func GetOptionWithDefault(
	name engine.Atom, options engine.Term, defaultValue engine.Term, env *engine.Env,
) (engine.Term, error) {
	if term, err := GetOption(name, options, env); err != nil {
		return nil, err
	} else if term != nil {
		return term, nil
	}
	return defaultValue, nil
}

// GetOptionAsAtomWithDefault is a helper function that returns the value of the first option with the given name in the
// given options.
func GetOptionAsAtomWithDefault(
	name engine.Atom, options engine.Term, defaultValue engine.Term, env *engine.Env,
) (engine.Atom, error) {
	term, err := GetOptionWithDefault(name, options, defaultValue, env)
	if err != nil {
		return AtomEmpty, err
	}
	atom, err := AssertAtom(env, term)
	if err != nil {
		return AtomEmpty, err
	}

	return atom, nil
}

// ConvertFunc is a function mapping a domain which is a list of terms with a codomain which is a set of terms.
// Domains and co-domains can have different cardinalities.
// options is a list of options that can be used to parameterize the conversion.
// All the terms provided are fully instantiated (i.e. no variables).
type ConvertFunc func(value []engine.Term, options engine.Term, env *engine.Env) ([]engine.Term, error)

// UnifyFunctional is a generic unification which unifies a set of input terms with a set of output terms, using the
// given conversion functions maintaining the function's relationship.
//
// The aim of this function is to simplify the implementation of a wide range of predicates which are essentially
// functional, like hash functions, encoding functions, etc.
//
// The semantic of the unification is as follows:
//  1. first all the variables are resolved
//  2. if there's variables in the input and the output,
//     the conversion is not possible and a not sufficiently instantiated error is returned.
//  3. if there's no variables in the input,
//     then the conversion is attempted from the input to the output and the result is unified with the output.
//  4. if there's no variables in the output,
//     then the conversion is attempted from the output to the input and the result is unified with the input.
//
// The following table summarizes the behavior, where:
// - fi = fully instantiated (i.e. no variables)
// - !fi = not fully instantiated (i.e. at least one variable)
//
// | input | output | result                               |
// |-------|--------|--------------------------------------|
// | !fi   | !fi    | error: not sufficiently instantiated |
// |  fi   | !fi    | unify(forward(input), output)        |
// |  fi   |  fi    | unify(forward(input), output)        |
// | !fi   |  fi    | unify(input,backward(output))        |
//
// Conversion functions may produce an error in scenarios where the conversion is unsuccessful or infeasible due to
// the inherent characteristics of the function's relationship, such as the absence of a one-to-one correspondence
// (e.g. hash functions).
func UnifyFunctional(
	vm *engine.VM,
	in,
	out []engine.Term,
	options engine.Term,
	forwardConverter ConvertFunc,
	backwardConverter ConvertFunc,
	cont engine.Cont,
	env *engine.Env,
) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		isInFI, isOutFi := AreFullyInstantiated(in, env), AreFullyInstantiated(out, env)
		if !isInFI && !isOutFi {
			return engine.Error(engine.InstantiationError(env))
		}

		var err error
		from, to := in, out
		if isInFI {
			from, err = forwardConverter(in, options, env)
			if err != nil {
				return engine.Error(err)
			}
		} else {
			to, err = backwardConverter(out, options, env)
			if err != nil {
				return engine.Error(err)
			}
		}
		return engine.Unify(
			vm,
			Tuple(from...),
			Tuple(to...),
			cont,
			env,
		)
	})
}
