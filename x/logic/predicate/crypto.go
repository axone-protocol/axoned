package predicate

import (
	"context"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/ichiban/prolog/engine"

	cometcrypto "github.com/cometbft/cometbft/crypto"

	"github.com/okp4/okp4d/x/logic/util"
)

// SHAHash is a predicate that computes the Hash of the given Data.
//
// The signature is as follows:
//
//	sha_hash(+Data, -Hash) is det
//	sha_hash(+Data, +Hash) is det
//
// Where:
//   - Data represents the data to be hashed with the SHA-256 algorithm.
//   - Hash is the variable that will contain Hashed value of Data.
//
// Note: Due to the principles of the hash algorithm (pre-image resistance), this predicate can only compute the hash
// value from input data, and cannot compute the original input data from the hash value.
//
// Examples:
//
//	# Compute the hash of the given data and unify it with the given Hash.
//	- sha_hash("Hello OKP4", Hash).
func SHAHash(vm *engine.VM, data, hash engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var result []byte
		switch d := env.Resolve(data).(type) {
		case engine.Atom:
			result = cometcrypto.Sha256([]byte(d.String()))
			return engine.Unify(vm, hash, BytesToList(result), cont, env)
		default:
			return engine.Error(fmt.Errorf("sha_hash/2: invalid data type: %T, should be Atom", d))
		}
	})
}

// HexBytes is a predicate that unifies hexadecimal encoded bytes to a list of bytes.
//
// The signature is as follows:
//
//	hex_bytes(?Hex, ?Bytes) is det
//
// Where:
//   - Hex is an Atom, string or list of characters in hexadecimal encoding.
//   - Bytes is the list of numbers between 0 and 255 that represent the sequence of bytes.
//
// Examples:
//
//	# Convert hexadecimal atom to list of bytes.
//	- hex_bytes('2c26b46b68ffc68ff99b453c1d3041341342d706483bfa0f98a5e886266e7ae', Bytes).
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

// EdDSAVerify return `true` if the Signature can be verified as the EdDSA signature for Data, using the given PubKey
// as bytes.
//
// eddsa_verify(+PubKey, +Data, +Signature, +Options) is semidet
//
// Where:
// - PubKey is a list of bytes representing the public key.
// - Data is the hash of the signed message could be an Atom or List of bytes.
// - Signature is the signature of the Data, as list of bytes.
// - Options allow to give option to the predicates, available options are:
//   - encoding(+Format): Encoding to use for the given Data. Possible values are:
//     -- `hex` (default): hexadecimal encoding represented as an atom.
//     -- `octet`: plain bytes encoding represented as a list of integers between 0 and 255.
//   - type(+Alg): Algorithm to use in the EdDSA family. Supported algorithms are:
//     -- `ed25519` (default): the EdDSA signature scheme using SHA-512 (SHA-2) and Curve25519.
//
// Examples:
//
// # Verify the signature of given hexadecimal data.
// - eddsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', [23, 56, ...], [encoding(hex), type(ed25519)]).
//
// # Verify the signature of given binary data.
// - eddsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(ed25519)]).
func EdDSAVerify(_ *engine.VM, key, data, sig, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return xVerify("eddsa_verify/4", key, data, sig, options, util.Ed25519, []util.Alg{util.Ed25519}, cont, env)
}

// ECDSAVerify return `true` if the Signature can be verified as the ECDSA signature for Data, using the given PubKey
// as bytes.
//
// ecdsa_verify(+PubKey, +Data, +Signature, +Options) is semidet
//
// Where:
// - PubKey is a list of bytes representing the public key.
// - Data is the hash of the signed message could be an Atom or List of bytes.
// - Signature is the signature of the Data, as list of bytes.
// - Options allow to give option to the predicates, available options are:
//   - encoding(+Format): Encoding to use for the given Data. Possible values are:
//     -- `hex` (default): hexadecimal encoding represented as an atom.
//     -- `octet`: plain bytes encoding represented as a list of integers between 0 and 255.
//   - type(+Alg): Algorithm to use in the EdDSA family. Supported algorithms are:
//     -- `secp256r1` (default):
//     -- `secp256k1`:
//
// Examples:
//
// # Verify the signature of given hexadecimal data as ECDSA secp256r1 algorithm.
// - ecdsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', [23, 56, ...], encoding(hex)).
//
// # Verify the signature of given binary data as ECDSA secp256k1 algorithm.
// - ecdsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(secp256k1)]).
func ECDSAVerify(_ *engine.VM, key, data, sig, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return xVerify("ecdsa_verify/4", key, data, sig, options, util.Secp256r1, []util.Alg{util.Secp256r1, util.Secp256k1}, cont, env)
}

// xVerify return `true` if the Signature can be verified as the signature for Data, using the given PubKey for a
// considered algorithm.
// This is a generic predicate implementation that can be used to verify any signature.
func xVerify(functor string, key, data, sig, options engine.Term, defaultAlgo util.Alg, algos []util.Alg, cont engine.Cont, env *engine.Env) *engine.Promise {
	typeOpt := engine.NewAtom("type")
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		typeTerm, err := util.GetOptionWithDefault(typeOpt, options, engine.NewAtom(defaultAlgo.String()), env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: %w", functor, err))
		}
		typeAtom, err := util.ResolveToAtom(env, typeTerm)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: %w", functor, err))
		}

		if idx := slices.IndexFunc(algos, func(a util.Alg) bool { return a.String() == typeAtom.String() }); idx == -1 {
			return engine.Error(fmt.Errorf("%s: invalid type: %s. Possible values: %s",
				functor,
				typeAtom.String(),
				strings.Join(util.Map(algos, func(a util.Alg) string { return a.String() }), ", ")))
		}

		decodedKey, err := TermToBytes(key, AtomEncoding.Apply(AtomOctet), env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode public key: %w", functor, err))
		}

		decodedData, err := TermToBytes(data, options, env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode data: %w", functor, err))
		}

		decodedSignature, err := TermToBytes(sig, AtomEncoding.Apply(AtomOctet), env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode signature: %w", functor, err))
		}

		r, err := util.VerifySignature(util.Alg(typeAtom.String()), decodedKey, decodedData, decodedSignature)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to verify signature: %w", functor, err))
		}

		if !r {
			return engine.Bool(false)
		}

		return cont(env)
	})
}
