package keys

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/keys"
)

const (
	flagListNames = "list-names"
)

// KeyOutput is the output format for keys when listing them.
// It is an improved copy of the KeyOutput from the keys module (github.com/cosmos/cosmos-sdk/client/keys/types.go).
type KeyOutput struct {
	keys.KeyOutput
	DID string `json:"did,omitempty" yaml:"did"`
}

// runListCmd retrieves all keys from the keyring and prints them to the console.
// This is an improved copy of the runListCmd from the keys module of the cosmos-sdk.
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
