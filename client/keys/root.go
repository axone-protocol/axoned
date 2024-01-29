package keys

import (
	"github.com/spf13/cobra"
)

// Enhance augment the given command which is assumed to be the root command of the 'list' command.
// It will:
// - add the 'did' command.
// - replace the original 'list' command with our own 'list' command which will list all did:key.
func Enhance(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(
		DIDCmd(),
	)

	for i, c := range cmd.Commands() {
		if c.Name() == ListCmdName {
			cmd.RemoveCommand(cmd.Commands()[i])
			cmd.AddCommand(ListKeysCmd())
			break
		}
	}

	return cmd
}
