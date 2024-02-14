package credential

import (
	"github.com/spf13/cobra"
)

// Commands registers a sub-tree of commands to manage Verifiable Credentials.
func Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credential",
		Short: "W3C Verifiable Credential",
		Long: `W3C Verifiable Credentials management commands.

This command provides a set of sub-commands to manage W3C Verifiable Credentials, including signing and verification operations.`,
	}

	cmd.AddCommand(
		SignCmd(),
	)

	return cmd
}
