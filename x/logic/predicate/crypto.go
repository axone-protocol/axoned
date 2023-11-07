package predicate

import (
	"context"
	"encoding/hex"
	"fmt"
	"slices"
	"strings"

	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/x/logic/util"
)

// SHAHash is a predicate that computes the Hash of the given Data.
//
// Deprecated: sha_hash/2 should not be used anymore as it will be removed in a future release.
// Use the new crypto_data_hash/3 predicate instead.
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
//	- sha_hash('Hello OKP4', Hash).
func SHAHash(vm *engine.VM, data, hash engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		switch d := env.Resolve(data).(type) {
		case engine.Atom:
			result, err := util.Hash(util.HashAlgSha256, []byte(d.String()))
			if err != nil {
				engine.Error(fmt.Errorf("sha_hash/2: failed to hash data: %w", err))
			}

			return engine.Unify(vm, hash, BytesToList(result), cont, env)
		default:
			return engine.Error(fmt.Errorf("sha_hash/2: invalid data type: %T, should be Atom", d))
		}
	})
}

// CryptoDataHash is a predicate that computes the Hash of the given Data using different algorithms.
//
// The signature is as follows:
//
//	crypto_data_hash(+Data, -Hash, +Options) is det
//	crypto_data_hash(+Data, +Hash, +Options) is det
//
// Where:
//   - Data represents the data to be hashed with the SHA-256 algorithm, given as an atom, or code-list.
//   - Hash represents the Hashed value of Data, which can be given as an atom or a variable..
//   - Options are additional configurations for the hashing process. Supported options include:
//     encoding(+Format) which specifies the encoding used for the Data, and algorithm(+Alg) which chooses the hashing
//     algorithm among the supported ones (see below for details).
//
// For Format, the supported encodings are:
//
//   - utf8 (default), the UTF-8 encoding represented as an atom.
//   - hex, the hexadecimal encoding represented as an atom.
//   - octet, the raw byte encoding depicted as a list of integers ranging from 0 to 255.
//
// For Alg, the supported algorithms are:
//
//   - sha256 (default): The SHA-256 algorithm.
//
// Note: Due to the principles of the hash algorithm (pre-image resistance), this predicate can only compute the hash
// value from input data, and cannot compute the original input data from the hash value.
//
// Examples:
//
//		# Compute the SHA-256 hash of the given data and unify it with the given Hash.
//		- crypto_data_hash('Hello OKP4', Hash).
//
//		# Compute the SHA-256 hash of the given hexadecimal data and unify it with the given Hash.
//		- crypto_data_hash('9b038f8ef6918cbb56040dfda401b56b...', Hash, encoding(hex)).
//
//	 # Compute the SHA-256 hash of the given hexadecimal data and unify it with the given Hash.
//	 - crypto_data_hash([127, ...], Hash, encoding(octet)).
func CryptoDataHash(
	vm *engine.VM, data, hash, options engine.Term, cont engine.Cont, env *engine.Env,
) *engine.Promise {
	functor := "crypto_data_hash/3"
	algorithmOpt := engine.NewAtom("algorithm")

	return engine.Delay(func(ctx context.Context) *engine.Promise {
		algorithm, err := getOptionAsAtomWithDefault(algorithmOpt, options, engine.NewAtom("sha256"), env, functor)
		if err != nil {
			return engine.Error(err)
		}
		decodedData, err := TermToBytes(data, options, AtomUtf8, env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode data: %w", functor, err))
		}

		switch algorithm.String() {
		case util.HashAlgSha256.String():
			result, err := util.Hash(util.HashAlgSha256, decodedData)
			if err != nil {
				engine.Error(fmt.Errorf("sha_hash/2: failed to hash data: %w", err))
			}

			return engine.Unify(vm, hash, BytesToList(result), cont, env)
		default:
			return engine.Error(fmt.Errorf("%s: invalid algorithm: %s. Possible values: %s",
				functor,
				algorithm.String(),
				util.HashAlgNames()))
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

// EDDSAVerify determines if a given signature is valid as per the EdDSA algorithm for the provided data, using the
// specified public key.
//
// The signature is as follows:
//
//	eddsa_verify(+PubKey, +Data, +Signature, +Options) is semi-det
//
// Where:
//   - PubKey is the encoded public key as a list of bytes.
//   - Data is the message to verify, represented as either a hexadecimal atom or a list of bytes.
//     It's important that the message isn't pre-hashed since the Ed25519 algorithm processes
//     messages in two passes when signing.
//   - Signature represents the signature corresponding to the data, provided as a list of bytes.
//   - Options are additional configurations for the verification process. Supported options include:
//     encoding(+Format) which specifies the encoding used for the Data, and type(+Alg) which chooses the algorithm
//     within the EdDSA family (see below for details).
//
// For Format, the supported encodings are:
//
//   - hex (default), the hexadecimal encoding represented as an atom.
//   - octet, the plain byte encoding depicted as a list of integers ranging from 0 to 255.
//
// For Alg, the supported algorithms are:
//
//   - ed25519 (default): The EdDSA signature scheme using SHA-512 (SHA-2) and Curve25519.
//
// Examples:
//
//	# Verify a signature for a given hexadecimal data.
//	- eddsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56b...', [23, 56, ...], [encoding(hex), type(ed25519)])
//
//	# Verify a signature for binary data.
//	- eddsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(ed25519)])
func EDDSAVerify(_ *engine.VM, key, data, sig, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return xVerify("eddsa_verify/4", key, data, sig, options, util.Ed25519, []util.KeyAlg{util.Ed25519}, cont, env)
}

// ECDSAVerify determines if a given signature is valid as per the ECDSA algorithm for the provided data, using the
// specified public key.
//
// The signature is as follows:
//
//	ecdsa_verify(+PubKey, +Data, +Signature, +Options), which is semi-deterministic.
//
// Where:
//
//   - PubKey is the 33-byte compressed public key, as specified in section 4.3.6 of ANSI X9.62.
//
//   - Data is the hash of the signed message, which can be either an atom or a list of bytes.
//
//   - Signature represents the ASN.1 encoded signature corresponding to the Data.
//
//   - Options are additional configurations for the verification process. Supported options include:
//     encoding(+Format) which specifies the encoding used for the data, and type(+Alg) which chooses the algorithm
//     within the ECDSA family (see below for details).
//
// For Format, the supported encodings are:
//
//   - hex (default), the hexadecimal encoding represented as an atom.
//   - octet, the plain byte encoding depicted as a list of integers ranging from 0 to 255.
//
// For Alg, the supported algorithms are:
//
//   - secp256r1 (default): Also known as P-256 and prime256v1.
//   - secp256k1: The Koblitz elliptic curve used in Bitcoin's public-key cryptography.
//
// Examples:
//
//	# Verify a signature for hexadecimal data using the ECDSA secp256r1 algorithm.
//	- ecdsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56b...', [23, 56, ...], encoding(hex))
//
//	# Verify a signature for binary data using the ECDSA secp256k1 algorithm.
//	- ecdsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(secp256k1)])
func ECDSAVerify(_ *engine.VM, key, data, sig, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return xVerify("ecdsa_verify/4", key, data, sig, options, util.Secp256r1, []util.KeyAlg{util.Secp256r1, util.Secp256k1}, cont, env)
}

// xVerify return `true` if the Signature can be verified as the signature for Data, using the given PubKey for a
// considered algorithm.
// This is a generic predicate implementation that can be used to verify any signature.
func xVerify(functor string, key, data, sig, options engine.Term, defaultAlgo util.KeyAlg,
	algos []util.KeyAlg, cont engine.Cont, env *engine.Env,
) *engine.Promise {
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

		if idx := slices.IndexFunc(algos, func(a util.KeyAlg) bool { return a.String() == typeAtom.String() }); idx == -1 {
			return engine.Error(fmt.Errorf("%s: invalid type: %s. Possible values: %s",
				functor,
				typeAtom.String(),
				strings.Join(util.Map(algos, func(a util.KeyAlg) string { return a.String() }), ", ")))
		}

		decodedKey, err := TermToBytes(key, AtomEncoding.Apply(AtomOctet), AtomHex, env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode public key: %w", functor, err))
		}

		decodedData, err := TermToBytes(data, options, AtomHex, env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode data: %w", functor, err))
		}

		decodedSignature, err := TermToBytes(sig, AtomEncoding.Apply(AtomOctet), AtomHex, env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode signature: %w", functor, err))
		}

		r, err := util.VerifySignature(util.KeyAlg(typeAtom.String()), decodedKey, decodedData, decodedSignature)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to verify signature: %w", functor, err))
		}

		if !r {
			return engine.Bool(false)
		}

		return cont(env)
	})
}

// getOptionAsAtomWithDefault is a helper function that returns the value of the first option with the given name in the
// given options.
func getOptionAsAtomWithDefault(algorithmOpt engine.Atom, options engine.Term, defaultValue engine.Term, env *engine.Env,
	functor string,
) (engine.Atom, error) {
	term, err := util.GetOptionWithDefault(algorithmOpt, options, defaultValue, env)
	if err != nil {
		return util.AtomEmpty, fmt.Errorf("%s: %w", functor, err)
	}
	atom, err := util.ResolveToAtom(env, term)
	if err != nil {
		return util.AtomEmpty, fmt.Errorf("%s: %w", functor, err)
	}

	return atom, nil
}
