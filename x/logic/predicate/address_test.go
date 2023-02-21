package predicate

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
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
				query: `bech32_address(Hrp, Address, 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantResult: []types.TermResults{{
					"Hrp":     "okp4",
					"Address": "[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]",
				}},
				wantSuccess: true,
			},
			{
				query: `bech32_address('okp4', Address, 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantResult: []types.TermResults{{
					"Address": "[163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30]",
				}},
				wantSuccess: true,
			},
			{
				query:       `bech32_address('okp5', Address, 'okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn').`,
				wantSuccess: false,
			},
			{
				query: `bech32_address('okp4', [163,167,23,244,162,175,49,162,170,15,181,141,68,134,141,168,18,56,247,30], Bech32).`,
				wantResult: []types.TermResults{{
					"Bech32": "okp415wn30a9z4uc692s0kkx5fp5d4qfr3ac7sj9dqn",
				}},
				wantSuccess: true,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := tmdb.NewMemDB()
					stateStore := store.NewCommitMultiStore(db)
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewInterpreterMust(ctx)
						interpreter.Register3(engine.NewAtom("bech32_address"), Bech32Address)

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
