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
		Use:   "bang-dataspace [name] [description]",
		Short: "Broadcast message bang-dataspace",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argName := args[0]
			argDescription := args[1]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgBangDataspace(
				clientCtx.GetFromAddress().String(),
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
