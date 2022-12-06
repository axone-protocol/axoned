package main

import (
	"errors"
	"log"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra/doc"

	"github.com/okp4/okp4d/cmd/okp4d/cmd"
)

func main() {
	err := generateDocumentation("command")
	if err != nil {
		log.Printf("failed to generate documentation: %s\n", err)

		var codeErr *server.ErrorCode
		switch {
		case errors.As(err, &codeErr):
			os.Exit(codeErr.Code)
		default:
			os.Exit(1)
		}
	}
}

func generateDocumentation(folder string) error {
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.DisableAutoGenTag = true

	err := os.Mkdir(folder, 0o750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	err = doc.GenMarkdownTree(rootCmd, "command")
	if err != nil {
		return err
	}

	return nil
}
