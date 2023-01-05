package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/spf13/cobra"
)

func CmdQueryAsk() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ask [query] [program]",
		Short: "executes a logic query and returns the solutions found.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Executes the [query] for the given [program] file and return the solution(s) found.

Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee
is charged for this, but the execution is constrained by the current limits configured in the module.

Example:
$ %s %s query ask "immortal(X)." program.txt
`,
				version.AppName,
				types.ModuleName,
			),
		),
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			query := args[0]
			program, err := ReadProgramFromFile(args[1])

			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(clientCtx)

			res, err := queryClient.Ask(context.Background(), &types.QueryServiceAskRequest{
				Program: program,
				Query:   query,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// ReadProgramFromFile read text from the given filename. Can pass "-" to read from stdin.
func ReadProgramFromFile(filename string) (v string, err error) {
	var bytes []byte

	if filename == "-" {
		bytes, err = io.ReadAll(os.Stdin)
	} else {
		bytes, err = os.ReadFile(filename)
	}

	if err != nil {
		return
	}

	v = string(bytes)
	return
}
