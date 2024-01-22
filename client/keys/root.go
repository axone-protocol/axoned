package keys

import (
	"github.com/spf13/cobra"
)

func Install(cmd *cobra.Command) *cobra.Command {
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
