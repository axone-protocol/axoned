package main

import (
	"errors"
	"os"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/server"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "gendoc",
		Short: "Simple CLI to generate documentation for the project",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SetOut(cmd.OutOrStdout())
			cmd.SetErr(cmd.ErrOrStderr())
		},
	}

	rootCmd.AddCommand(generateCommandDocumentationCommand())
	rootCmd.AddCommand(generatePredicateDocumentationCommand())

	if err := rootCmd.Execute(); err != nil {
		var codeErr *server.ErrorCode
		switch {
		case errors.As(err, &codeErr):
			os.Exit(codeErr.Code)
		default:
			os.Exit(1)
		}
	}
}

func generateCommandDocumentationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "command",
		Short: "Generate command documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return GenerateCommandDocumentation()
		},
	}
	return cmd
}

func generatePredicateDocumentationCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "predicate",
		Short: "Generate predicate documentation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return GeneratePredicateDocumentation()
		},
	}
	return cmd
}
