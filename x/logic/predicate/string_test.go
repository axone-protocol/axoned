//nolint:gocognit
package predicate

import (
	"fmt"
	"strings"
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReadString(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			input       string
			program     string
			query       string
			wantResult  []types.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				input:   "foo",
				program: "read_input(String) :- current_input(Stream), read_string(Stream, _, String).",
				query:   `read_input(String).`,
				wantResult: []types.TermResults{{
					"String": "foo",
				}},
				wantSuccess: true,
			},
			{
				input:   "foo bar",
				program: "read_input(String) :- current_input(Stream), read_string(Stream, _, String).",
				query:   `read_input(String).`,
				wantResult: []types.TermResults{{
					"String": "'foo bar'",
				}},
				wantSuccess: true,
			},
			{
				input:   "foo bar",
				program: "read_input(String, Len) :- current_input(Stream), read_string(Stream, Len, String).",
				query:   `read_input(String, Len).`,
				wantResult: []types.TermResults{{
					"String": "'foo bar'",
					"Len":    "7",
				}},
				wantSuccess: true,
			},
			{
				input:   "foo bar",
				program: "read_input(String, Len) :- current_input(Stream), read_string(Stream, Len, String).",
				query:   `read_input(String, 3).`,
				wantResult: []types.TermResults{{
					"String": "foo",
				}},
				wantSuccess: true,
			},
			{
				input:   "foo bar",
				program: "read_input(String, Len) :- current_input(Stream), read_string(Stream, Len, String).",
				query:   `read_input(String, 7).`,
				wantResult: []types.TermResults{{
					"String": "'foo bar'",
				}},
				wantSuccess: true,
			},
			{
				input:   "foo bar 🧙",
				program: "read_input(String, Len) :- current_input(Stream), read_string(Stream, Len, String).",
				query:   `read_input(String, _).`,
				wantResult: []types.TermResults{{
					"String": "'foo bar 🧙'",
				}},
				wantSuccess: true,
			},
			{
				input:   "foo bar 🧙",
				program: "read_input(String, Len) :- current_input(Stream), read_string(Stream, Len, String).",
				query:   `read_input(String, Len).`,
				wantResult: []types.TermResults{{
					"String": "'foo bar 🧙'",
					"Len":    "12",
				}},
				wantSuccess: true,
			},
			{
				input:   "🧙",
				program: "read_input(String, Len) :- current_input(Stream), read_string(Stream, Len, String).",
				query:   `read_input(String, Len).`,
				wantResult: []types.TermResults{{
					"String": "'🧙'",
					"Len":    "4",
				}},
				wantSuccess: true,
			},
			{
				input:   "🧙",
				program: "read_input(String, Len) :- current_input(Stream), read_string(Stream, Len, String).",
				query:   `read_input(String, 1).`,
				wantResult: []types.TermResults{{
					"String": "'🧙'",
				}},
				wantSuccess: false,
			},
			{
				input:   "Hello World!",
				program: "read_input(String, Len) :- current_input(Stream), read_string(Stream, Len, String).",
				query:   `read_input(String, 15).`,
				wantResult: []types.TermResults{{
					"String": "'Hello World!'",
				}},
				wantSuccess: false,
			},
			{
				input:       "Hello World!",
				program:     "read_input(String, Len) :- current_input(Stream), read_string(foo, Len, String).",
				query:       `read_input(String, Len).`,
				wantError:   fmt.Errorf("read_string/3: invalid domain for given stream"),
				wantSuccess: false,
			},
			{
				input:       "Hello World!",
				query:       `read_string(Stream, Len, data).`,
				wantError:   fmt.Errorf("read_string/3: stream cannot be a variable"),
				wantSuccess: false,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := tmdb.NewMemDB()
					stateStore := store.NewCommitMultiStore(db)
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewComprehensiveInterpreterMust(ctx)
						interpreter.Register3(engine.NewAtom("read_string"), ReadString)

						interpreter.SetUserInput(engine.NewInputTextStream(strings.NewReader(tc.input)))

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
