package prolog

import "github.com/ichiban/prolog/engine"

// Tuple is a predicate which unifies the given term with a tuple of the given arity.
func Tuple(args ...engine.Term) engine.Term {
	return engine.Atom("").Apply(args...)
}
