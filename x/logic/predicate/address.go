package predicate

import (
	"github.com/axone-protocol/prolog/v2/engine"

	bech322 "github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/axone-protocol/axoned/v10/x/logic/prolog"
)

// Bech32Address is a predicate that convert a [bech32] encoded string into [base64] bytes and give the address prefix,
// or convert a prefix (HRP) and [base64] encoded bytes to [bech32] encoded string.
//
// # Signature
//
//	bech32_address(-Address, +Bech32) is det
//	bech32_address(+Address, -Bech32) is det
//
// where:
//   - Address is a pair of the HRP (Human-Readable Part) which holds the address prefix and a list of numbers
//     ranging from 0 to 255 that represent the base64 encoded bech32 address string.
//   - Bech32 is an Atom or string representing the bech32 encoded string address
//
// [bech32]: https://docs.cosmos.network/main/build/spec/addresses/bech32#hrp-table
// [base64]: https://fr.wikipedia.org/wiki/Base64
func Bech32Address(_ *engine.VM, address, bech32 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	forwardConverter := func(value []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		hrpTerm, dataTerm, err := prolog.AssertPair(value[0], env)
		if err != nil {
			return nil, err
		}
		data, err := prolog.ByteListTermToBytes(dataTerm, env)
		if err != nil {
			return nil, err
		}
		hrp, err := prolog.AssertAtom(hrpTerm, env)
		if err != nil {
			return nil, err
		}

		b, err := bech322.ConvertAndEncode(hrp.String(), data)
		if err != nil {
			return nil, prolog.WithError(engine.DomainError(prolog.ValidEncoding("bech32"), value[0], env), err, env)
		}

		return []engine.Term{engine.NewAtom(b)}, nil
	}
	backwardConverter := func(value []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		b, err := prolog.AssertAtom(value[0], env)
		if err != nil {
			return nil, err
		}
		h, a, err := bech322.DecodeAndConvert(b.String())
		if err != nil {
			return nil, prolog.WithError(engine.DomainError(prolog.ValidEncoding("bech32"), value[0], env), err, env)
		}
		var r engine.Term = engine.NewAtom(h)
		pair := prolog.AtomPair.Apply(r, prolog.BytesToByteListTerm(a))
		return []engine.Term{pair}, nil
	}
	return prolog.UnifyFunctionalPredicate(
		[]engine.Term{address}, []engine.Term{bech32}, prolog.AtomEmpty, forwardConverter, backwardConverter, cont, env)
}
