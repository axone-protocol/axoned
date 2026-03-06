package keeper_test

import (
	"errors"
	"testing"

	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/collections"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/axone-protocol/axoned/v14/x/mint"
	"github.com/axone-protocol/axoned/v14/x/mint/keeper"
	minttestutil "github.com/axone-protocol/axoned/v14/x/mint/testutil"
	"github.com/axone-protocol/axoned/v14/x/mint/types"
)

var minterAcc = authtypes.NewEmptyModuleAccount(types.ModuleName, authtypes.Minter)

type genesisTestContext struct {
	sdkCtx        sdk.Context
	keeper        keeper.Keeper
	cdc           codec.BinaryCodec
	accountKeeper types.AccountKeeper
	key           *storetypes.KVStoreKey
}

func setupGenesisTestContext(t *testing.T) *genesisTestContext {
	t.Helper()

	key := storetypes.NewKVStoreKey(types.StoreKey)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(mint.AppModuleBasic{})

	// gomock initializations
	ctrl := gomock.NewController(t)
	cdc := codec.NewProtoCodec(encCfg.InterfaceRegistry)

	stakingKeeper := minttestutil.NewMockStakingKeeper(ctrl)
	accountKeeper := minttestutil.NewMockAccountKeeper(ctrl)
	bankKeeper := minttestutil.NewMockBankKeeper(ctrl)

	accountKeeper.EXPECT().GetModuleAddress(minterAcc.Name).Return(minterAcc.GetAddress())
	accountKeeper.EXPECT().GetModuleAccount(testCtx.Ctx, minterAcc.Name).Return(minterAcc)

	return &genesisTestContext{
		sdkCtx: testCtx.Ctx,
		keeper: keeper.NewKeeper(
			cdc,
			runtime.NewKVStoreService(key),
			stakingKeeper,
			accountKeeper,
			bankKeeper,
			"",
			"cosmos10d07y265gmmuvt4z0w9aw880jnsr700j6zn9kn",
		),
		cdc:           cdc,
		accountKeeper: accountKeeper,
		key:           key,
	}
}

func TestImportExportGenesis(t *testing.T) {
	Convey("Given a mint keeper genesis state", t, func() {
		testCtx := setupGenesisTestContext(t)

		genesisState := types.DefaultGenesisState()
		genesisState.Minter = types.NewMinter(math.LegacyNewDecWithPrec(20, 2), math.LegacyNewDec(1))
		genesisState.Params = types.NewParams(
			"testDenom",
			math.LegacyNewDecWithPrec(69, 2),
			uint64(60*60*8766/5),
		)

		testCtx.keeper.InitGenesis(testCtx.sdkCtx, testCtx.accountKeeper, genesisState)

		minter, err := testCtx.keeper.Minter.Get(testCtx.sdkCtx)
		So(err, ShouldBeNil)
		So(minter, ShouldResemble, genesisState.Minter)

		invalidCtx := testutil.DefaultContextWithDB(t, testCtx.key, storetypes.NewTransientStoreKey("transient_test"))
		_, err = testCtx.keeper.Minter.Get(invalidCtx.Ctx)
		So(errors.Is(err, collections.ErrNotFound), ShouldBeTrue)

		params, err := testCtx.keeper.Params.Get(testCtx.sdkCtx)
		So(err, ShouldBeNil)
		So(params, ShouldResemble, genesisState.Params)

		genesisState2 := testCtx.keeper.ExportGenesis(testCtx.sdkCtx)
		So(genesisState2, ShouldResemble, genesisState)
	})
}
