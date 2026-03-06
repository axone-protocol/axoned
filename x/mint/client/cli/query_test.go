package cli_test

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	rpcclientmock "github.com/cometbft/cometbft/rpc/client/mock"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	testutilmod "github.com/cosmos/cosmos-sdk/types/module/testutil"

	"github.com/axone-protocol/axoned/v14/x/mint"
	mintcli "github.com/axone-protocol/axoned/v14/x/mint/client/cli"
)

func TestGetCmdQueryParams(t *testing.T) {
	Convey("Given the params query command", t, func() {
		encCfg := testutilmod.MakeTestEncodingConfig(mint.AppModuleBasic{})
		kr := keyring.NewInMemory(encCfg.Codec)
		baseCtx := client.Context{}.
			WithKeyring(kr).
			WithTxConfig(encCfg.TxConfig).
			WithCodec(encCfg.Codec).
			WithClient(clitestutil.MockCometRPC{Client: rpcclientmock.Client{}}).
			WithAccountRetriever(client.MockAccountRetriever{}).
			WithOutput(io.Discard).
			WithChainID("test-chain")

		cmd := mintcli.GetCmdQueryParams()

		testCases := []struct {
			name           string
			flagArgs       []string
			expCmdOutput   string
			expectedOutput string
		}{
			{
				"json output",
				[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", flags.FlagOutput)},
				`[--height=1 --output=json]`,
				`{"mint_denom":"","inflation_coef":"0","blocks_per_year":"0","inflation_max":null,"inflation_min":null}`,
			},
			{
				"text output",
				[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=text", flags.FlagOutput)},
				`[--height=1 --output=text]`,
				`blocks_per_year: "0"
inflation_coef: "0"
inflation_max: null
inflation_min: null
mint_denom: ""`,
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				ctx := svrcmd.CreateExecuteContext(context.Background())

				cmd.SetOut(io.Discard)
				So(cmd, ShouldNotBeNil)

				cmd.SetContext(ctx)
				cmd.SetArgs(tc.flagArgs)

				So(client.SetCmdClientContextHandler(baseCtx, cmd), ShouldBeNil)

				if len(tc.flagArgs) != 0 {
					So(fmt.Sprint(cmd), ShouldContainSubstring, "params [] [] Query the current minting parameters")
					So(fmt.Sprint(cmd), ShouldContainSubstring, tc.expCmdOutput)
				}

				out, err := clitestutil.ExecTestCLICmd(baseCtx, cmd, tc.flagArgs)
				So(err, ShouldBeNil)
				So(strings.TrimSpace(out.String()), ShouldEqual, tc.expectedOutput)
			})
		}
	})
}

func TestGetCmdQueryInflation(t *testing.T) {
	Convey("Given the inflation query command", t, func() {
		encCfg := testutilmod.MakeTestEncodingConfig(mint.AppModuleBasic{})
		kr := keyring.NewInMemory(encCfg.Codec)
		baseCtx := client.Context{}.
			WithKeyring(kr).
			WithTxConfig(encCfg.TxConfig).
			WithCodec(encCfg.Codec).
			WithClient(clitestutil.MockCometRPC{Client: rpcclientmock.Client{}}).
			WithAccountRetriever(client.MockAccountRetriever{}).
			WithOutput(io.Discard).
			WithChainID("test-chain")

		cmd := mintcli.GetCmdQueryInflation()

		testCases := []struct {
			name           string
			flagArgs       []string
			expCmdOutput   string
			expectedOutput string
		}{
			{
				"json output",
				[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", flags.FlagOutput)},
				`[--height=1 --output=json]`,
				`<nil>`,
			},
			{
				"text output",
				[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=text", flags.FlagOutput)},
				`[--height=1 --output=text]`,
				`<nil>`,
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				ctx := svrcmd.CreateExecuteContext(context.Background())

				cmd.SetOut(io.Discard)
				So(cmd, ShouldNotBeNil)

				cmd.SetContext(ctx)
				cmd.SetArgs(tc.flagArgs)

				So(client.SetCmdClientContextHandler(baseCtx, cmd), ShouldBeNil)

				if len(tc.flagArgs) != 0 {
					So(fmt.Sprint(cmd), ShouldContainSubstring, "inflation [] [] Query the current minting inflation value")
					So(fmt.Sprint(cmd), ShouldContainSubstring, tc.expCmdOutput)
				}

				out, err := clitestutil.ExecTestCLICmd(baseCtx, cmd, tc.flagArgs)
				So(err, ShouldBeNil)
				So(strings.TrimSpace(out.String()), ShouldEqual, tc.expectedOutput)
			})
		}
	})
}

func TestGetCmdQueryAnnualProvisions(t *testing.T) {
	Convey("Given the annual provisions query command", t, func() {
		encCfg := testutilmod.MakeTestEncodingConfig(mint.AppModuleBasic{})
		kr := keyring.NewInMemory(encCfg.Codec)
		baseCtx := client.Context{}.
			WithKeyring(kr).
			WithTxConfig(encCfg.TxConfig).
			WithCodec(encCfg.Codec).
			WithClient(clitestutil.MockCometRPC{Client: rpcclientmock.Client{}}).
			WithAccountRetriever(client.MockAccountRetriever{}).
			WithOutput(io.Discard).
			WithChainID("test-chain")

		cmd := mintcli.GetCmdQueryAnnualProvisions()

		testCases := []struct {
			name           string
			flagArgs       []string
			expCmdOutput   string
			expectedOutput string
		}{
			{
				"json output",
				[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=json", flags.FlagOutput)},
				`[--height=1 --output=json]`,
				`<nil>`,
			},
			{
				"text output",
				[]string{fmt.Sprintf("--%s=1", flags.FlagHeight), fmt.Sprintf("--%s=text", flags.FlagOutput)},
				`[--height=1 --output=text]`,
				`<nil>`,
			},
		}

		for _, tc := range testCases {
			Convey(tc.name, func() {
				ctx := svrcmd.CreateExecuteContext(context.Background())

				cmd.SetOut(io.Discard)
				So(cmd, ShouldNotBeNil)

				cmd.SetContext(ctx)
				cmd.SetArgs(tc.flagArgs)

				So(client.SetCmdClientContextHandler(baseCtx, cmd), ShouldBeNil)

				if len(tc.flagArgs) != 0 {
					So(fmt.Sprint(cmd), ShouldContainSubstring, "annual-provisions [] [] Query the current minting annual provisions value")
					So(fmt.Sprint(cmd), ShouldContainSubstring, tc.expCmdOutput)
				}

				out, err := clitestutil.ExecTestCLICmd(baseCtx, cmd, tc.flagArgs)
				So(err, ShouldBeNil)
				So(strings.TrimSpace(out.String()), ShouldEqual, tc.expectedOutput)
			})
		}
	})
}
