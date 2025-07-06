package keeper_test

import (
	gocontext "context"
	"fmt"
	"io/fs"
	"testing"

	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v12/x/logic"
	"github.com/axone-protocol/axoned/v12/x/logic/keeper"
	logictestutil "github.com/axone-protocol/axoned/v12/x/logic/testutil"
	"github.com/axone-protocol/axoned/v12/x/logic/types"
)

func TestGRPCParams(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			params types.Params
		}{
			{
				params: types.NewParams(
					types.NewInterpreter(
						types.WithBootstrap("bootstrap"),
						types.WithPredicatesBlacklist([]string{"halt/1"}),
						types.WithPredicatesWhitelist([]string{"source_file/1"}),
						types.WithVirtualFilesBlacklist([]string{"file1"}),
						types.WithVirtualFilesWhitelist([]string{"file2"}),
					),
					types.NewLimits(
						types.WithMaxSize(2),
						types.WithMaxResultCount(3),
						types.WithMaxUserOutputSize(4),
						types.WithMaxVariables(5),
					),
					types.GasPolicy{},
				),
			},
			{
				params: types.NewParams(
					types.NewInterpreter(
						types.WithBootstrap("bootstrap"),
						types.WithPredicatesBlacklist([]string{"halt/1"}),
						types.WithPredicatesWhitelist([]string{"source_file/1"}),
						types.WithVirtualFilesBlacklist([]string{"file1"}),
						types.WithVirtualFilesWhitelist([]string{"file2"}),
					),
					types.NewLimits(
						types.WithMaxSize(2),
						types.WithMaxResultCount(3),
						types.WithMaxUserOutputSize(4),
						types.WithMaxVariables(5),
					),
					types.GasPolicy{
						WeightingFactor:      2,
						DefaultPredicateCost: 1,
						PredicateCosts: []types.PredicateCost{
							{Predicate: "foo", Cost: 1},
						},
					},
				),
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d with params: %v", nc, tc.params), func() {
					encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
					key := storetypes.NewKVStoreKey(types.StoreKey)
					testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

					// gomock initializations
					ctrl := gomock.NewController(t)
					accountKeeper := logictestutil.NewMockAccountKeeper(ctrl)
					authQueryService := logictestutil.NewMockAuthQueryService(ctrl)
					bankKeeper := logictestutil.NewMockBankKeeper(ctrl)
					fsProvider := logictestutil.NewMockFS(ctrl)

					logicKeeper := keeper.NewKeeper(
						encCfg.Codec,
						encCfg.InterfaceRegistry,
						key,
						key,
						authtypes.NewModuleAddress(govtypes.ModuleName),
						accountKeeper,
						authQueryService,
						bankKeeper,
						func(_ gocontext.Context) fs.FS {
							return fsProvider
						})

					Convey("and given params to the keeper", func() {
						err := logicKeeper.SetParams(testCtx.Ctx, tc.params)
						So(err, ShouldBeNil)

						queryHelper := baseapp.NewQueryServerTestHelper(testCtx.Ctx, encCfg.InterfaceRegistry)
						types.RegisterQueryServiceServer(queryHelper, logicKeeper)

						queryClient := types.NewQueryServiceClient(queryHelper)

						Convey("when the grpc query params is called", func() {
							params, err := queryClient.Params(gocontext.Background(), &types.QueryServiceParamsRequest{})

							Convey("Then it should return the expected params set to the keeper", func() {
								So(err, ShouldBeNil)
								So(params.Params, ShouldResemble, tc.params)
							})
						})
					})
				})
		}
	})

	Convey("Given a keeper", t, func() {
		encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
		key := storetypes.NewKVStoreKey(types.StoreKey)
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

		logicKeeper := keeper.NewKeeper(
			encCfg.Codec,
			encCfg.InterfaceRegistry,
			key,
			key,
			authtypes.NewModuleAddress(govtypes.ModuleName),
			nil,
			nil,
			nil,
			nil)

		Convey("When the query params is called with a nil query", func() {
			params, err := logicKeeper.Params(testCtx.Ctx, nil)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "rpc error: code = InvalidArgument desc = invalid request")
				So(params, ShouldBeNil)
			})
		})
	})
}
