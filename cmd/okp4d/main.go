package main

import (
	_ "embed"
	"os"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/ignite/cli/ignite/pkg/cosmoscmd"
	"github.com/ignite/cli/ignite/pkg/xstrings"
	"github.com/okp4/okp4d/app"
)

func main() {
	rootCmd, _ := cosmoscmd.NewRootCmd(
		app.Name,
		app.AccountAddressPrefix,
		app.DefaultNodeHome,
		xstrings.NoDash(app.Name),
		app.ModuleBasics,
		app.New,
	)

	rootCmd.Use = app.Name
	rootCmd.Short = Resource.Short
	rootCmd.Long = Resource.Long

	for _, element := range rootCmd.Commands() {
		if element.Name() == "add-genesis-account" {
			rootCmd.RemoveCommand(element)
		}
	}
	rootCmd.AddCommand(AddGenesisAccountCmd(app.DefaultNodeHome))

	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
