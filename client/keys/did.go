package keys

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/okp4/okp4d/x/logic/util"
)

var (
	flagPubKeyType = "type"
)

func DIDCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   fmt.Sprintf("did [pubkey] -t [{%s, %s}]", util.KeyAlgEd25519, util.KeyAlgSecp256k1),
		Short: fmt.Sprintf("Give the did:key from a %s or %s pubkey (hex, base64)", util.KeyAlgEd25519, util.KeyAlgSecp256k1),
		Long: fmt.Sprintf(`Give the did:key from a %s or %s pubkey given as hex or base64 encoded string.

Example:
$ %s keys did "AtD+mbIUqu615Grk1loWI6ldnQzs1X1nP35MmhmsB1K8" -t %s
$ %s keys did 02d0fe99b214aaeeb5e46ae4d65a1623a95d9d0cecd57d673f7e4c9a19ac0752bc -t %s
			`, util.KeyAlgEd25519, util.KeyAlgSecp256k1, version.AppName, util.KeyAlgSecp256k1, version.AppName, util.KeyAlgSecp256k1),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			pubkeyType, err := cmd.Flags().GetString(flagPubKeyType)
			if err != nil {
				return err
			}
			pubkeyAlgo, err := util.ParseKeyAlg(pubkeyType)
			if err != nil {
				return errorsmod.Wrapf(errors.ErrInvalidType,
					"invalid pubkey type; expected oneof %+q", []util.KeyAlg{util.KeyAlgSecp256k1, util.KeyAlgEd25519})
			}
			bs, err := getBytesFromString(args[0])
			if err != nil {
				return err
			}
			pubKey, err := bytesToPubkey(bs, pubkeyAlgo)
			if err != nil {
				return err
			}
			did, err := util.CreateDIDKeyByPubKey(pubKey)
			if err != nil {
				return errorsmod.Wrapf(errors.ErrInvalidPubKey, "failed to make did:key from %s; %s", args[0], err)
			}

			cmd.Println(did)

			return nil
		},
	}
	cmd.Flags().StringP(flagPubKeyType, "t", util.KeyAlgSecp256r1.String(),
		fmt.Sprintf("Pubkey type to decode (oneof %s, %s)", util.KeyAlgEd25519, util.KeyAlgSecp256k1))
	return cmd
}

func getBytesFromString(pubKey string) ([]byte, error) {
	if bz, err := hex.DecodeString(pubKey); err == nil {
		return bz, nil
	}

	if bz, err := base64.StdEncoding.DecodeString(pubKey); err == nil {
		return bz, nil
	}

	return nil, errorsmod.Wrapf(errors.ErrInvalidPubKey,
		"pubkey '%s' invalid; expected hex or base64 encoding of correct size", pubKey)
}

func bytesToPubkey(bz []byte, keytype util.KeyAlg) (cryptotypes.PubKey, error) {
	switch keytype {
	case util.KeyAlgEd25519:
		if len(bz) != ed25519.PubKeySize {
			return nil, errorsmod.Wrapf(errors.ErrInvalidPubKey,
				"invalid pubkey size; expected %d, got %d", ed25519.PubKeySize, len(bz))
		}
		return &ed25519.PubKey{Key: bz}, nil
	case util.KeyAlgSecp256k1:
		if len(bz) != secp256k1.PubKeySize {
			return nil, errorsmod.Wrapf(errors.ErrInvalidPubKey,
				"invalid pubkey size; expected %d, got %d", secp256k1.PubKeySize, len(bz))
		}
		return &secp256k1.PubKey{Key: bz}, nil
	case util.KeyAlgSecp256r1:
	}

	return nil, errorsmod.Wrapf(errors.ErrInvalidType,
		"invalid pubkey type; expected oneof %+q", []util.KeyAlg{util.KeyAlgSecp256k1, util.KeyAlgEd25519})
}
