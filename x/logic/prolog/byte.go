package prolog

import (
	"github.com/ichiban/prolog/engine"
)

// ByteListTermToBytes try to convert a given list of bytes into native golang []byte.
func ByteListTermToBytes(term engine.Term, env *engine.Env) ([]byte, error) {
	iter, err := ListIterator(term, env)
	if err != nil {
		return nil, err
	}
	var bs []byte

	for iter.Next() {
		b, err := AssertByte(env, iter.Current())
		if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}

// BytesToByteListTerm try to convert a given golang []byte into a list of bytes.
func BytesToByteListTerm(in []byte) engine.Term {
	terms := make([]engine.Term, 0, len(in))
	for _, b := range in {
		terms = append(terms, engine.Integer(b))
	}
	return engine.List(terms...)
}
