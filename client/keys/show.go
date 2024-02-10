package keys

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	errorsmod "cosmossdk.io/errors"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/ledger"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerr "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	flagMultiSigThreshold = "multisig-threshold"
)

func runShowCmd(cmd *cobra.Command, args []string) error {
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return err
	}

	var keyRecord *keyring.Record
	if len(args) == 1 {
		keyRecord, err = fetchKey(clientCtx.Keyring, args[0])
		if err != nil {
			return fmt.Errorf("%s is not a valid name or address: %w", args[0], err)
		}
	} else {
		keyRecord, err = fetchMultiSigKey(cmd, clientCtx, args)
		if err != nil {
			return err
		}
	}

	if err := checkFlagCompatibility(cmd); err != nil {
		return err
	}

	bechPrefix, _ := cmd.Flags().GetString(keys.FlagBechPrefix)
	bechKeyOut, err := getBechKeyOut(bechPrefix)
	if err != nil {
		return err
	}

	if err := processOutput(cmd, clientCtx, keyRecord, bechKeyOut); err != nil {
		return err
	}

	return nil
}

func fetchMultiSigKey(cmd *cobra.Command, clientCtx client.Context, args []string) (*keyring.Record, error) {
	pks := make([]cryptotypes.PubKey, len(args))
	for i, keyRef := range args {
		k, err := fetchKey(clientCtx.Keyring, keyRef)
		if err != nil {
			return nil, fmt.Errorf("%s is not a valid name or address: %w", keyRef, err)
		}
		pubKey, err := k.GetPubKey()
		if err != nil {
			return nil, err
		}
		pks[i] = pubKey
	}

	multisigThreshold, _ := cmd.Flags().GetInt(flagMultiSigThreshold)
	if err := validateMultisigThreshold(multisigThreshold, len(args)); err != nil {
		return nil, err
	}

	multiKey := multisig.NewLegacyAminoPubKey(multisigThreshold, pks)
	return keyring.NewMultiRecord(args[0], multiKey)
}

func checkFlagCompatibility(cmd *cobra.Command) error {
	isShowAddr, _ := cmd.Flags().GetBool(keys.FlagAddress)
	isShowPubKey, _ := cmd.Flags().GetBool(keys.FlagPublicKey)
	if isShowAddr && isShowPubKey {
		return errors.New("cannot use both --address and --pubkey at once")
	}

	isOutputSet := cmd.Flag(flags.FlagOutput) != nil && cmd.Flag(flags.FlagOutput).Changed
	if isOutputSet && (isShowAddr || isShowPubKey) {
		return errors.New("cannot use --output with --address or --pubkey")
	}

	return nil
}

func processOutput(cmd *cobra.Command, clientCtx client.Context, k *keyring.Record, bechKeyOut bechKeyOutFn) error {
	isShowAddr, _ := cmd.Flags().GetBool(keys.FlagAddress)
	isShowPubKey, _ := cmd.Flags().GetBool(keys.FlagPublicKey)
	isShowDevice, _ := cmd.Flags().GetBool(keys.FlagDevice)

	if isShowDevice {
		return handleDeviceOutput(k)
	}

	if isShowAddr || isShowPubKey {
		ko, err := bechKeyOut(k)
		if err != nil {
			return err
		}
		out := ko.Address
		if isShowPubKey {
			out = ko.PubKey
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), out)
		return err
	}

	outputFormat := clientCtx.OutputFormat
	return printKeyringRecord(cmd.OutOrStdout(), k, bechKeyOut, outputFormat)
}

func handleDeviceOutput(k *keyring.Record) error {
	if k.GetType() != keyring.TypeLedger {
		return fmt.Errorf("the device flag (-d) can only be used for ledger keys")
	}
	ledgerItem := k.GetLedger()
	if ledgerItem == nil {
		return errors.New("unable to get ledger item")
	}
	pk, err := k.GetPubKey()
	if err != nil {
		return err
	}
	bechPrefix := sdk.GetConfig().GetBech32AccountAddrPrefix()
	return ledger.ShowAddress(*ledgerItem.Path, pk, bechPrefix)
}

func fetchKey(kb keyring.Keyring, keyref string) (*keyring.Record, error) {
	k, err := kb.Key(keyref)

	if err == nil || !errorsmod.IsOf(err, sdkerr.ErrIO, sdkerr.ErrKeyNotFound) {
		return k, err
	}

	accAddr, err := sdk.AccAddressFromBech32(keyref)
	if err != nil {
		return k, err
	}

	k, err = kb.KeyByAddress(accAddr)
	return k, errorsmod.Wrap(err, "Invalid key")
}

func validateMultisigThreshold(k, nKeys int) error {
	if k <= 0 {
		return fmt.Errorf("threshold must be a positive integer")
	}
	if nKeys < k {
		return fmt.Errorf(
			"threshold k of n multisignature: %d < %d", nKeys, k)
	}
	return nil
}

func getBechKeyOut(bechPrefix string) (bechKeyOutFn, error) {
	switch bechPrefix {
	case sdk.PrefixAccount:
		return keys.MkAccKeyOutput, nil
	case sdk.PrefixValidator:
		return keys.MkValKeyOutput, nil
	case sdk.PrefixConsensus:
		return keys.MkConsKeyOutput, nil
	}

	return nil, fmt.Errorf("invalid Bech32 prefix encoding provided: %s", bechPrefix)
}
