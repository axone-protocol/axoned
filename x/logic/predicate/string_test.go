//nolint:gocognit
package predicate

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
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

func TestStringBytes(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program     string
			query       string
			wantError   error
			wantSuccess bool
		}{ /*
				// inspired from https://github.com/SWI-Prolog/swipl-devel/blob/V9.1.21/src/Tests/core/test_string.pl#L91
				{
					query:       "string_bytes(aap, [97, 97, 112], ascii).",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(aap, B, utf8), B == [97, 97, 112].`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(S, [97, 97, 112], utf8), S == "aap".`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(aap, B, 'utf-16be'), B == [0, 97, 0, 97, 0, 112].`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(S, [0, 97, 0, 97, 0, 112], 'utf-16be'), S == "aap".`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(aap, B, 'utf-16le'), B ==[97, 0, 97, 0, 112, 0].`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(S, [97, 0, 97, 0, 112, 0], 'utf-16le'), S == "aap".`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(今日は, B, utf8), B == [228,187,138,230,151,165,227,129,175].`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(S, [228,187,138,230,151,165,227,129,175], utf8), S == "今日は".`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(今日は, B, 'utf-16le'), B == [202,78,229,101,111,48].`,
					query:       "test.",
					wantSuccess: true,
				},
				{
					program:     `test :- string_bytes(S, [202,78,229,101,111,48], 'utf-16le'), S == "今日は".`,
					query:       "test.",
					wantSuccess: true,
				},
				// error cases

				{
					query:       `string_bytes(_, [202,78,229,101,111,48], foo).`,
					wantSuccess: false,
					wantError:   fmt.Errorf("string_bytes/3: invalid encoding: foo"),
				},*/
			{
				query:       `string_bytes(_, _, foo).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("string_bytes/3: error(instantiation_error,string_bytes/3)"),
			}, /*
				{
					query:       `string_bytes(_, wtf, utf8).`,
					wantSuccess: false,
					wantError:   fmt.Errorf("string_bytes/3: error(type_error(list,wtf),string_bytes/3)"),
				},
				{
					query:       `string_bytes(foo(bar), _, utf8).`,
					wantSuccess: false,
					wantError:   fmt.Errorf("string_bytes/3: invalid compound term: expected a list of character_code or integer"),
				},
				{
					query:       `string_bytes(_, foo(bar), utf8).`,
					wantSuccess: false,
					wantError:   fmt.Errorf("string_bytes/3: error(type_error(list,foo(bar)),string_bytes/3)"),
				},*/
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := tmdb.NewMemDB()
					stateStore := store.NewCommitMultiStore(db)
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewLightInterpreterMust(ctx)
						interpreter.Register3(engine.NewAtom("string_bytes"), StringBytes)

						Convey("and a program", func() {
							err := interpreter.Compile(ctx, tc.program)
							So(err, ShouldBeNil)

							Convey("When the predicate is called", func() {
								sols, err := interpreter.QueryContext(ctx, tc.query)
								Reset(func() {
									So(sols.Close(), ShouldBeNil)
								})

								Convey("Then the error should be nil", func() {
									So(err, ShouldBeNil)
									So(sols, ShouldNotBeNil)

									Convey("and the result should be as expected", func() {
										if tc.wantError != nil {
											sols.Next()
											So(sols.Err(), ShouldNotBeNil)
											So(sols.Err().Error(), ShouldEqual, tc.wantError.Error())
										} else {
											nb := 0
											for sols.Next() {
												m := types.TermResults{}
												So(sols.Scan(m), ShouldBeNil)
												nb++
											}
											So(sols.Err(), ShouldBeNil)
											if tc.wantSuccess {
												So(nb, ShouldEqual, 1)
											} else {
												So(nb, ShouldEqual, 0)
											}
										}
									})
								})
							})
						})
					})
				})
			})
		}
	})
}
