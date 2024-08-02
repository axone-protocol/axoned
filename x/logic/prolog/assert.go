package prolog

import (
	"strings"
	"unicode/utf8"

	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"

	"github.com/axone-protocol/axoned/v9/x/logic/util"
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
func IsList(term engine.Term, env *engine.Env) bool {
	_, err := AssertList(term, env)
	return err == nil
}

// IsEmptyList returns true if the given term is an empty list.
func IsEmptyList(term engine.Term, env *engine.Env) bool {
	if v, ok := env.Resolve(term).(engine.Atom); ok {
		return v == AtomEmptyList
	}
	return false
}

// IsGround returns true if the given term holds no free variables.
func IsGround(term engine.Term, env *engine.Env) bool {
	_, err := AssertIsGround(term, env)
	return err == nil
}

func AreGround(terms []engine.Term, env *engine.Env) bool {
	return lo.EveryBy(terms, func(t engine.Term) bool {
		return IsGround(t, env)
	})
}

// AssertIsGround resolves a term and returns it if it is ground.
// If the term is not ground, the function returns nil and the instantiation error.
func AssertIsGround(term engine.Term, env *engine.Env) (engine.Term, error) {
	switch term := env.Resolve(term).(type) {
	case engine.Variable:
		return nil, engine.InstantiationError(env)
	case engine.Compound:
		args := make([]engine.Term, term.Arity())
		for i := 0; i < term.Arity(); i++ {
			arg, err := AssertIsGround(term.Arg(i), env)
			if err != nil {
				return nil, err
			}
			args[i] = arg
		}
		return term.Functor().Apply(args...), nil
	default:
		return term, nil
	}
}

// AssertAtom resolves a term and attempts to convert it into an engine.Atom if possible.
// If conversion fails, the function returns the empty atom and the error.
func AssertAtom(term engine.Term, env *engine.Env) (engine.Atom, error) {
	switch term := env.Resolve(term).(type) {
	case engine.Atom:
		return term, nil
	case engine.Variable:
		return AtomEmpty, engine.InstantiationError(env)
	default:
		return AtomEmpty, engine.TypeError(AtomTypeAtom, term, env)
	}
}

// AssertCharacterCode resolves a term and attempts to convert it into a rune if possible.
// If conversion fails, the function returns the zero value and the error.
func AssertCharacterCode(term engine.Term, env *engine.Env) (rune, error) {
	switch term := env.Resolve(term).(type) {
	case engine.Integer:
		if term >= 0 && term <= utf8.MaxRune {
			return rune(term), nil
		}
	case engine.Variable:
		return utf8.RuneError, engine.InstantiationError(env)
	}

	return utf8.RuneError, engine.TypeError(AtomTypeCharacterCode, term, env)
}

// AssertCharacter resolves a term and attempts to convert it into an engine.Atom if possible.
// If conversion fails, the function returns the empty atom and the error.
func AssertCharacter(term engine.Term, env *engine.Env) (rune, error) {
	switch term := env.Resolve(term).(type) {
	case engine.Atom:
		runes := []rune(term.String())
		if len(runes) == 1 {
			return runes[0], nil
		}
	case engine.Variable:
		return utf8.RuneError, engine.InstantiationError(env)
	}

	return utf8.RuneError, engine.TypeError(AtomTypeCharacter, term, env)
}

// AssertByte resolves a term and attempts to convert it into a byte if possible.
// If conversion fails, the function returns the zero value and the error.
func AssertByte(term engine.Term, env *engine.Env) (byte, error) {
	switch term := env.Resolve(term).(type) {
	case engine.Integer:
		if term >= 0 && term <= 255 {
			return byte(term), nil
		}
	case engine.Variable:
		return 0, engine.InstantiationError(env)
	}
	return 0, engine.TypeError(AtomTypeByte, term, env)
}

// AssertList resolves a term as a list and returns it as a engine.Compound.
// If conversion fails, the function returns nil and the error.
func AssertList(term engine.Term, env *engine.Env) (engine.Term, error) {
	switch term := env.Resolve(term).(type) {
	case engine.Compound:
		if term.Functor() == AtomDot && term.Arity() == 2 {
			return term, nil
		}
	case engine.Atom:
		if term == AtomEmptyList {
			return term, nil
		}
	}

	return nil, engine.TypeError(AtomTypeList, term, env)
}

// AssertPair resolves a term as a pair and returns the pair components.
// If conversion fails, the function returns nil and the error.
func AssertPair(term engine.Term, env *engine.Env) (engine.Term, engine.Term, error) {
	term, err := AssertIsGround(term, env)
	if err != nil {
		return nil, nil, err
	}
	if term, ok := term.(engine.Compound); ok && term.Functor() == AtomPair && term.Arity() == 2 {
		return term.Arg(0), term.Arg(1), nil
	}

	return nil, nil, engine.TypeError(AtomTypePair, term, env)
}

// AssertURIComponent resolves a term as a URI component and returns it as an URIComponent.
func AssertURIComponent(term engine.Term, env *engine.Env) (util.URIComponent, error) {
	switch v := env.Resolve(term); v {
	case AtomQueryValue:
		return util.QueryValueComponent, nil
	case AtomFragment:
		return util.FragmentComponent, nil
	case AtomPath:
		return util.PathComponent, nil
	case AtomSegment:
		return util.SegmentComponent, nil
	default:
		return 0, engine.TypeError(AtomTypeURIComponent, term, env)
	}
}
