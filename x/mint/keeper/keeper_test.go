package keeper_test

import (
	"reflect"
	"testing"

	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v14/x/mint"
	"github.com/axone-protocol/axoned/v14/x/mint/keeper"
	minttestutil "github.com/axone-protocol/axoned/v14/x/mint/testutil"
	"github.com/axone-protocol/axoned/v14/x/mint/types"
)

type integrationTestContext struct {
	mintKeeper    keeper.Keeper
	ctx           sdk.Context
	msgServer     types.MsgServer
	stakingKeeper *minttestutil.MockStakingKeeper
	bankKeeper    *minttestutil.MockBankKeeper
}

func setupIntegrationTestContext(t *testing.T) *integrationTestContext {
	t.Helper()

	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})
	key := storetypes.NewKVStoreKey(types.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

	// gomock initializations
	ctrl := gomock.NewController(t)
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)

	accountKeeper.EXPECT().GetModuleAddress(types.ModuleName).Return(sdk.AccAddress{})

	mintKeeper := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		stakingKeeper,
		accountKeeper,
		bankKeeper,
		authtypes.FeeCollectorName,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)

	if !reflect.DeepEqual(testCtx.Ctx.Logger().With("module", "x/"+types.ModuleName), mintKeeper.Logger(testCtx.Ctx)) {
		t.Fatal("unexpected module logger")
	}

	if err := mintKeeper.Params.Set(testCtx.Ctx, types.DefaultParams()); err != nil {
		t.Fatalf("set params: %v", err)
	}
	if err := mintKeeper.Minter.Set(testCtx.Ctx, types.DefaultInitialMinter()); err != nil {
		t.Fatalf("set minter: %v", err)
	}

	return &integrationTestContext{
		mintKeeper:    mintKeeper,
		ctx:           testCtx.Ctx,
		msgServer:     keeper.NewMsgServerImpl(mintKeeper),
		stakingKeeper: stakingKeeper,
		bankKeeper:    bankKeeper,
	}
}

func TestAliasFunctions(t *testing.T) {
	Convey("Given a mint keeper", t, func() {
		tc := setupIntegrationTestContext(t)

		stakingTokenSupply := math.NewIntFromUint64(100000000000)
		tc.stakingKeeper.EXPECT().StakingTokenSupply(tc.ctx).Return(stakingTokenSupply, nil)
		tokenSupply, err := tc.mintKeeper.StakingTokenSupply(tc.ctx)
		So(err, ShouldBeNil)
		So(tokenSupply, ShouldResemble, stakingTokenSupply)

		bondedRatio := math.LegacyNewDecWithPrec(15, 2)
		tc.stakingKeeper.EXPECT().BondedRatio(tc.ctx).Return(bondedRatio, nil)
		ratio, err := tc.mintKeeper.BondedRatio(tc.ctx)
		So(err, ShouldBeNil)
		So(ratio, ShouldResemble, bondedRatio)

		coins := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1000000)))
		tc.bankKeeper.EXPECT().MintCoins(tc.ctx, types.ModuleName, coins).Return(nil)
		So(tc.mintKeeper.MintCoins(tc.ctx, sdk.NewCoins()), ShouldBeNil)
		So(tc.mintKeeper.MintCoins(tc.ctx, coins), ShouldBeNil)

		fees := sdk.NewCoins(sdk.NewCoin("stake", math.NewInt(1000)))
		tc.bankKeeper.EXPECT().SendCoinsFromModuleToModule(tc.ctx, types.ModuleName, authtypes.FeeCollectorName, fees).Return(nil)
		So(tc.mintKeeper.AddCollectedFees(tc.ctx, fees), ShouldBeNil)
	})
}
