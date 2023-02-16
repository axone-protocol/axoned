package predicate

import (
	"context"
	"encoding/hex"
	"fmt"

	"github.com/ichiban/prolog/engine"
	"github.com/tendermint/tendermint/crypto"
)

func CryptoHash(vm *engine.VM, data, hash engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		switch d := env.Resolve(data).(type) {
		case engine.Atom:
			result := crypto.Sha256([]byte(d.String()))
			return engine.Unify(vm, hash, engine.NewAtom(hex.EncodeToString(result)), cont, env)
		default:
			return engine.Error(fmt.Errorf("crypto_hash/2: cannot unify %s from %s", data, hash))
		}
	})
}
