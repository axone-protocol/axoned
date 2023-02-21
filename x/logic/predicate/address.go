package predicate

import (
	"context"
	"fmt"

	bech322 "github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
)

func Bech32Address(vm *engine.VM, hrp, address, bech32 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		switch b := env.Resolve(bech32).(type) {
		case engine.Variable:
			return nil
		case engine.Atom:
			h, a, err := bech322.DecodeAndConvert(b.String())
			if err != nil {
				return engine.Error(fmt.Errorf("bech32_address/3: failed convert bech32 encoded string to base64: %w", err))
			}
			return engine.Unify(vm, Tuple(hrp, address), Tuple(util.StringToTerm(h), BytesToList(a)), cont, env)
		default:
			return engine.Error(fmt.Errorf("bech32_address/3: invalid data type: %T, should be Atom or Variable", b))
		}
	})
}
