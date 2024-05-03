package main

import (
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/axone/axoned/v7/cmd/axoned/cmd"
)

func generateCommandDocumentation() error {
	if err := os.Setenv("DAEMON_NAME", "axoned"); err != nil {
		return err
	}

	targetPath := "docs/command"
	rootCmd := cmd.NewRootCmd()
	rootCmd.DisableAutoGenTag = true

	err := os.Mkdir(targetPath, 0o750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	err = doc.GenMarkdownTree(rootCmd, targetPath)
	if err != nil {
		return err
	}

	return normalizeMarkdownFiles(targetPath)
}
