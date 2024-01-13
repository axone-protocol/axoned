package predicate

import (
	"context"
	"encoding/hex"

	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/x/logic/prolog"
)

// HexBytes is a predicate that unifies hexadecimal encoded bytes to a list of bytes.
//
// The signature is as follows:
//
//	hex_bytes(?Hex, ?Bytes) is det
//
// Where:
//   - Hex is an Atom, string or list of characters in hexadecimal encoding.
//   - Bytes is the list of numbers between 0 and 255 that represent the sequence of bytes.
//
// Examples:
//
//	# Convert hexadecimal atom to list of bytes.
//	- hex_bytes('2c26b46b68ffc68ff99b453c1d3041341342d706483bfa0f98a5e886266e7ae', Bytes).
func HexBytes(vm *engine.VM, hexa, bts engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var result []byte

		switch h := env.Resolve(hexa).(type) {
		case engine.Variable:
		case engine.Atom:
			src := []byte(h.String())
			result = make([]byte, hex.DecodedLen(len(src)))
			_, err := hex.Decode(result, src)
			if err != nil {
				return engine.Error(
					prolog.WithError(
						engine.DomainError(prolog.ValidEncoding("hex"), hexa, env), err, env))
			}
		default:
			return engine.Error(engine.TypeError(prolog.AtomTypeAtom, hexa, env))
		}

		switch b := env.Resolve(bts).(type) {
		case engine.Variable:
			if result == nil {
				return engine.Error(engine.InstantiationError(env))
			}
			return engine.Unify(vm, bts, prolog.BytesToByteListTerm(result), cont, env)
		case engine.Compound:
			src, err := prolog.ByteListTermToBytes(b, env)
			if err != nil {
				return engine.Error(err)
			}
			dst := hex.EncodeToString(src)
			var r engine.Term = engine.NewAtom(dst)
			return engine.Unify(vm, hexa, r, cont, env)
		default:
			return engine.Error(engine.TypeError(prolog.AtomTypeText, bts, env))
		}
	})
}
