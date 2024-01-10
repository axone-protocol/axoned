package prolog

import "github.com/ichiban/prolog/engine"

// ListHead returns the first element of the given list.
func ListHead(list engine.Term, env *engine.Env) engine.Term {
	iter := engine.ListIterator{List: list, Env: env}
	if !iter.Next() {
		return nil
	}
	return iter.Current()
}
