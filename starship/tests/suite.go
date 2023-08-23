package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/cosmos/cosmos-sdk/types/module"

	starship "github.com/cosmology-tech/starship/clients/go/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	lens "github.com/strangelove-ventures/lens/client"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

var configFile = "../configs/devnet.yaml"

type TestSuite struct {
	suite.Suite

	config       *starship.Config
	chainClients starship.ChainClients
}

func (s *TestSuite) SetupTest() {
	s.T().Log("setting up e2e integration test suite...")

	// read config file from yaml
	yamlFile, err := os.ReadFile(configFile)
	s.Require().NoError(err)
	config := &starship.Config{}
	err = yaml.Unmarshal(yamlFile, config)
	s.Require().NoError(err)
	s.config = config

	// create chain clients
	modulesBasicMap := map[string][]module.AppModuleBasic{}
	for _, chain := range config.Chains {
		modulesBasicMap[chain.Name] = lens.ModuleBasics
	}
	chainClients, err := starship.NewChainClients(zap.L(), config, modulesBasicMap)
	s.Require().NoError(err)
	s.chainClients = chainClients
}

func (s *TestSuite) MakeRequest(req *http.Request, expCode int) io.Reader {
	resp, err := http.DefaultClient.Do(req)
	s.Require().NoError(err, "trying to make request", zap.Any("request", req))

	s.Require().Equal(expCode, resp.StatusCode, "response code did not match")

	return resp.Body
}

// WaitForHeight will wait till the chain reaches the block height
func (s *TestSuite) WaitForHeight(chain *starship.ChainClient, height int64) {
	s.Require().Eventuallyf(
		func() bool {
			curHeight, err := chain.GetHeight()
			s.Assert().NoError(err)
			if curHeight >= height {
				return true
			}
			return false
		},
		300*time.Second,
		5*time.Second,
		"waited for too long, still height did not reach desired block height",
	)
}

func (s *TestSuite) TransferTokens(chain *starship.ChainClient, addr string, amount int, denom string) {
	coin, err := sdk.ParseCoinNormalized(fmt.Sprintf("%d%s", amount, denom))
	s.Require().NoError(err)

	// Build transaction message
	req := &banktypes.MsgSend{
		FromAddress: chain.Address,
		ToAddress:   addr,
		Amount:      sdk.Coins{coin},
	}

	res, err := chain.Client.SendMsg(context.Background(), req, "Transfer tokens for e2e tests")
	s.Require().NoError(err)
	s.Require().NotEmpty(res)
}

func (s *TestSuite) IBCTransferTokens(chain1, chain2 *starship.ChainClient, chain2Addr string, amount int) {
	channel, err := chain1.GetIBCChannel(chain2.GetChainID())
	s.Require().NoError(err)

	denom, err := chain1.GetChainDenom()
	s.Require().NoError(err)

	coin := sdk.Coin{Denom: denom, Amount: sdk.NewInt(int64(amount))}
	req := &transfertypes.MsgTransfer{
		SourcePort:       channel.Chain_2.PortId,
		SourceChannel:    channel.Chain_2.ChannelId,
		Token:            coin,
		Sender:           chain1.Address,
		Receiver:         chain2Addr,
		TimeoutHeight:    clienttypes.NewHeight(12300, 45600),
		TimeoutTimestamp: 0,
		Memo:             fmt.Sprintf("testsetup: transfer token from %s to %s", chain1.GetChainID(), chain2.GetChainID()),
	}

	res, err := chain1.Client.SendMsg(context.Background(), req, "")
	s.Require().NoError(err)
	s.Require().NotEmpty(res)

	time.Sleep(5 * time.Second)
}
