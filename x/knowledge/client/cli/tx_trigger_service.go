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

func CmdTriggerService() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "trigger-service [uri]",
		Short: "Trigger a service execution from an invocation URI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			argURI := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgTriggerService(
				clientCtx.GetFromAddress().String(),
				argURI,
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
