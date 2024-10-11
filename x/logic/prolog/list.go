package prolog

import (
	"github.com/axone-protocol/prolog/engine"
)

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

// ForEach iterates over the elements of the given list and calls the given function for each element.
func ForEach(list engine.Term, env *engine.Env, f func(v engine.Term, hasNext bool) error) error {
	iter, err := ListIterator(list, env)
	if err != nil {
		return err
	}

	if !iter.Next() {
		return nil
	}

	for {
		elem := iter.Current()
		hasNext := iter.Next()

		if err := f(elem, hasNext); err != nil {
			return err
		}

		if !hasNext {
			break
		}
	}

	return nil
}
