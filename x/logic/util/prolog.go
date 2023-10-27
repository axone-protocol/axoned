package util

import (
	"fmt"
	"strings"

	"github.com/ichiban/prolog/engine"
)

var (
	// AtomDot is the term used to represent the dot in a list.
	AtomDot = engine.NewAtom(".")

	// AtomEmpty is the term used to represent empty.
	AtomEmpty = engine.NewAtom("")
)

// StringToTerm converts a string to a term.
func StringToTerm(s string) engine.Term {
	return engine.NewAtom(s)
}

// ResolveToAtom resolves a term and attempts to convert it into an engine.Atom if possible.
// If conversion fails, the function returns the empty atom and the error.
func ResolveToAtom(env *engine.Env, t engine.Term) (engine.Atom, error) {
	switch t := env.Resolve(t).(type) {
	case engine.Atom:
		return t, nil
	default:
		return AtomEmpty,
			fmt.Errorf("invalid term '%s' - expected engine.Atom but got %T", t, t)
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

// IsList returns true if the given compound is a list.
func IsList(compound engine.Compound) bool {
	return compound.Functor() == AtomDot && compound.Arity() == 2
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

	compound, ok := resolvedTerm.(engine.Compound)
	if ok && IsList(compound) {
		iter := engine.ListIterator{List: compound, Env: env}

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
func GetOptionWithDefault(name engine.Atom, options engine.Term, defaultValue engine.Term, env *engine.Env) (engine.Term, error) {
	if term, err := GetOption(name, options, env); err != nil {
		return nil, err
	} else if term != nil {
		return term, nil
	}
	return defaultValue, nil
}
