package prolog

import (
	"strings"

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
	if t, ok := env.Resolve(t).(engine.Atom); ok {
		return t, nil
	}
	return AtomEmpty, engine.TypeError(AtomAtom, t, env)
}

// AssertList resolves a term as a list and returns it as a engine.Compound.
// If conversion fails, the function returns nil and the error.
func AssertList(env *engine.Env, t engine.Term) (engine.Compound, error) {
	if t, ok := env.Resolve(t).(engine.Compound); ok && IsList(t) {
		return t, nil
	}

	return nil, engine.TypeError(AtomList, t, env)
}
