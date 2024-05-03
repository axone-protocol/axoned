package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	sdkmath "cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/axone-protocol/axoned/v7/x/logic/types"
)

var (
	program string
	limit   uint64
)

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

			limit := sdkmath.NewUint(limit)
			res, err := queryClient.Ask(context.Background(), &types.QueryServiceAskRequest{
				Program: program,
				Query:   query,
				Limit:   &limit,
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
	//nolint:lll
	cmd.Flags().Uint64Var(
		&limit,
		"limit",
		1,
		`limit the maximum number of solutions to return.
This parameter is constrained by the 'max_result_count' setting in the module configuration, which specifies the maximum number of results that can be requested per query.`)

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
