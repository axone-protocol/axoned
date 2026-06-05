package prolog

import (
	"unicode/utf8"

	"github.com/axone-protocol/prolog/v3/engine"
)

// IsList returns true if the given term is a list.
func IsList(term engine.Term, env *engine.Env) bool {
	_, err := AssertList(term, env)
	return err == nil
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

// AssertKeyValue resolves a term as a key-value and returns its components, the key as an atom,
// and the value as a term.
// If conversion fails, the function returns nil and the error.
func AssertKeyValue(term engine.Term, env *engine.Env) (engine.Atom, engine.Term, error) {
	k, v, err := assertTuple2WithFunctor(term, AtomKeyValue, AtomTypeKeyValue, env)
	if err != nil {
		return AtomEmpty, nil, err
	}

	key, err := AssertAtom(k, env)
	if err != nil {
		return AtomEmpty, nil, err
	}

	return key, v, err
}

// assertTuple2WithFunctor resolves a term as a tuple and returns the tuple components based on the given functor.
// If conversion fails, the function returns nil and an error.
func assertTuple2WithFunctor(
	term engine.Term, functor engine.Atom, functorType engine.Atom, env *engine.Env,
) (engine.Term, engine.Term, error) {
	term, err := AssertIsGround(term, env)
	if err != nil {
		return nil, nil, err
	}
	if compound, ok := term.(engine.Compound); ok && compound.Functor() == functor && compound.Arity() == 2 {
		return compound.Arg(0), compound.Arg(1), nil
	}

	return nil, nil, engine.TypeError(functorType, term, env)
}
