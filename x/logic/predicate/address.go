package predicate

import (
	"context"
	"fmt"

	bech322 "github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
)

func Bech32Address(vm *engine.VM, address, bech32 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		switch b := env.Resolve(bech32).(type) {
		case engine.Variable:
		case engine.Atom:
			h, a, err := bech322.DecodeAndConvert(b.String())
			if err != nil {
				return engine.Error(fmt.Errorf("bech32_address/2: failed convert bech32 encoded string to base64: %w", err))
			}
			pair := AtomPair.Apply(util.StringToTerm(h), BytesToList(a))
			return engine.Unify(vm, address, pair, cont, env)
		default:
			return engine.Error(fmt.Errorf("bech32_address/2: invalid Bech32 type: %T, should be Atom or Variable", b))
		}

		switch addressPair := env.Resolve(address).(type) {
		case engine.Compound:
			bech32Decoded, err := AddressPairToBech32(addressPair, env)
			if err != nil {
				return engine.Error(fmt.Errorf("bech32_address/2: %w", err))
			}
			return engine.Unify(vm, bech32, util.StringToTerm(bech32Decoded), cont, env)
		default:
			return engine.Error(fmt.Errorf("bech32_address/2: you should give at least on instantiated value (Address or Bech32)"))
		}
	})
}

func AddressPairToBech32(addressPair engine.Compound, env *engine.Env) (string, error) {
	if addressPair.Functor() != AtomPair || addressPair.Arity() != 2 {
		return "", fmt.Errorf("address should be a Pair '-(Hrp, Address)'")
	}

	switch a := env.Resolve(addressPair.Arg(1)).(type) {
	case engine.Compound:
		if a.Arity() != 2 || a.Functor().String() != "." {
			return "", fmt.Errorf("address should be a List of bytes, give %s/%d", a.Functor().String(), a.Arity())
		}

		iter := engine.ListIterator{List: a, Env: env}
		data, err := ListToBytes(iter, env)
		if err != nil {
			return "", fmt.Errorf("failed convert term to bytes list: %w", err)
		}
		hrp, ok := env.Resolve(addressPair.Arg(0)).(engine.Atom)
		if !ok {
			return "", fmt.Errorf("HRP should be instantiated when trying convert bytes to bech32")
		}
		b, err := bech322.ConvertAndEncode(hrp.String(), data)
		if err != nil {
			return "", fmt.Errorf("failed convert base64 encoded address to bech32 string encoded: %w", err)
		}

		return b, nil
	default:
		return "", fmt.Errorf("address should be a Pair with a List of bytes in arity 2, give %T", addressPair.Arg(1))
	}
}
