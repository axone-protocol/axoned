package main

import (
	"errors"
	"os"

	"github.com/cosmos/cosmos-sdk/server"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	"github.com/okp4/okp4d/cmd/okp4d/cmd"

	"github.com/okp4/okp4d/app"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, "", app.DefaultNodeHome); err != nil {
		var codeErr *server.ErrorCode
		switch {
		case errors.As(err, &codeErr):
			os.Exit(codeErr.Code)
		default:
			os.Exit(1)
		}
	}
}
