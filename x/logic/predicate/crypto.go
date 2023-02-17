package predicate

import (
	"context"
	"fmt"

	"github.com/ichiban/prolog/engine"
	"github.com/tendermint/tendermint/crypto"
)

func CryptoHash(vm *engine.VM, data, hash engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		switch h := env.Resolve(hash).(type) {
		case engine.Variable:
			switch d := env.Resolve(data).(type) {
			case engine.Atom:
				result := crypto.Sha256([]byte(d.String()))
				return engine.Unify(vm, hash, BytesToList(result), cont, env)
			default:
				return engine.Error(fmt.Errorf("crypto_hash/2: invalid data type: %T, should be Atom", d))
			}
		default:
			return engine.Error(fmt.Errorf("crypto_hash/2: invalid hash type: %T, should be Variable", h))
		}
	})
}
