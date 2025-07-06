package keys

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/axone-protocol/axoned/v12/x/logic/util"
)

var flagPubKeyType = "type"

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
		RunE: runDidCmd(),
	}
	cmd.Flags().StringP(flagPubKeyType, "t", util.KeyAlgSecp256k1.String(),
		fmt.Sprintf("Pubkey type to decode (oneof %s, %s)", util.KeyAlgEd25519, util.KeyAlgSecp256k1))
	return cmd
}

func runDidCmd() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
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
		pubKey, err := util.BytesToPubKey(bs, pubkeyAlgo)
		if err != nil {
			return errorsmod.Wrapf(errors.ErrInvalidPubKey, "failed to make pubkey from %s; %s", args[0], err)
		}
		did, err := util.CreateDIDKeyByPubKey(pubKey)
		if err != nil {
			return errorsmod.Wrapf(errors.ErrInvalidPubKey, "failed to make did:key from %s; %s", args[0], err)
		}

		cmd.Println(did)

		return nil
	}
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
