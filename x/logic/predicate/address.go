package predicate

import (
	"context"
	"fmt"

	bech322 "github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
)

// Bech32Address is a predicate that convert a bech32 encoded string into base64 bytes and give the address prefix, or
// convert a prefix (HRP) and base64 encoded bytes to bech32 encoded string. The signature is as follows:
//
//	bech32_address(-Address, +Bech32)
//	bech32_address(+Address, -Bech32)
//	bech32_address(+Address, +Bech32)
//
// where:
//   - Address is a pair of, HRP (Human-Readable Part) containing the address prefix and the list of integers
//     between 0 and 255 (byte) that represent the base64 encoded bech32 address string. Represented like this : -(HRP, Address)
//   - Bech32 is an Atom or string representing the bech32 encoded string address
//
// # Example:
//
// 	- Convert the given bech32 address into base64 encoded byte by unify the prefix of given address (Hrp) and
// 	the base64 encoded value (Address).
//
//	bech32_address(-(Hrp, Address), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').
//
//  - Convert the given pair of HRP and base64 encoded address byte by unify the Bech32 string encoded value.
//
//	bech32_address(-('okp4', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).
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
