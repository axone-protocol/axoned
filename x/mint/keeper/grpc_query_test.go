package keeper_test

import (
	gocontext "context"
	"testing"

	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v15/x/mint"
	"github.com/axone-protocol/axoned/v15/x/mint/keeper"
	minttestutil "github.com/axone-protocol/axoned/v15/x/mint/testutil"
	"github.com/axone-protocol/axoned/v15/x/mint/types"
)

type grpcQueryTestContext struct {
	ctx         sdk.Context
	queryClient types.QueryClient
	mintKeeper  keeper.Keeper
}

func setupGRPCQueryTestContext(t *testing.T) *grpcQueryTestContext {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

	// gomock initializations
	ctrl := gomock.NewController(t)
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)

	accountKeeper.EXPECT().GetModuleAddress("mint").Return(sdk.AccAddress{})

	mintKeeper := keeper.NewKeeper(
		encCfg.Codec,
		runtime.NewKVStoreService(key),
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	if err := mintKeeper.Params.Set(testCtx.Ctx, types.DefaultParams()); err != nil {
		t.Fatalf("set params: %v", err)
	}
	if err := mintKeeper.Minter.Set(testCtx.Ctx, types.DefaultInitialMinter()); err != nil {
		t.Fatalf("set minter: %v", err)
	}

	queryHelper := baseapp.NewQueryServerTestHelper(testCtx.Ctx, encCfg.InterfaceRegistry)
	types.RegisterQueryServer(queryHelper, keeper.NewQueryServerImpl(mintKeeper))

	return &grpcQueryTestContext{
		ctx:         testCtx.Ctx,
		mintKeeper:  mintKeeper,
		queryClient: types.NewQueryClient(queryHelper),
	}
}

func TestGRPCParams(t *testing.T) {
	Convey("Given a mint query service", t, func() {
		testCtx := setupGRPCQueryTestContext(t)

		params, err := testCtx.queryClient.Params(gocontext.Background(), &types.QueryParamsRequest{})
		So(err, ShouldBeNil)
		keeperParams, err := testCtx.mintKeeper.Params.Get(testCtx.ctx)
		So(err, ShouldBeNil)
		So(params.Params, ShouldResemble, keeperParams)

		inflation, err := testCtx.queryClient.Inflation(gocontext.Background(), &types.QueryInflationRequest{})
		So(err, ShouldBeNil)
		minter, err := testCtx.mintKeeper.Minter.Get(testCtx.ctx)
		So(err, ShouldBeNil)
		So(inflation.Inflation, ShouldResemble, minter.Inflation)

		annualProvisions, err := testCtx.queryClient.AnnualProvisions(gocontext.Background(), &types.QueryAnnualProvisionsRequest{})
		So(err, ShouldBeNil)
		minter, err = testCtx.mintKeeper.Minter.Get(testCtx.ctx)
		So(err, ShouldBeNil)
		So(annualProvisions.AnnualProvisions, ShouldResemble, minter.AnnualProvisions)
	})
}
