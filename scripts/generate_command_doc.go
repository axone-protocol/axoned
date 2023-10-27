package main

import (
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/okp4/okp4d/cmd/okp4d/cmd"
)

func GenerateCommandDocumentation() error {
	if err := os.Setenv("DAEMON_NAME", "okp4d"); err != nil {
		return err
	}

	targetPath := "docs/command"
	rootCmd, _ := cmd.NewRootCmd()
	rootCmd.DisableAutoGenTag = true

	err := os.Mkdir(targetPath, 0o750)
	if err != nil && !os.IsExist(err) {
		return err
	}

	return doc.GenMarkdownTree(rootCmd, targetPath)
}
