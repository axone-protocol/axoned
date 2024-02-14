package keeper_test

import (
	gocontext "context"
	"fmt"
	"io/fs"
	"testing"

	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

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

func TestGRPCAsk(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program        string
			query          string
			expectedAnswer *types.Answer
			expectedError  string
		}{
			{
				program: "foo.",
				query:   "foo.",
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   false,
					Variables: nil,
					Results:   []types.Result{{Substitutions: nil}},
				},
			},
			{
				program: "father(bob, alice).",
				query:   "father(bob, X).",
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   false,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable:   "X",
						Expression: "alice",
					}}}},
				},
			},
			{
				program: `father("bob", "alice").`,
				query:   `father("bob", X).`,
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   false,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable:   "X",
						Expression: "[a,l,i,c,e]",
					}}}},
				},
			},
			{
				program: "father(bob, alice). father(bob, john).",
				query:   "father(bob, X).",
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   true,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable:   "X",
						Expression: "alice",
					}}}},
				},
			},
			{
				program: "father(bob, alice).",
				query:   "father(bob, john).",
				expectedAnswer: &types.Answer{
					Success:   false,
					HasMore:   false,
					Variables: nil,
					Results:   nil,
				},
			},
			{
				program: "",
				query:   "block_height(X).",
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   false,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable:   "X",
						Expression: "0",
					}}}},
				},
			},
			{
				program: "father(bob, 'élodie').",
				query:   "father(bob, X).",
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   false,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable:   "X",
						Expression: "élodie",
					}}}},
				},
			},
			{
				program: "foo(foo(bar)).",
				query:   "foo(X).",
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   false,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable:   "X",
						Expression: "foo(bar)",
					}}}},
				},
			},
			{
				program: "father(bob, alice).",
				query:   "father(A, B).",
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   false,
					Variables: []string{"A", "B"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable:   "A",
						Expression: "bob",
					}, {
						Variable:   "B",
						Expression: "alice",
					}}}},
				},
			},
			{
				program: "father(bob, alice).",
				query:   "father(B, A).",
				expectedAnswer: &types.Answer{
					Success:   true,
					HasMore:   false,
					Variables: []string{"B", "A"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable:   "B",
						Expression: "bob",
					}, {
						Variable:   "A",
						Expression: "alice",
					}}}},
				},
			},
			{
				program: "father(bob, alice).",
				query:   "father(bob, X, O).",
				expectedAnswer: &types.Answer{
					Success:   false,
					Error:     "error(existence_error(procedure,father/3),root)",
					HasMore:   false,
					Variables: []string{"X", "O"},
					Results:   nil,
				},
			},
			{
				program:       "father°(bob, alice).",
				query:         "father(bob, X).",
				expectedError: "error compiling query: unexpected token: invalid(°): invalid argument",
			},
			{
				program:       "father(bob, alice).",
				query:         "father°(bob, X).",
				expectedError: "error executing query: unexpected token: invalid(°): invalid argument",
			},
			{
				program:       `father("bob", "alice").`,
				query:         `father("bob"", X).`,
				expectedError: "error executing query: EOF: invalid argument",
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d with program: %v and query: %v", nc, tc.program, tc.query),
				func() {
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
					err := logicKeeper.SetParams(testCtx.Ctx, types.DefaultParams())

					So(err, ShouldBeNil)

					Convey("and given a query with program and query to grpc", func() {
						queryHelper := baseapp.NewQueryServerTestHelper(testCtx.Ctx, encCfg.InterfaceRegistry)
						types.RegisterQueryServiceServer(queryHelper, logicKeeper)

						queryClient := types.NewQueryServiceClient(queryHelper)

						query := types.QueryServiceAskRequest{
							Program: tc.program,
							Query:   tc.query,
						}

						Convey("when the grpc query ask is called", func() {
							result, err := queryClient.Ask(gocontext.Background(), &query)

							Convey("Then it should return the expected answer", func() {
								if tc.expectedError != "" {
									So(err, ShouldNotBeNil)
									So(err.Error(), ShouldEqual, tc.expectedError)
									So(result, ShouldBeNil)
								} else {
									So(err, ShouldBeNil)
									So(result, ShouldNotBeNil)
									So(result.Answer, ShouldResemble, tc.expectedAnswer)
								}
							})
						})
					})
				})
		}
	})
}
