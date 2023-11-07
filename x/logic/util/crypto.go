//go:generate go-enum --names
package util

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/sha256"
	"fmt"

	"github.com/dustinxie/ecc"
)

const (
	Secp256k1 KeyAlg = "secp256k1"
	Secp256r1 KeyAlg = "secp256r1"
	Ed25519   KeyAlg = "ed25519"
)

// KeyAlg is the type of key algorithm supported by the crypto util functions.
type KeyAlg string

// String returns the string representation of the key algorithm.
func (a KeyAlg) String() string {
	return string(a)
}

// HashAlg is the type of hash algorithm supported by the crypto util functions.
// ENUM(sha256)
type HashAlg int

// VerifySignature verifies the signature of the given message with the given public key using the given algorithm.
func VerifySignature(alg KeyAlg, pubKey []byte, msg, sig []byte) (_ bool, err error) {
	defer func() {
		if recoveredErr := recover(); recoveredErr != nil {
			err = fmt.Errorf("%s", recoveredErr)
		}
	}()

	switch alg {
	case Ed25519:
		return ed25519.Verify(pubKey, msg, sig), nil
	case Secp256r1:
		return verifySignatureWithCurve(elliptic.P256(), pubKey, msg, sig)
	case Secp256k1:
		return verifySignatureWithCurve(ecc.P256k1(), pubKey, msg, sig)
	default:
		return false, fmt.Errorf("algo %s not supported", alg)
	}
}

// Hash hashes the given data using the given algorithm.
func Hash(alg HashAlg, bytes []byte) ([]byte, error) {
	switch alg {
	case HashAlgSha256:
		hasher := sha256.New()
		hasher.Write(bytes)
		return hasher.Sum(nil), nil
	default:
		return nil, fmt.Errorf("algo %s not supported", alg)
	}
}

// verifySignatureWithCurve verifies the ASN1 signature of the given message with the given
// public key (in compressed form specified in section 4.3.6 of ANSI X9.62.) using the given
// elliptic curve.
func verifySignatureWithCurve(curve elliptic.Curve, pubKey, msg, sig []byte) (bool, error) {
	x, y := ecc.UnmarshalCompressed(curve, pubKey)
	if x == nil || y == nil {
		return false, fmt.Errorf("failed to parse compressed public key (first 10 bytes): %x", pubKey[:10])
	}

	pk := &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	return ecc.VerifyASN1(pk, msg, sig), nil
}
