package util

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"fmt"

	"github.com/dustinxie/ecc"
)

// Alg is the type of algorithm supported by the crypto util functions.
type Alg string

// String returns the string representation of the algorithm.
func (a Alg) String() string {
	return string(a)
}

const (
	Secp256k1 Alg = "secp256k1"
	Secp256r1 Alg = "secp256r1"
	Ed25519   Alg = "ed25519"
)

// VerifySignature verifies the signature of the given message with the given public key using the given algorithm.
func VerifySignature(alg Alg, pubKey []byte, msg, sig []byte) (_ bool, err error) {
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
