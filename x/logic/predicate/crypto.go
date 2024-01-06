package predicate

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/x/logic/util"
)

// CryptoDataHash is a predicate that computes the Hash of the given Data using different algorithms.
//
// The signature is as follows:
//
//	crypto_data_hash(+Data, -Hash, +Options) is det
//	crypto_data_hash(+Data, +Hash, +Options) is det
//
// Where:
//   - Data represents the data to be hashed, given as an atom, or code-list.
//   - Hash represents the Hashed value of Data, which can be given as an atom or a variable.
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
//   - sha512: The SHA-512 algorithm.
//   - md5: (insecure) The MD5 algorithm.
//
// Note: Due to the principles of the hash algorithm (pre-image resistance), this predicate can only compute the hash
// value from input data, and cannot compute the original input data from the hash value.
//
// Examples:
//
//	# Compute the SHA-256 hash of the given data and unify it with the given Hash.
//	- crypto_data_hash('Hello OKP4', Hash).
//
//	# Compute the SHA-256 hash of the given hexadecimal data and unify it with the given Hash.
//	- crypto_data_hash('9b038f8ef6918cbb56040dfda401b56b...', Hash, encoding(hex)).
//
//	# Compute the SHA-256 hash of the given hexadecimal data and unify it with the given Hash.
//	- crypto_data_hash([127, ...], Hash, encoding(octet)).
func CryptoDataHash(
	vm *engine.VM, data, hash, options engine.Term, cont engine.Cont, env *engine.Env,
) *engine.Promise {
	functor := "crypto_data_hash/3"
	algorithmOpt := engine.NewAtom("algorithm")

	return engine.Delay(func(ctx context.Context) *engine.Promise {
		algorithmAtom, err := util.GetOptionAsAtomWithDefault(algorithmOpt, options, engine.NewAtom("sha256"), env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: %w", functor, err))
		}
		algorithm, err := util.ParseHashAlg(algorithmAtom.String())
		if err != nil {
			return engine.Error(fmt.Errorf("%s: invalid algorithm: %s. Possible values: %s",
				functor,
				algorithmAtom.String(),
				util.HashAlgNames()))
		}
		decodedData, err := termToBytes(data, options, AtomUtf8, env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode data: %w", functor, err))
		}

		result, err := util.Hash(algorithm, decodedData)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to hash data: %w", functor, err))
		}

		return engine.Unify(vm, hash, util.BytesToCodepointListTermWithDefault(result), cont, env)
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
		typeAtom, err := util.AssertAtom(env, typeTerm)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: %w", functor, err))
		}

		if idx := slices.IndexFunc(algos, func(a util.KeyAlg) bool { return a.String() == typeAtom.String() }); idx == -1 {
			return engine.Error(fmt.Errorf("%s: invalid type: %s. Possible values: %s",
				functor,
				typeAtom.String(),
				strings.Join(util.Map(algos, func(a util.KeyAlg) string { return a.String() }), ", ")))
		}

		decodedKey, err := termToBytes(key, AtomEncoding.Apply(AtomOctet), AtomHex, env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode public key: %w", functor, err))
		}

		decodedData, err := termToBytes(data, options, AtomHex, env)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: failed to decode data: %w", functor, err))
		}

		decodedSignature, err := termToBytes(sig, AtomEncoding.Apply(AtomOctet), AtomHex, env)
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

func termToBytes(term, options, defaultEncoding engine.Term, env *engine.Env) ([]byte, error) {
	encodingTerm, err := util.GetOptionWithDefault(AtomEncoding, options, defaultEncoding, env)
	if err != nil {
		return nil, err
	}
	encodingAtom, err := util.AssertAtom(env, encodingTerm)
	if err != nil {
		return nil, err
	}

	switch encodingAtom {
	case AtomHex:
		return util.TermHexToBytes(term, env)
	case AtomOctet, AtomUtf8:
		return util.StringTermToBytes(term, "", env)
	default:
		return nil, fmt.Errorf("invalid encoding: %s. Possible values: hex, octet", encodingAtom.String())
	}
}
