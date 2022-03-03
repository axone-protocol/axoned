package cli

import (
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/okp4/okp4d/x/knowledge/types"
	"github.com/spf13/cobra"
)

var _ = strconv.Itoa(0)

func CmdBangDataspace() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bang-dataspace [id] [name] [description]",
		Short: "Create a new dataspace given its identifier (unique), name and description.",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argID := args[0]
			argName := args[1]
			argDescription := args[2]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgBangDataspace(
				clientCtx.GetFromAddress().String(),
				argID,
				argName,
				argDescription,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
