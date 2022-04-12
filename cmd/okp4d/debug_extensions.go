package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/proto/tendermint/types"
)

func ExtendDebugCmd(getCmd func(name string) (*cobra.Command, error)) error {
	debugCmd, err := getCmd("debug")
	if err != nil {
		return err
	}

	debugCmd.AddCommand(DecodeBlocksCmd())
	return nil
}

func DecodeBlocksCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "decode-blocks",
		Short: "Decode base64 protobuf encoded blocks in JSON.",
		Long:  "Read base64 protobuf encoded blocks from stdin and write the corresponding JSON representation on stdout.",
		RunE: func(cmd *cobra.Command, args []string) error {
			scanner := bufio.NewScanner(cmd.InOrStdin())
			for scanner.Scan() {
				strB64 := scanner.Text()
				if strB64[0] == '"' {
					strB64 = strB64[1:]
				}
				if strB64[len(strB64)-1] == '"' {
					strB64 = strB64[:len(strB64)-1]
				}

				bytes, err := base64.StdEncoding.DecodeString(strB64)
				if err != nil {
					cmd.PrintErrln("could not decode base64 block:", err.Error())
					continue
				}

				block := new(types.Block)
				if err := block.Unmarshal(bytes); err != nil {
					cmd.PrintErrln("could not decode block:", err.Error())
					continue
				}

				json, err := json.Marshal(block)
				if err != nil {
					cmd.PrintErrln("could not marshal block in JSON:", err.Error())
					continue
				}
				cmd.Println(string(json))
			}

			return nil
		},
	}
}
