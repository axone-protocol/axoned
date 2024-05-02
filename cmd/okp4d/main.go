package main

import (
	"os"

	"cosmossdk.io/log"

	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"

	"github.com/axone/axoned/v7/app"
	"github.com/axone/axoned/v7/cmd/okp4d/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		log.NewLogger(rootCmd.OutOrStderr()).Error("failure when running app", "err", err)
		os.Exit(1)
	}
}
