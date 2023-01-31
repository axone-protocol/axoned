package util

import "github.com/ichiban/prolog/engine"

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
