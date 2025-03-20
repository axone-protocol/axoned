package predicate

import (
	"github.com/axone-protocol/prolog/v2/engine"

	bech322 "github.com/cosmos/cosmos-sdk/types/bech32"

	"github.com/axone-protocol/axoned/v11/x/logic/prolog"
)

// Bech32Address is a predicate that converts a Bech32-encoded string into a prefix (HRP) and Base64-encoded bytes,
// or constructs a Bech32-encoded string from a prefix and Base64 bytes.
//
// This predicate handles Bech32 address encoding and decoding as per the Cosmos specification. In the Cosmos ecosystem,
// most chains (e.g., Cosmos Hub, Akash) share the BIP-44 coin type 118', allowing HRP conversion (e.g., 'cosmos' to 'akash')
// to produce valid addresses from the same underlying key.
//
// # Signature
//
//	bech32_address(-Address, +Bech32) is det
//	bech32_address(+Address, -Bech32) is det
//
// where:
//
//   - Address: A pair `HRP-Base64Bytes`, where: HRP is an atom representing the Human-Readable Part (e.g. 'cosmos', 'akash',
//     'axone'), and Base64Bytes is a list of integers (0-255) representing the Base64-encoded bytes git statof the address.
//   - Bech32: An atom or string representing the Bech32-encoded address (e.g., 'cosmos17sc02mcgjzdv5l4jwnzffxw7g60y5ta4pggcp4').
//
// # Limitations
//
// Conversion between HRPs is only valid for chains sharing the same BIP-44 coin type (e.g., 118'). For chains with
// distinct coin types (e.g., Secret: 529', Bitsong: 639'), this predicate cannot derive the correct address from another
// chain’s Bech32 string.
//
// # References
//
//   - [Bech32 on Cosmos]
//
//   - [Base64 Encoding]
//
//   - [Cosmos Chain Registry]
//
//   - [BIP 44]
//
// [Bech32 on Cosmos]: https://docs.cosmos.network/main/build/spec/addresses/bech32
// [Base64 Encoding]: https://fr.wikipedia.org/wiki/Base64
// [Cosmos Chain Registry]: https://github.com/cosmos/chain-registry
// [BIP 44]: https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
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
