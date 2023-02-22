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
		case engine.Atom:
			h, a, err := bech322.DecodeAndConvert(b.String())
			if err != nil {
				return engine.Error(fmt.Errorf("bech32_address/3: failed convert bech32 encoded string to base64: %w", err))
			}
			return engine.Unify(vm, Tuple(hrp, address), Tuple(util.StringToTerm(h), BytesToList(a)), cont, env)
		default:
			return engine.Error(fmt.Errorf("bech32_address/3: invalid data type: %T, should be Atom or Variable", b))
		}

		switch a := env.Resolve(address).(type) {
		case engine.Compound:
			if a.Arity() != 2 || a.Functor().String() != "." {
				return engine.Error(fmt.Errorf("bech32_address/3: Address should be a List of bytes, give %s/%d", a.Functor().String(), a.Arity()))
			}

			iter := engine.ListIterator{List: a, Env: env}
			data, err := ListToBytes(iter, env)
			if err != nil {
				return engine.Error(fmt.Errorf("bech32_address/3: failed convert term to bytes list: %w", err))
			}
			h, ok := env.Resolve(hrp).(engine.Atom)
			if !ok {
				return engine.Error(fmt.Errorf("bech32_address/3: Hrp should be instantiated in Address convertion context"))
			}
			b, err := bech322.ConvertAndEncode(h.String(), data)
			if err != nil {
				return engine.Error(fmt.Errorf("bech32_address/3: failed convert base64 encoded address to bech32 string encoded: %w", err))
			}

			return engine.Unify(vm, bech32, util.StringToTerm(b), cont, env)
		default:
			return engine.Error(fmt.Errorf("bech32_address/3: Address should be a List of bytes when bech32 string encoded value is given"))
		}
	})
}
