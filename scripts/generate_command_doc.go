package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra/doc"

	"github.com/okp4/okp4d/cmd/okp4d/cmd"
)

func main() {
	err := generateDocumentaton("command")
	if err != nil {
		fmt.Printf("failed to generate documentation: %e", err)

		var codeErr *server.ErrorCode
		switch {
		case errors.As(err, &codeErr):
			os.Exit(codeErr.Code)
		default:
			os.Exit(1)
		}
	}
}

func generateDocumentaton(folder string) error {
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.DisableAutoGenTag = true

	err := os.Mkdir(folder, 0750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	err = doc.GenMarkdownTree(rootCmd, "command")
	if err != nil {
		return err
	}

	return nil
}
