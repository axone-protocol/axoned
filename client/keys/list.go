package keys

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
)

const (
	flagListNames = "list-names"
	listKeysCmd   = "list"
)

// EnhanceListCmd replaces the original 'list' command implementation with our own 'list' command which
// will allow us to list did:key of the keys as well as the original keys.
func EnhanceListCmd(cmd *cobra.Command) {
	for _, c := range cmd.Commands() {
		if c.Name() == listKeysCmd {
			c.RunE = runListCmd
			break
		}
	}
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
