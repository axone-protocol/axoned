package cli

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"

	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

func CmdQueryPrograms() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "programs",
		Short:   "lists stored programs",
		Example: fmt.Sprintf("$ %s query %s programs --limit 10", version.AppName, types.ModuleName),
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(clientCtx)
			res, err := queryClient.Programs(context.Background(), &types.QueryProgramsRequest{
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "programs")

	return cmd
}

func CmdQueryProgramsByPublisher() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "programs-by-publisher [publisher]",
		Short:   "lists stored programs published by an address",
		Example: fmt.Sprintf("$ %s query %s programs-by-publisher axone1... --limit 10", version.AppName, types.ModuleName),
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			pageReq, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryServiceClient(clientCtx)
			res, err := queryClient.ProgramsByPublisher(context.Background(), &types.QueryProgramsByPublisherRequest{
				Publisher:  args[0],
				Pagination: pageReq,
			})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "programs-by-publisher")

	return cmd
}
