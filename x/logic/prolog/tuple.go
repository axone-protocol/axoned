package prolog

import "github.com/axone-protocol/prolog/v2/engine"

// Tuple is a predicate which unifies the given term with a tuple of the given arity.
func Tuple(args ...engine.Term) engine.Term {
	return engine.Atom("").Apply(args...)
}
