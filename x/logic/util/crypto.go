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
func VerifySignature(alg Alg, pubKey []byte, msg, sig []byte) (r bool, err error) {
	defer func() {
		if recoveredErr := recover(); recoveredErr != nil {
			err = fmt.Errorf("%s", recoveredErr)
		}
	}()

	switch alg {
	case Ed25519:
		r = ed25519.Verify(pubKey, msg, sig)
	case Secp256r1:
		curve := elliptic.P256()
		x, y := ecc.UnmarshalCompressed(curve, pubKey)
		if x == nil || y == nil {
			err = fmt.Errorf("failed to parse compressed public key")
			break
		}

		pk := &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		}

		r = ecc.VerifyASN1(pk, msg, sig)
	case Secp256k1:
		curve := ecc.P256k1()
		x, y := ecc.UnmarshalCompressed(curve, pubKey)
		if x == nil || y == nil {
			err = fmt.Errorf("failed to parse compressed public key")
			break
		}

		pk := &ecdsa.PublicKey{
			Curve: curve,
			X:     x,
			Y:     y,
		}

		r = ecc.VerifyASN1(pk, msg, sig)
	default:
		err = fmt.Errorf("algo %s not supported", alg)
	}

	return r, err
}
