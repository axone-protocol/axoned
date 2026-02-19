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

	"github.com/axone-protocol/axoned/v14/x/logic"
	"github.com/axone-protocol/axoned/v14/x/logic/keeper"
	logictestutil "github.com/axone-protocol/axoned/v14/x/logic/testutil"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

//nolint:lll,gocognit
func TestGRPCAsk(t *testing.T) {
	emptySolution := types.Result{}
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program               string
			query                 string
			limit                 uint64
			maxResultCount        uint64
			maxSize               uint64
			predicateBlacklist    []string
			virtualFilesWhitelist []string
			virtualFilesBlacklist []string
			maxGas                uint64
			maxVariables          uint64
			predicateCosts        map[string]uint64
			expectedAnswer        *types.Answer
			expectedError         string
		}{
			{
				program: "foo.",
				query:   "foo.",
				expectedAnswer: &types.Answer{
					Results: []types.Result{emptySolution},
				},
			},
			{
				program:        "father(bob, alice).",
				query:          "father(bob, john).",
				expectedAnswer: &types.Answer{},
			},
			{
				program: "father(bob, alice).",
				query:   "father(bob, X).",
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "alice",
					}}}},
				},
			},
			{
				program: `father("bob", "alice").`,
				query:   `father("bob", X).`,
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "[a,l,i,c,e]",
					}}}},
				},
			},
			{
				program: "father(bob, alice). father(bob, john).",
				query:   "father(bob, X).",
				expectedAnswer: &types.Answer{
					HasMore:   true,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "alice",
					}}}},
				},
			},
			{
				program:        "father(bob, alice). father(bob, john).",
				query:          "father(bob, X).",
				maxResultCount: 5,
				expectedAnswer: &types.Answer{
					HasMore:   true,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "alice",
					}}}},
				},
			},
			{
				program:        "father(bob, alice). father(bob, john).",
				query:          "father(bob, X).",
				limit:          2,
				maxResultCount: 5,
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "alice",
					}}}, {Substitutions: []types.Substitution{{
						Variable: "X", Expression: "john",
					}}}},
				},
			},
			{
				program:        "father(bob, alice). father(bob, john).",
				query:          "father(bob, X).",
				limit:          2,
				maxResultCount: 1,
				expectedAnswer: &types.Answer{
					HasMore:   true,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "alice",
					}}}},
				},
			},
			{
				program:       "father(bob, alice). father(bob, john).",
				query:         "father(bob, X).",
				maxSize:       5,
				expectedError: "query: 15 > MaxSize: 5: limit exceeded",
			},
			{
				program: "father(bob, alice). father(bob, john).",
				query:   "father(bob, X).",
				maxSize: 0,
				expectedAnswer: &types.Answer{
					HasMore:   true,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "alice",
					}}}},
				},
			},
			{
				program: "block_height(X) :- block_header(Header), X = Header.height.",
				query:   "block_height(X).",
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "0",
					}}}},
				},
			},
			{
				program:       "block_height(X) :- block_header(Header), X = Header.height.",
				query:         "block_height(X).",
				maxGas:        1000,
				expectedError: "out of gas: logic <ReadPerByte> (1018/1000): limit exceeded",
			},
			{
				program: "block_height(X) :- block_header(Header), X = Header.height.",
				query:   "block_height(X).",
				maxGas:  3000,
				predicateCosts: map[string]uint64{
					"block_header/1": 10000,
				},
				expectedError: "out of gas: logic <block_header/1> (11141/3000): limit exceeded",
			},
			{
				program:       "recursionOfDeath :- recursionOfDeath.",
				query:         "recursionOfDeath.",
				maxGas:        3000,
				expectedError: "out of gas: logic <recursionOfDeath/0> (3001/3000): limit exceeded",
			},
			{
				program:       "backtrackOfDeath :- repeat, fail.",
				query:         "backtrackOfDeath.",
				maxGas:        3014,
				expectedError: "out of gas: logic <fail/0> (3015/3014): limit exceeded",
			},
			{
				query:         "length(List, 100000).",
				maxVariables:  1000,
				expectedError: "maximum number of variables reached: limit exceeded",
			},
			{
				program:       "l(L) :- length(L, 1). l(L) :- length(L, 1000).",
				query:         "l(L).",
				limit:         2,
				maxVariables:  1000,
				expectedError: "maximum number of variables reached: limit exceeded",
			},
			{
				program: "father(bob, 'élodie').",
				query:   "father(bob, X).",
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "élodie",
					}}}},
				},
			},
			{
				program: "foo(foo(bar)).",
				query:   "foo(X).",
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "foo(bar)",
					}}}},
				},
			},
			{
				program: "father(bob, alice).",
				query:   "father(A, B).",
				expectedAnswer: &types.Answer{
					Variables: []string{"A", "B"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "A", Expression: "bob",
					}, {
						Variable: "B", Expression: "alice",
					}}}},
				},
			},
			{
				program: "father(bob, alice).",
				query:   "father(B, A).",
				expectedAnswer: &types.Answer{
					Variables: []string{"B", "A"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "B", Expression: "bob",
					}, {
						Variable: "A", Expression: "alice",
					}}}},
				},
			},
			{
				program: "father(bob, X) :- true.",
				query:   "father(B, X).",
				expectedAnswer: &types.Answer{
					Variables: []string{"B", "X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "B", Expression: "bob",
					}}}},
				},
			},
			{
				program: "father(bob, alice).",
				query:   "father(bob, X, O).",
				expectedAnswer: &types.Answer{
					Variables: []string{"X", "O"},
					Results:   []types.Result{{Error: "error(existence_error(procedure,father/3),root)"}},
				},
			},
			{
				program:            "",
				query:              "block_header(X).",
				predicateBlacklist: []string{"block_header/1"},
				expectedAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{"X"},
					Results:   []types.Result{{Error: "error(permission_error(execute,forbidden_predicate,block_header/1),root)"}},
				},
			},
			{
				program:            "contains_forbidden_predicate(X) :- block_header(X).",
				query:              "contains_forbidden_predicate(X).",
				predicateBlacklist: []string{"block_header/1"},
				expectedAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{"X"},
					Results:   []types.Result{{Error: "error(permission_error(execute,forbidden_predicate,block_header/1),contains_forbidden_predicate/1)"}},
				},
			},
			{
				program:            "cannot_be_blacklisted(X) :- X = 42.",
				query:              "cannot_be_blacklisted(X).",
				predicateBlacklist: []string{"cannot_be_blacklisted/1"},
				expectedAnswer: &types.Answer{
					HasMore:   false,
					Variables: []string{"X"},
					Results: []types.Result{{Substitutions: []types.Substitution{{
						Variable: "X", Expression: "42",
					}}}},
				},
			},
			{
				program:               `test :- open('cosmwasm:okp4-objectarium:okp41ffzp0xmjhwkltuxcvccl0z9tyfuu7txp5ke0tpkcjpzuq9fcj3pqrteqt3?query=%7B%22object_data%22%3A%7B%22id%22%3A%22content1%22%7D%7D', read, _, []).`,
				query:                 "test.",
				virtualFilesBlacklist: []string{"cosmwasm:"},
				expectedAnswer: &types.Answer{
					HasMore: false,
					Results: []types.Result{{Error: "error(permission_error(open,source_sink,cosmwasm:okp4-objectarium:okp41ffzp0xmjhwkltuxcvccl0z9tyfuu7txp5ke0tpkcjpzuq9fcj3pqrteqt3?query=%7B%22object_data%22%3A%7B%22id%22%3A%22content1%22%7D%7D),open/4)"}},
				},
			},
			{
				program:               `test :- open('https://example.com/data.pl', read, _, []).`,
				query:                 "test.",
				virtualFilesWhitelist: []string{"cosmwasm:"},
				expectedAnswer: &types.Answer{
					HasMore: false,
					Results: []types.Result{{Error: "error(permission_error(open,source_sink,https://example.com/data.pl),open/4)"}},
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
			{
				program: `
				foo(a1).
				foo(a2).
				foo(a3) :- throw(error(resource_error(foo))).
				foo(a4).
				`,
				query:          `foo(X).`,
				maxResultCount: 1,
				expectedAnswer: &types.Answer{
					HasMore:   true,
					Variables: []string{"X"},
					Results: []types.Result{
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a1"}}},
					},
				},
			},
			{
				program: `
				foo(a1).
				foo(a2).
				foo(a3) :- throw(error(resource_error(foo))).
				foo(a4).
				`,
				query:          `foo(X).`,
				limit:          2,
				maxResultCount: 3,
				expectedAnswer: &types.Answer{
					HasMore:   true,
					Variables: []string{"X"},
					Results: []types.Result{
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a1"}}},
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a2"}}},
					},
				},
			},
			{
				program: `
				foo(a1).
				foo(a2).
				foo(a3) :- throw(error(resource_error(foo))).
				foo(a4).
				`,
				query:          `foo(X).`,
				limit:          3,
				maxResultCount: 5,
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a1"}}},
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a2"}}},
						{Error: "error(resource_error(foo))"},
					},
				},
			},
			{
				program: `
				foo(a1).
				foo(a2).
				foo(a3) :- throw(error(resource_error(foo))).
				foo(a4).
				`,
				query:          `foo(X).`,
				limit:          5,
				maxResultCount: 5,
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a1"}}},
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a2"}}},
						{Error: "error(resource_error(foo))"},
					},
				},
			},
			{
				program: `
				foo(a1).
				foo(a2).
				foo(a3) :- throw(error(resource_error(foo))).
				foo(a4).
				`,
				query:          `foo(X).`,
				limit:          5,
				maxResultCount: 0,
				expectedAnswer: &types.Answer{
					Variables: []string{"X"},
					Results: []types.Result{
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a1"}}},
						{Substitutions: []types.Substitution{{Variable: "X", Expression: "a2"}}},
						{Error: "error(resource_error(foo))"},
					},
				},
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
						func(_ gocontext.Context) (fs.FS, error) {
							return fsProvider, nil
						})

					params := types.DefaultParams()
					params.Limits.MaxResultCount = tc.maxResultCount
					params.Limits.MaxSize = tc.maxSize
					params.Limits.MaxVariables = tc.maxVariables

					if tc.predicateBlacklist != nil {
						params.Interpreter.PredicatesFilter.Blacklist = tc.predicateBlacklist
					}
					if tc.virtualFilesWhitelist != nil {
						params.Interpreter.VirtualFilesFilter.Whitelist = tc.virtualFilesWhitelist
					}
					if tc.virtualFilesBlacklist != nil {
						params.Interpreter.VirtualFilesFilter.Blacklist = tc.virtualFilesBlacklist
					}
					if tc.predicateCosts != nil {
						predicateCosts := make([]types.PredicateCost, 0, len(tc.predicateCosts))
						for predicate, cost := range tc.predicateCosts {
							predicateCosts = append(predicateCosts, types.PredicateCost{
								Predicate: predicate,
								Cost:      cost,
							})
						}
						params.GasPolicy.PredicateCosts = predicateCosts
					}
					err := logicKeeper.SetParams(testCtx.Ctx, params)

					So(err, ShouldBeNil)

					if tc.maxGas != 0 {
						testCtx.Ctx = testCtx.Ctx.WithGasMeter(storetypes.NewGasMeter(tc.maxGas))
					} else {
						testCtx.Ctx = testCtx.Ctx.WithGasMeter(storetypes.NewInfiniteGasMeter())
					}

					Convey("and given a query with program and query to grpc", func() {
						queryHelper := baseapp.NewQueryServerTestHelper(testCtx.Ctx, encCfg.InterfaceRegistry)
						types.RegisterQueryServiceServer(queryHelper, logicKeeper)

						queryClient := types.NewQueryServiceClient(queryHelper)

						query := types.QueryServiceAskRequest{
							Program: tc.program,
							Query:   tc.query,
							Limit:   tc.limit,
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

		Convey("When the query ask is called with a nil query", func() {
			response, err := logicKeeper.Ask(testCtx.Ctx, nil)

			Convey("Then it should return an error", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "request is nil: invalid argument")
				So(response, ShouldBeNil)
			})
		})
	})
}
