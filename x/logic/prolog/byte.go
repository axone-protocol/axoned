package prolog

import (
	"github.com/axone-protocol/prolog/v3/engine"
)

// ByteListTermToBytes try to convert a given list of bytes into native golang []byte.
func ByteListTermToBytes(term engine.Term, env *engine.Env) ([]byte, error) {
	iter, err := ListIterator(term, env)
	if err != nil {
		return nil, err
	}
	var bs []byte

	for iter.Next() {
		b, err := AssertByte(iter.Current(), env)
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

// BytesToAtom converts a given golang []byte into an Atom.
func BytesToAtom(in []byte) engine.Atom {
	return engine.NewAtom(string(in))
}
