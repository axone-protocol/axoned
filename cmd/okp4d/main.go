package main

import (
	"fmt"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/okp4/okp4d/app"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd, _ := NewRootCmd()

	if err := Extend(rootCmd); err != nil {
		os.Exit(1)
	}

	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}

func Extend(cmd *cobra.Command) error {
	cmdGetter := func(name string) (*cobra.Command, error) {
		return getSubCommand(cmd, name)
	}

	return ExtendDebugCmd(cmdGetter)
}

func getSubCommand(cmd *cobra.Command, name string) (*cobra.Command, error) {
	for i, v := range cmd.Commands() {
		if v.Name() == name {
			return cmd.Commands()[i], nil
		}
	}

	return nil, fmt.Errorf("cannot find '%s' command", name)
}
