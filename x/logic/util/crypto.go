//go:generate go-enum --names
package util

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/md5" //nolint:gosec
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"

	"github.com/dustinxie/ecc"
)

// KeyAlg is the type of key algorithm supported by the crypto util functions.
// ENUM(secp256k1,secp256r1,ed25519).
type KeyAlg int

// HashAlg is the type of hash algorithm supported by the crypto util functions.
// ENUM(md5,sha256,sha512).
type HashAlg int

// Hasher returns a new hash.Hash for the given algorithm.
func (a HashAlg) Hasher() (hash.Hash, error) {
	switch a {
	case HashAlgMd5:
		//nolint:gosec // md5 is used for hashing, not for cryptography
		return md5.New(), nil
	case HashAlgSha256:
		return sha256.New(), nil
	case HashAlgSha512:
		return sha512.New(), nil
	default:
		return nil, fmt.Errorf("algo %s not supported", a.String())
	}
}

// VerifySignature verifies the signature of the given message with the given public key using the given algorithm.
func VerifySignature(alg KeyAlg, pubKey []byte, msg, sig []byte) (_ bool, err error) {
	defer func() {
		if recoveredErr := recover(); recoveredErr != nil {
			err = fmt.Errorf("%s", recoveredErr)
		}
	}()

	switch alg {
	case KeyAlgEd25519:
		return ed25519.Verify(pubKey, msg, sig), nil
	case KeyAlgSecp256r1:
		return verifySignatureWithCurve(elliptic.P256(), pubKey, msg, sig)
	case KeyAlgSecp256k1:
		return verifySignatureWithCurve(ecc.P256k1(), pubKey, msg, sig)
	default:
		return false, fmt.Errorf("algo %s not supported", alg)
	}
}

// Hash hashes the given data using the given algorithm.
func Hash(alg HashAlg, bytes []byte) ([]byte, error) {
	hasher, err := alg.Hasher()
	if err != nil {
		return nil, err
	}

	hasher.Write(bytes)
	return hasher.Sum(nil), nil
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
