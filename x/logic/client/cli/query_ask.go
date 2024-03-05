package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/okp4/okp4d/v7/x/logic/types"
)

var program string

func CmdQueryAsk() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ask [query]",
		Short: "executes a logic query and returns the solutions found.",
		Long: `Executes the [query] and return the solution(s) found.
 Optionally, a program can be transmitted, which will be interpreted before the query is processed.
 Since the query is without any side-effect, the query is not executed in the context of a transaction and no fee
 is charged for this, but the execution is constrained by the current limits configured in the module (that you can
 query).`,
		Example: fmt.Sprintf(`$ %s query %s ask "chain_id(X)." # returns the chain-id`,
			version.AppName,
			types.ModuleName),
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
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
		`reads the program from the given string.`)

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
