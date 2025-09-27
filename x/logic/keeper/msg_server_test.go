package keeper_test

import (
	gocontext "context"
	"fmt"
	"io/fs"
	"testing"

	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v13/x/logic"
	"github.com/axone-protocol/axoned/v13/x/logic/keeper"
	logictestutil "github.com/axone-protocol/axoned/v13/x/logic/testutil"
	"github.com/axone-protocol/axoned/v13/x/logic/types"
)

func TestUpdateParams(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			name      string
			request   *types.MsgUpdateParams
			expectErr bool
		}{
			{
				name: "set invalid authority",
				request: &types.MsgUpdateParams{
					Authority: "foo",
				},
				expectErr: true,
			},
			{
				name: "set full valid params",
				request: &types.MsgUpdateParams{
					Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
					Params:    types.DefaultParams(),
				},
				expectErr: false,
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d: %v, with request: %v", nc, tc.name, tc.request), func() {
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

					msgServer := keeper.NewMsgServerImpl(*logicKeeper)

					Convey("when call msg server to update params", func() {
						res, err := msgServer.UpdateParams(testCtx.Ctx, tc.request)

						Convey("then it should return the expected result", func() {
							if tc.expectErr {
								So(err, ShouldNotBeNil)
								So(res, ShouldBeNil)
							} else {
								So(err, ShouldBeNil)
								So(res, ShouldNotBeNil)
							}
						})
					})
				})
		}
	})
}
