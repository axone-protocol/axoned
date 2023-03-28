package util

import (
	"strings"

	"github.com/ichiban/prolog/engine"
)

// StringToTerm converts a string to a term.
// If the string is empty, it returns a variable.
func StringToTerm(s string) engine.Term {
	if s == "" {
		return engine.NewVariable()
	}

	return engine.NewAtom(s)
}

// Resolve resolves a term and returns the resolved term and a boolean indicating whether the term is instantiated.
func Resolve(env *engine.Env, t engine.Term) (engine.Atom, bool) {
	switch t := env.Resolve(t).(type) {
	case engine.Atom:
		return t, true
	default:
		return engine.NewAtom(""), false
	}
}

// PredicateEq returns a function that matches the given predicate against the given other predicate.
// If the other predicate contains a slash, it is matched as is. Otherwise, the other predicate is matched against the
// first part of the given predicate.
// For example:
//   - matchPredicate("foo/0")("foo/0") -> true
//   - matchPredicate("foo/0")("foo/1") -> false
//   - matchPredicate("foo/0")("foo") -> true
//   - matchPredicate("foo/0")("bar") -> false
func PredicateEq(predicate string) func(b string) bool {
	return func(other string) bool {
		if strings.Contains(other, "/") {
			return predicate == other
		}
		return strings.Split(predicate, "/")[0] == other
	}
}
