package util

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
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
		if !ecc.VerifyBytes(pk, msg, sig, ecc.Normal) {
			return false, nil
		}
	default:
		err = fmt.Errorf("algo %s not supported", alg)
	}

	return r, err
}
