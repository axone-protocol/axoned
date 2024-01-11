package keeper_test

import (
	gocontext "context"
	"fmt"
	"io/fs"
	"testing"

	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/math"

	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/okp4/okp4d/x/logic"
	"github.com/okp4/okp4d/x/logic/keeper"
	logictestutil "github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
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
						types.WithMaxGas(math.NewUint(1)),
						types.WithMaxSize(math.NewUint(2)),
						types.WithMaxResultCount(math.NewUint(3)),
						types.WithMaxUserOutputSize(math.NewUint(4)),
					),
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
					bankKeeper := logictestutil.NewMockBankKeeper(ctrl)
					fsProvider := logictestutil.NewMockFS(ctrl)

					logicKeeper := keeper.NewKeeper(
						encCfg.Codec,
						key,
						key,
						authtypes.NewModuleAddress(govtypes.ModuleName),
						accountKeeper,
						bankKeeper,
						func(ctx gocontext.Context) fs.FS {
							return fsProvider
						},
					)

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
}
