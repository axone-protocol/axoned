package predicate

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"

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

type Alg string

const (
	Secp256k1 Alg = "secp256k1"
	Secp256r1 Alg = "secp256r1"
	Ed25519   Alg = "ed25519"
)

// ED25519Verify return `true` if the Signature can be verified as the ED25519 signature for Data, using the given PubKey
// as bytes.
//
// ed25519_verify(+PubKey, +Data, +Signature, +Options) is semidet
//
// Where:
// - PubKey is a list of bytes representing the public key.
// - Data is the hash of the signed message could be an Atom or List of bytes.
// - Signature is the signature of the Data, as list of bytes.
// - Options allow to give option to the predicates, available options are:
//   - encoding(+Encoding): Encoding to use for the given Data. Default is `hex`. Can be `hex` or `octet`.
//
// Examples:
//
// # Verify the signature of given hexadecimal data.
// - ed25519_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', [23, 56, ...], encoding(hex)).
//
// # Verify the signature of given binary data.
// - ed25519_verify([127, ...], [56, 90, ..], [23, 56, ...], encoding(octet)).
func ED25519Verify(vm *engine.VM, key, data, sig, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		r, err := cryptoVerify(Ed25519, key, data, sig, options, env)
		if err != nil {
			return engine.Error(fmt.Errorf("ed25519_verify/4: %w", err))
		}

		if !r {
			return engine.Bool(false)
		}
		return cont(env)
	})
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
//   - encoding(+Encoding): Encoding to use for the given Data. Default is `hex`. Can be `hex` or `octet`.
//   - type(+Alg): Alg to use for verify the signature. Default is `secp256r1`. Can be `secp256r1` or `secp256k1`.
//
// Examples:
//
// # Verify the signature of given hexadecimal data as ECDSA secp256r1 algorithm.
// - ecdsa_verify([127, ...], '9b038f8ef6918cbb56040dfda401b56bb1ce79c472e7736e8677758c83367a9d', [23, 56, ...], encoding(hex)).
//
// # Verify the signature of given binary data as ECDSA secp256k1 algorithm.
// - ecdsa_verify([127, ...], [56, 90, ..], [23, 56, ...], [encoding(octet), type(secp256k1)]).
func ECDSAVerify(vm *engine.VM, key, data, sig, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		// TODO: Get good algo in options
		r, err := cryptoVerify(Secp256r1, key, data, sig, options, env)
		if err != nil {
			return engine.Error(fmt.Errorf("ecdsa_verify/4: %w", err))
		}

		if !r {
			return engine.Bool(false)
		}
		return cont(env)
	})
}

func cryptoVerify(alg Alg, key, data, sig, options engine.Term, env *engine.Env) (bool, error) {
	pubKey, err := TermToBytes(key, AtomEncoding.Apply(engine.NewAtom("octet")), env)
	if err != nil {
		return false, fmt.Errorf("decoding public key: %w", err)
	}

	msg, err := TermToBytes(data, options, env)
	if err != nil {
		return false, fmt.Errorf("decoding data: %w", err)
	}

	signature, err := TermToBytes(sig, AtomEncoding.Apply(engine.NewAtom("octet")), env)
	if err != nil {
		return false, fmt.Errorf("decoding signature: %w", err)
	}

	r, err := verifySignature(alg, pubKey, msg, signature)
	if err != nil {
		return false, fmt.Errorf("failed verify signature: %w", err)
	}
	return r, nil
}

func verifySignature(alg Alg, pubKey []byte, msg, sig []byte) (r bool, err error) {
	defer func() {
		if recoveredErr := recover(); recoveredErr != nil {
			err = fmt.Errorf("%s", recoveredErr)
		}
	}()

	switch alg {
	case Ed25519:
		r = ed25519.Verify(pubKey, msg, sig)
	case Secp256r1:
		block, _ := pem.Decode(pubKey)
		if block == nil {
			err = fmt.Errorf("failed decode PEM public key")
			break
		}
		genericPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			break
		}
		pk := genericPublicKey.(*ecdsa.PublicKey)
		r = ecdsa.VerifyASN1(pk, msg, sig)
	case Secp256k1:
		err = fmt.Errorf("secp256k1 public key not implemented yet")
	default:
		err = fmt.Errorf("pub key format not implemented")
	}

	return r, err
}
