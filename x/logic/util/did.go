package util

import (
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

const (
	ED25519PubKeyMultiCodec   = 0xed
	SECP256k1PubKeyMultiCodec = 0xe7
)

// BytesToPubKey converts bytes to a PubKey given a key type.
// Supported key types: secp256k1, ed25519.
func BytesToPubKey(bz []byte, keytype KeyAlg) (cryptotypes.PubKey, error) {
	invalidPubKey := func(expectedSize int, bs []byte) error {
		return fmt.Errorf("invalid pubkey size; expected %d, got %d", expectedSize, len(bs))
	}
	switch keytype {
	case KeyAlgEd25519:
		if len(bz) != ed25519.PubKeySize {
			return nil, invalidPubKey(ed25519.PubKeySize, bz)
		}
		return &ed25519.PubKey{Key: bz}, nil
	case KeyAlgSecp256k1:
		if len(bz) != secp256k1.PubKeySize {
			return nil, invalidPubKey(secp256k1.PubKeySize, bz)
		}
		return &secp256k1.PubKey{Key: bz}, nil
	case KeyAlgSecp256r1:
	}

	return nil, fmt.Errorf("invalid pubkey type: %s; expected oneof %+q",
		keytype, []KeyAlg{KeyAlgSecp256k1, KeyAlgEd25519})
}

// CreateDIDKeyByPubKey creates a did:key ID using the given public key.
// The multicodec key fingerprint is determined by the key type and complies with the did:key format spec found at:
// https://w3c-ccg.github.io/did-method-key/#format.
func CreateDIDKeyByPubKey(pubKey cryptotypes.PubKey) (string, error) {
	code, err := multicodecFromPubKey(pubKey)
	if err != nil {
		return "", err
	}

	didKey, _ := fingerprint.CreateDIDKeyByCode(code, pubKey.Bytes())
	return didKey, nil
}

// CreateDIDKeyIDByPubKey creates a DID key ID using the given public key.
// The function is similar to CreateDIDKeyByPubKey but returns the DID with the key ID as hash fragment.
// The multicodec key fingerprint is determined by the key type and complies with the did:key format spec found at:
// https://w3c-ccg.github.io/did-method-key/#format.
func CreateDIDKeyIDByPubKey(pubKey cryptotypes.PubKey) (string, error) {
	code, err := multicodecFromPubKey(pubKey)
	if err != nil {
		return "", err
	}

	_, keyID := fingerprint.CreateDIDKeyByCode(code, pubKey.Bytes())
	return keyID, nil
}

// multicodecFromPubKey returns the multicodec and the error message for the given public key.
// Supported key types: secp256k1, ed25519.
func multicodecFromPubKey(pubKey cryptotypes.PubKey) (uint64, error) {
	var code uint64
	switch pubKey.(type) {
	case *ed25519.PubKey:
		code = ED25519PubKeyMultiCodec
	case *secp256k1.PubKey:
		code = SECP256k1PubKeyMultiCodec
	default:
		return 0, fmt.Errorf("invalid pubkey type: %s; expected oneof %+q",
			pubKey.Type(), []string{(&ed25519.PubKey{}).Type(), (&secp256k1.PubKey{}).Type()})
	}
	return code, nil
}
