package keeper_test

import (
	gocontext "context"
	"fmt"
	"io/fs"
	"testing"

	"github.com/samber/lo"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v12/x/logic"
	"github.com/axone-protocol/axoned/v12/x/logic/keeper"
	v1beta2types "github.com/axone-protocol/axoned/v12/x/logic/legacy/v1beta2/types"
	logictestutil "github.com/axone-protocol/axoned/v12/x/logic/testutil"
	"github.com/axone-protocol/axoned/v12/x/logic/types"
)

func TestMigrateStoreV10ToV11(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			params v1beta2types.Params
			expect types.Params
		}{
			{
				params: v1beta2types.Params{
					Interpreter: v1beta2types.Interpreter{
						PredicatesFilter:   v1beta2types.Filter{},
						Bootstrap:          "",
						VirtualFilesFilter: v1beta2types.Filter{},
					},
					GasPolicy: v1beta2types.GasPolicy{},
					Limits:    v1beta2types.Limits{},
				},
				expect: types.Params{
					Interpreter: types.Interpreter{
						PredicatesFilter:   types.Filter{},
						Bootstrap:          "",
						VirtualFilesFilter: types.Filter{},
					},
					GasPolicy: types.GasPolicy{},
					Limits:    types.Limits{},
				},
			},
			{
				params: v1beta2types.Params{
					Interpreter: v1beta2types.Interpreter{
						PredicatesFilter: v1beta2types.Filter{
							Whitelist: []string{"foo/1", "bar/2"},
							Blacklist: []string{"baz/3"},
						},
						Bootstrap: "foo(bar).",
						VirtualFilesFilter: v1beta2types.Filter{
							Whitelist: []string{"foo://bar"},
							Blacklist: []string{"bar://baz"},
						},
					},
					GasPolicy: v1beta2types.GasPolicy{
						WeightingFactor:      lo.ToPtr(sdkmath.NewUint(42)),
						DefaultPredicateCost: lo.ToPtr(sdkmath.NewUint(66)),
						PredicateCosts: []v1beta2types.PredicateCost{
							{
								Predicate: "foo/1",
								Cost:      lo.ToPtr(sdkmath.NewUint(99)),
							},
						},
					},
					Limits: v1beta2types.Limits{
						MaxSize:           lo.ToPtr(sdkmath.NewUint(100)),
						MaxResultCount:    lo.ToPtr(sdkmath.NewUint(10)),
						MaxUserOutputSize: lo.ToPtr(sdkmath.NewUint(50)),
						MaxVariables:      lo.ToPtr(sdkmath.NewUint(5)),
					},
				},
				expect: types.Params{
					Interpreter: types.Interpreter{
						PredicatesFilter: types.Filter{
							Whitelist: []string{"foo/1", "bar/2"},
							Blacklist: []string{"baz/3"},
						},
						Bootstrap: "foo(bar).",
						VirtualFilesFilter: types.Filter{
							Whitelist: []string{"foo://bar"},
							Blacklist: []string{"bar://baz"},
						},
					},
					GasPolicy: types.GasPolicy{
						WeightingFactor:      42,
						DefaultPredicateCost: 66,
						PredicateCosts: []types.PredicateCost{
							{
								Predicate: "foo/1",
								Cost:      99,
							},
						},
					},
					Limits: types.Limits{
						MaxSize:           100,
						MaxResultCount:    10,
						MaxUserOutputSize: 50,
						MaxVariables:      5,
					},
				},
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given a mocked logic keeper for test case %d", nc), func() {
				encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
				key := storetypes.NewKVStoreKey(types.StoreKey)
				testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

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
				So(logicKeeper, ShouldNotBeNil)

				Convey("Given a store with v10 params", func() {
					store := testCtx.Ctx.KVStore(key)
					bz, err := encCfg.Codec.Marshal(&tc.params)
					So(err, ShouldBeNil)

					store.Set(types.ParamsKey, bz)

					Convey("When migrating store from v10 to v11", func() {
						migrateHandler := keeper.MigrateStoreV3ToV4(*logicKeeper)
						So(migrateHandler, ShouldNotBeNil)

						err := migrateHandler(testCtx.Ctx)
						So(err, ShouldBeNil)

						Convey("Then the store should have the expected v11 params", func() {
							params := logicKeeper.GetParams(testCtx.Ctx)
							So(err, ShouldBeNil)
							So(params, ShouldResemble, tc.expect)
						})
					})
				})
			})
		}
	})
}
