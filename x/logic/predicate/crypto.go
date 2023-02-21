package predicate

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
	"github.com/tendermint/tendermint/crypto"
)

func SHAHash(vm *engine.VM, data, hash engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var result []byte
		switch d := env.Resolve(data).(type) {
		case engine.Atom:
			result = crypto.Sha256([]byte(d.String()))
			return engine.Unify(vm, hash, BytesToList(result), cont, env)
		default:
			return engine.Error(fmt.Errorf("sha_hash/2: invalid data type: %T, should be Atom", d))
		}
	})
}

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
				return engine.Error(fmt.Errorf("hex_bytes/2: failed decode hexadecimal %w", err))
			}
		default:
			return engine.Error(fmt.Errorf("hex_bytes/2: invalid hex type: %T, should be Atom or Variable", h))
		}

		switch b := env.Resolve(bts).(type) {
		case engine.Variable:
			if result == nil {
				return engine.Error(fmt.Errorf("hex_bytes/2: nil hexadecimal conversion in input"))
			}
			return engine.Unify(vm, bts, BytesToList(result), cont, env)
		case engine.Compound:
			if b.Arity() != 2 || b.Functor().String() != "." {
				return engine.Error(fmt.Errorf("hex_bytes/2: bytes should be a List, give %T", b))
			}
			iter := engine.ListIterator{List: b, Env: env}

			src, err := ListToBytes(iter, env)
			if err != nil {
				return engine.Error(fmt.Errorf("hex_bytes/2: failed convert list into bytes: %w", err))
			}
			dst := hex.EncodeToString(src)
			return engine.Unify(vm, hexa, util.StringToTerm(dst), cont, env)
		default:
			return engine.Error(fmt.Errorf("hex_bytes/2: invalid hex type: %T, should be Variable or List", b))
		}
	})
}
