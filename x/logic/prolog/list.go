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

// ListIterator returns a list iterator.
func ListIterator(list engine.Term, env *engine.Env) (engine.ListIterator, error) {
	if !IsList(list, env) {
		return engine.ListIterator{}, engine.TypeError(AtomTypeList, list, env)
	}
	return engine.ListIterator{List: list, Env: env}, nil
}
