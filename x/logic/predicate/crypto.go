package predicate

import (
	"context"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"encoding/hex"
	"fmt"
	"math/big"

	cometcrypto "github.com/cometbft/cometbft/crypto"
	"github.com/ichiban/prolog/engine"

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

func ECDSAVerify(vm *engine.VM, key, data, sig, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		pubKey, err := TermToBytes(key, env)
		if err != nil {
			return engine.Error(fmt.Errorf("ecdsa_verify/4: decoding public key: %w", err))
		}

		msg, err := TermToBytes(data, env)
		if err != nil {
			return engine.Error(fmt.Errorf("ecdsa_verify/4: decoding data: %w", err))
		}

		signature, err := TermToBytes(sig, env)
		if err != nil {
			return engine.Error(fmt.Errorf("ecdsa_verify/4: decoding signature: %w", err))
		}

		// TODO: Create function hasDecoding option
		r, err := verifySignature(Ed25519, pubKey, msg, signature)
		if err != nil {
			return engine.Error(fmt.Errorf("ecdsa_verify/4: failed verify signature: %w", err))
		}

		if !r {
			return engine.Bool(false)
		}
		return cont(env)
	})
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
		//pub, err := secp256r1BytesToPublicKey(pubKey)
		//
		//var pub secp256r1.PubKey
		//pub = pubKey
		//secp256r1.PubKey{Key: pubKey}
		//if key, ok :=  pubKey.(); ok {
		//	return key.VerifySignature(msg, sig), nil
		//}
		err = fmt.Errorf("public key is not secp256r1 compatible")
	case Secp256k1:
		//if key, ok := pubKey.(secp256k1.PubKey); ok {
		//	return key.VerifySignature(msg, sig), nil
		//}
		err = fmt.Errorf("public key is not secp256k1 compatible")
	default:
		err = fmt.Errorf("pub key format not implemented")
	}

	return r, err
}

func secp256r1BytesToPublicKey(pubKeyBytes []byte) (*ecdsa.PublicKey, error) {
	curve := elliptic.P256()

	// La longueur de la séquence de bytes doit être 2 * la longueur des points de la courbe elliptique + 1 (pour le byte de préfixe non compressé)
	if len(pubKeyBytes) != 2*curve.Params().BitSize/8+1 {
		return nil, fmt.Errorf("public key size is not good: %d expected, got %d", 2*curve.Params().BitSize/8+1, len(pubKeyBytes))
	}

	xBytes := pubKeyBytes[1 : curve.Params().BitSize/8+1]
	yBytes := pubKeyBytes[curve.Params().BitSize/8+1:]

	x := new(big.Int).SetBytes(xBytes)
	y := new(big.Int).SetBytes(yBytes)

	publicKey := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	return publicKey, nil
}
