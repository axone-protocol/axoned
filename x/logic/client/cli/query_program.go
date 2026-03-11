package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

func CmdQueryProgram() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program [program-id]",
		Short: "shows stored program metadata",
		Example: fmt.Sprintf(
			"$ %s query %s program 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			version.AppName,
			types.ModuleName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryServiceClient(clientCtx)

			res, err := queryClient.Program(context.Background(), &types.QueryProgramRequest{
				ProgramId: args[0],
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

func CmdQueryProgramSource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "program-source [program-id]",
		Short: "shows stored program source",
		Example: fmt.Sprintf(
			"$ %s query %s program-source 0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			version.AppName,
			types.ModuleName,
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			queryClient := types.NewQueryServiceClient(clientCtx)

			res, err := queryClient.ProgramSource(context.Background(), &types.QueryProgramSourceRequest{
				ProgramId: args[0],
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
