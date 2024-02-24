package predicate

import (
	"github.com/ichiban/prolog/engine"

	bech322 "github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/okp4/okp4d/x/logic/prolog"
)

// Bech32Address is a predicate that convert a [bech32] encoded string into [base64] bytes and give the address prefix,
// or convert a prefix (HRP) and [base64] encoded bytes to [bech32] encoded string.
//
// The signature is as follows:
//
//	bech32_address(-Address, +Bech32)
//	bech32_address(+Address, -Bech32)
//	bech32_address(+Address, +Bech32)
//
// where:
//   - Address is a pair of the HRP (Human-Readable Part) which holds the address prefix and a list of numbers
//     ranging from 0 to 255 that represent the base64 encoded bech32 address string.
//   - Bech32 is an Atom or string representing the bech32 encoded string address
//
// # Examples:
//
//	# Convert the given bech32 address into base64 encoded byte by unify the prefix of given address (Hrp) and the
//	base64 encoded value (Address).
//	- bech32_address(-(Hrp, Address), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').
//
//	# Convert the given pair of HRP and base64 encoded address byte by unify the Bech32 string encoded value.
//	- bech32_address(-('okp4', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).
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
