package prolog

import (
	"strings"
	"unicode/utf8"

	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"
)

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

// IsGround returns true if the given term holds no free variables.
func IsGround(term engine.Term, env *engine.Env) bool {
	switch term := env.Resolve(term).(type) {
	case engine.Variable:
		return false
	case engine.Compound:
		for i := 0; i < term.Arity(); i++ {
			if !IsGround(term.Arg(i), env) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

func AreGround(terms []engine.Term, env *engine.Env) bool {
	return lo.EveryBy(terms, func(t engine.Term) bool {
		return IsGround(t, env)
	})
}

// AssertIsGround resolves a term and returns it if it is ground.
// If the term is not ground, the function returns nil and the instantiation error.
func AssertIsGround(env *engine.Env, t engine.Term) (engine.Term, error) {
	if IsGround(t, env) {
		return t, nil
	}
	return nil, engine.InstantiationError(env)
}

// AssertAtom resolves a term and attempts to convert it into an engine.Atom if possible.
// If conversion fails, the function returns the empty atom and the error.
func AssertAtom(env *engine.Env, t engine.Term) (engine.Atom, error) {
	_, err := AssertIsGround(env, t)
	if err != nil {
		return AtomEmpty, err
	}
	if t, ok := t.(engine.Atom); ok {
		return t, nil
	}
	return AtomEmpty, engine.TypeError(AtomTypeAtom, t, env)
}

// AssertCharacterCode resolves a term and attempts to convert it into a rune if possible.
// If conversion fails, the function returns the zero value and the error.
func AssertCharacterCode(env *engine.Env, t engine.Term) (rune, error) {
	_, err := AssertIsGround(env, t)
	if err != nil {
		return 0, err
	}

	if t, ok := t.(engine.Integer); ok {
		if t >= 0 && t <= utf8.MaxRune {
			return rune(t), nil
		}
	}

	return 0, engine.TypeError(AtomTypeCharacterCode, t, env)
}

// AssertCharacter resolves a term and attempts to convert it into an engine.Atom if possible.
// If conversion fails, the function returns the empty atom and the error.
func AssertCharacter(env *engine.Env, t engine.Term) (rune, error) {
	_, err := AssertIsGround(env, t)
	if err != nil {
		return utf8.RuneError, err
	}
	if t, ok := t.(engine.Atom); ok {
		runes := []rune(t.String())
		if len(runes) == 1 {
			return runes[0], nil
		}
	}
	return utf8.RuneError, engine.TypeError(AtomTypeCharacter, t, env)
}

// AssertByte resolves a term and attempts to convert it into a byte if possible.
// If conversion fails, the function returns the zero value and the error.
func AssertByte(env *engine.Env, t engine.Term) (byte, error) {
	_, err := AssertIsGround(env, t)
	if err != nil {
		return 0, err
	}
	if t, ok := t.(engine.Integer); ok {
		if t >= 0 && t <= 255 {
			return byte(t), nil
		}
	}
	return 0, engine.TypeError(AtomTypeByte, t, env)
}

// AssertList resolves a term as a list and returns it as a engine.Compound.
// If conversion fails, the function returns nil and the error.
func AssertList(env *engine.Env, t engine.Term) (engine.Term, error) {
	_, err := AssertIsGround(env, t)
	if err != nil {
		return nil, err
	}
	if IsList(t) {
		return t, nil
	}

	return nil, engine.TypeError(AtomTypeList, t, env)
}

// AssertPair resolves a term as a pair and returns the pair components.
// If conversion fails, the function returns nil and the error.
func AssertPair(env *engine.Env, t engine.Term) (engine.Term, engine.Term, error) {
	_, err := AssertIsGround(env, t)
	if err != nil {
		return nil, nil, err
	}
	if t, ok := t.(engine.Compound); ok && t.Functor() == AtomPair && t.Arity() == 2 {
		return t.Arg(0), t.Arg(1), nil
	}

	return nil, nil, engine.TypeError(AtomTypePair, t, env)
}
