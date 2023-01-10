package cli

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/spf13/cobra"
)

var (
	program     string
	programFile string
)

func CmdQueryAsk() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ask [query]",
		Short: "executes a logic query and returns the solutions found.",
		Long: `Executes the [query] for the given [program] file and return the solution(s) found.

Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee
is charged for this, but the execution is constrained by the current limits configured in the module (that you can
query).`,
		Example: fmt.Sprintf(`$ %s %s query ask "chain_id(X)." # returns the chain-id`,
			version.AppName,
			types.ModuleName),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			if programFile != "" {
				program, err = ReadProgramFromFile(programFile)
				if err != nil {
					return
				}
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			query := args[0]
			queryClient := types.NewQueryServiceClient(clientCtx)

			res, err := queryClient.Ask(context.Background(), &types.QueryServiceAskRequest{
				Program: program,
				Query:   query,
			})
			if err != nil {
				return
			}

			return clientCtx.PrintProto(res)
		},
	}

	cmd.Flags().StringVar(
		&program,
		"program",
		"",
		`reads the program from the given filename or from stdin if "-" is passed as the filename.`)
	cmd.Flags().StringVar(
		&programFile,
		"program-file",
		"",
		`reads the program from the given string.`)
	cmd.MarkFlagsMutuallyExclusive("program", "program-file")

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// ReadProgramFromFile reads text from the given filename or from stdin if "-" is passed as the filename.
// It returns the text as a string and an error value.
func ReadProgramFromFile(filename string) (string, error) {
	var r io.Reader

	if filename == "-" {
		r = os.Stdin
	} else {
		file, err := os.Open(filename)
		if err != nil {
			return "", err
		}
		defer func() {
			_ = file.Close()
		}()

		r = file
	}

	bytes, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
