//nolint:gocognit,lll
package predicate

import (
	"fmt"
	"strings"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
)

func TestBech32(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program     string
			query       string
			wantResult  []types.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				query: `bech32_address(-(Hrp, Address), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantResult: []types.TermResults{{
					"Hrp":     "okp4",
					"Address": "[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]",
				}},
				wantSuccess: true,
			},
			{
				query: `bech32_address(Address, 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantResult: []types.TermResults{{
					"Address": "okp4-[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]",
				}},
				wantSuccess: true,
			},
			{
				query:       `bech32_address(-('okp4', X), foo(bar)).`,
				wantError:   fmt.Errorf("error(type_error(atom,foo(bar)),bech32_address/2)"),
				wantSuccess: false,
			},
			{
				query: `bech32_address(-('okp4', Address), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantResult: []types.TermResults{{
					"Address": "[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]",
				}},
				wantSuccess: true,
			},
			{
				query:       `bech32_address(-('okp5', Address), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantSuccess: false,
			},
			{
				query: `bech32_address(-('okp4', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).`,
				wantResult: []types.TermResults{{
					"Bech32": "okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn",
				}},
				wantSuccess: true,
			},
			{
				query:       `bech32_address(-('okp4', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				query:       `bech32_address(-(Hrp, [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantResult:  []types.TermResults{{"Hrp": "okp4"}},
				wantSuccess: true,
			},
			{
				query:       `bech32_address(-(Hrp, [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantResult:  []types.TermResults{{"Hrp": "okp4"}},
				wantSuccess: true,
			},
			{
				query:       `bech32_address(foo(bar), Bech32).`,
				wantError:   fmt.Errorf("error(type_error(pair,foo(bar)),bech32_address/2)"),
				wantSuccess: false,
			},
			{
				query:       `bech32_address(-('okp4', ['163',167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).`,
				wantError:   fmt.Errorf("error(type_error(byte,163),bech32_address/2)"),
				wantSuccess: false,
			},
			{
				query:       `bech32_address(-('okp4', [163,'x',23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).`,
				wantError:   fmt.Errorf("error(type_error(byte,x),bech32_address/2)"),
				wantSuccess: false,
			},
			{
				query:       `bech32_address(-(Hrp, [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]), Bech32).`,
				wantError:   fmt.Errorf("error(instantiation_error,bech32_address/2)"),
				wantSuccess: false,
			},
			{
				query:       `bech32_address(-('okp4', hey(2)), Bech32).`,
				wantError:   fmt.Errorf("error(type_error(list,hey(2)),bech32_address/2)"),
				wantSuccess: false,
			},
			{
				query: `bech32_address(-('okp4', X), foo).`,
				wantError: fmt.Errorf("error(domain_error(encoding(bech32),foo),[%s],bech32_address/2)",
					strings.Join(strings.Split("decoding bech32 failed: invalid bech32 string length 3", ""), ",")),
				wantSuccess: false,
			},
			{
				query:       `bech32_address(Address, Bech32).`,
				wantError:   fmt.Errorf("error(instantiation_error,bech32_address/2)"),
				wantSuccess: false,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := dbm.NewMemDB()
					stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewLightInterpreterMust(ctx)
						interpreter.Register2(engine.NewAtom("bech32_address"), Bech32Address)

						err := interpreter.Compile(ctx, tc.program)
						So(err, ShouldBeNil)

						Convey("When the predicate is called", func() {
							sols, err := interpreter.QueryContext(ctx, tc.query)

							Convey("Then the error should be nil", func() {
								So(err, ShouldBeNil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									var got []types.TermResults
									for sols.Next() {
										m := types.TermResults{}
										err := sols.Scan(m)
										So(err, ShouldBeNil)

										got = append(got, m)
									}
									if tc.wantError != nil {
										So(sols.Err(), ShouldNotBeNil)
										So(sols.Err().Error(), ShouldEqual, tc.wantError.Error())
									} else {
										So(sols.Err(), ShouldBeNil)

										if tc.wantSuccess {
											So(len(got), ShouldBeGreaterThan, 0)
											So(len(got), ShouldEqual, len(tc.wantResult))
											for iGot, resultGot := range got {
												for varGot, termGot := range resultGot {
													So(testutil.ReindexUnknownVariables(termGot), ShouldEqual, tc.wantResult[iGot][varGot])
												}
											}
										} else {
											So(len(got), ShouldEqual, 0)
										}
									}
								})
							})
						})
					})
				})
			})
		}
	})
}
