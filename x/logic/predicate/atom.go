package predicate

import (
	"github.com/ichiban/prolog/engine"
)

// AtomPair are terms with principal functor (-)/2.
// For example, the term -(A, B) denotes the pair of elements A and B.
var AtomPair = engine.NewAtom("-")
