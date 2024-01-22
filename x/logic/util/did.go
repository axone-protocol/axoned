package util

import (
	"fmt"

	"github.com/hyperledger/aries-framework-go/pkg/vdr/fingerprint"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

const (
	ED25519PubKeyMultiCodec  = 0xed
	SECP256k1ubKeyMultiCodec = 0xe7
)

// CreateDIDKeyByPubKey creates a did:key ID using the given public key.
// The multicodec key fingerprint is determined by the key type and complies with the did:key format spec found at:
// https://w3c-ccg.github.io/did-method-key/#format.
func CreateDIDKeyByPubKey(pubKey cryptotypes.PubKey) (string, error) {
	var code uint64
	switch pubKey.(type) {
	case *ed25519.PubKey:
		code = ED25519PubKeyMultiCodec
	case *secp256k1.PubKey:
		code = SECP256k1ubKeyMultiCodec
	default:
		return "", fmt.Errorf("unsupported key type: %s", pubKey.Type())
	}

	did, _ := fingerprint.CreateDIDKeyByCode(code, pubKey.Bytes())
	return did, nil
}
