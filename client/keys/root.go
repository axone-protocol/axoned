package keys

import "github.com/spf13/cobra"

func Enhance(cmd *cobra.Command) *cobra.Command {
	cmd.AddCommand(
		DIDCmd(),
	)
	return cmd
}
