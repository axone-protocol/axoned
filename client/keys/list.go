package keys

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
	cryptokeyring "github.com/cosmos/cosmos-sdk/crypto/keyring"

	"github.com/okp4/okp4d/x/logic/util"
)

const (
	flagListNames = "list-names"
	ListCmdName   = "list"
)

// KeyOutput defines a structure wrapping around an Info object used for output
// functionality.
type KeyOutput struct {
	keys.KeyOutput
	DID string `json:"did,omitempty" yaml:"did"`
}

// ListKeysCmd lists all keys in the key store with additional info, such as the did:key equivalent of the public key.
// This is an improved copy of the ListKeysCmd from the keys module.
func ListKeysCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   ListCmdName,
		Short: "List all keys",
		Long: `Return a list of all public keys stored by this key manager
along with their associated name, address and decentralized identifier (for supported public key algorithms)`,
		RunE: runListCmd,
	}

	cmd.Flags().BoolP(flagListNames, "n", false, "List names only")
	return cmd
}

func runListCmd(cmd *cobra.Command, _ []string) error {
	clientCtx, err := client.GetClientQueryContext(cmd)
	if err != nil {
		return err
	}

	records, err := clientCtx.Keyring.List()
	if err != nil {
		return err
	}

	if len(records) == 0 && clientCtx.OutputFormat == flags.OutputFormatText {
		cmd.Println("No records were found in keyring")
		return nil
	}

	if ok, _ := cmd.Flags().GetBool(flagListNames); !ok {
		return printKeyringRecords(cmd.OutOrStdout(), records, clientCtx.OutputFormat)
	}

	for _, k := range records {
		cmd.Println(k.Name)
	}

	return nil
}

func printKeyringRecords(w io.Writer, records []*cryptokeyring.Record, output string) error {
	kos, err := mkKeyOutput(records)
	if err != nil {
		return err
	}

	switch output {
	case flags.OutputFormatText:
		if err := printTextRecords(w, kos); err != nil {
			return err
		}

	case flags.OutputFormatJSON:
		out, err := json.Marshal(kos)
		if err != nil {
			return err
		}

		if _, err := fmt.Fprintf(w, "%s", out); err != nil {
			return err
		}
	}

	return nil
}

func printTextRecords(w io.Writer, kos []KeyOutput) error {
	out, err := yaml.Marshal(&kos)
	if err != nil {
		return err
	}

	if _, err := fmt.Fprintln(w, string(out)); err != nil {
		return err
	}

	return nil
}

func mkKeyOutput(records []*cryptokeyring.Record) ([]KeyOutput, error) {
	kos := make([]KeyOutput, len(records))

	for i, r := range records {
		kko, err := keys.MkAccKeyOutput(r)
		if err != nil {
			return nil, err
		}
		pk, err := r.GetPubKey()
		if err != nil {
			return nil, err
		}
		did, _ := util.CreateDIDKeyByPubKey(pk)

		kos[i] = KeyOutput{
			KeyOutput: kko,
			DID:       did,
		}
	}

	return kos, nil
}
