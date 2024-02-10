package keys

import (
	"github.com/spf13/cobra"
)

const (
	listKeysCmd = "list"
	showKeysCmd = "show"
)

// Enhance augment the given command which is assumed to be the root command of the 'list' command.
// It will:
// - add the 'did' command.
// - replace the original 'list' command implementation with our own 'list' command which will list all did:key.
// - replace the original 'show' command implementation with our own 'show' command which will show the did:key of the key.
func Enhance(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(
		DIDCmd(),
	)

	for _, c := range cmd.Commands() {
		switch c.Name() {
		case listKeysCmd:
			c.RunE = runListCmd
		case showKeysCmd:
			c.RunE = runShowCmd
		default:
		}
	}

	return cmd
}
