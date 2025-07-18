package keeper_test

import (
	gocontext "context"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v12/x/mint"
	"github.com/axone-protocol/axoned/v12/x/mint/keeper"
	minttestutil "github.com/axone-protocol/axoned/v12/x/mint/testutil"
	"github.com/axone-protocol/axoned/v12/x/mint/types"
)

type MintTestSuite struct {
	suite.Suite

	ctx         sdk.Context
	queryClient types.QueryClient
	mintKeeper  keeper.Keeper
}

func (suite *MintTestSuite) SetupTest() {
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(suite.T(), key, storetypes.NewTransientStoreKey("transient_test"))
	suite.ctx = testCtx.Ctx

	// gomock initializations
	ctrl := gomock.NewController(suite.T())
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)

	accountKeeper.EXPECT().GetModuleAddress("mint").Return(sdk.AccAddress{})

	suite.mintKeeper = keeper.NewKeeper(
		encCfg.Codec,
		runtime.NewKVStoreService(key),
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	err := suite.mintKeeper.Params.Set(suite.ctx, types.DefaultParams())
	suite.Require().NoError(err)
	err = suite.mintKeeper.Minter.Set(suite.ctx, types.DefaultInitialMinter())
	suite.Require().NoError(err)

	queryHelper := baseapp.NewQueryServerTestHelper(testCtx.Ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(suite.mintKeeper))

	suite.queryClient = types.NewQueryClient(queryHelper)
}

func (suite *MintTestSuite) TestGRPCParams() {
	params, err := suite.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
	suite.Require().NoError(err)
	keeperParams, err := suite.mintKeeper.Params.Get(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(params.Params, keeperParams)

	inflation, err := suite.queryClient.Inflation(gocontext.Background(), &types.QueryInflationRequest{})
	suite.Require().NoError(err)
	minter, err := suite.mintKeeper.Minter.Get(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(inflation.Inflation, minter.Inflation)

	annualProvisions, err := suite.queryClient.AnnualProvisions(gocontext.Background(), &types.QueryAnnualProvisionsRequest{})
	suite.Require().NoError(err)
	minter, err = suite.mintKeeper.Minter.Get(suite.ctx)
	suite.Require().NoError(err)
	suite.Require().Equal(annualProvisions.AnnualProvisions, minter.AnnualProvisions)
}

func TestMintTestSuite(t *testing.T) {
	suite.Run(t, new(MintTestSuite))
}
