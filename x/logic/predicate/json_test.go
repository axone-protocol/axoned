//nolint:gocognit,lll
package predicate

import (
	"fmt"
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

func TestJsonProlog(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			description string
			program     string
			query       string
			wantResult  []types.TermResults
			wantError   error
			wantSuccess bool
		}{
			// ** JSON -> Prolog **
			// String
			{
				description: "convert direct string (valid json) into prolog",
				query:       `json_prolog('"foo"', Term).`,
				wantResult: []types.TermResults{{
					"Term": "foo",
				}},
				wantSuccess: true,
			},
			{
				description: "convert direct string with space (valid json) into prolog",
				query:       `json_prolog('"a string with space"', Term).`,
				wantResult: []types.TermResults{{
					"Term": "'a string with space'",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Object
			{
				description: "convert json object into prolog",
				query:       `json_prolog('{"foo": "bar"}', Term).`,
				wantResult: []types.TermResults{{
					"Term": "json([foo-bar])",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with multiple attribute into prolog",
				query:       `json_prolog('{"foo": "bar", "foobar": "bar foo"}', Term).`,
				wantResult: []types.TermResults{{
					"Term": "json([foo-bar,foobar-'bar foo'])",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with attribute with a space into prolog",
				query:       `json_prolog('{"string with space": "bar"}', Term).`,
				wantResult: []types.TermResults{{
					"Term": "json(['string with space'-bar])",
				}},
				wantSuccess: true,
			},
			{
				description: "ensure determinism on object attribute key sorted alphabetically",
				query:       `json_prolog('{"b": "a", "a": "b"}', Term).`,
				wantResult: []types.TermResults{{
					"Term": "json([a-b,b-a])",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Number
			{
				description: "convert json number into prolog",
				query:       `json_prolog('10', Term).`,
				wantResult: []types.TermResults{{
					"Term": "10",
				}},
				wantSuccess: true,
			},
			{
				description: "convert large json number into prolog",
				query:       `json_prolog('100000000000000000000', Term).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("json_prolog/2: could not convert number '100000000000000000000' into integer term, overflow"),
			},
			{
				description: "decimal number not compatible yet",
				query:       `json_prolog('10.4', Term).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("json_prolog/2: could not convert number '10.4' into integer term, decimal number is not handled yet"),
			},
			// ** JSON -> Prolog **
			// Bool
			{
				description: "convert json true boolean into prolog",
				query:       `json_prolog('true', Term).`,
				wantResult: []types.TermResults{{
					"Term": "@(true)",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json false boolean into prolog",
				query:       `json_prolog('false', Term).`,
				wantResult: []types.TermResults{{
					"Term": "@(false)",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Null
			{
				description: "convert json null value into prolog",
				query:       `json_prolog('null', Term).`,
				wantResult: []types.TermResults{{
					"Term": "@(null)",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Array
			{
				description: "convert json array into prolog",
				query:       `json_prolog('["foo", "bar"]', Term).`,
				wantResult: []types.TermResults{{
					"Term": "[foo,bar]",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json string array into prolog",
				query:       `json_prolog('["string with space", "bar"]', Term).`,
				wantResult: []types.TermResults{{
					"Term": "['string with space',bar]",
				}},
				wantSuccess: true,
			},

			// ** Prolog -> JSON **
			// String
			{
				description: "convert string term to json",
				query:       `json_prolog(Json, 'foo').`,
				wantResult: []types.TermResults{{
					"Json": "'\"foo\"'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert string with space to json",
				query:       `json_prolog(Json, 'foo bar').`,
				wantResult: []types.TermResults{{
					"Json": "'\"foo bar\"'",
				}},
				wantSuccess: true,
			},
			// ** Prolog -> JSON **
			// Object
			{
				description: "convert json object from prolog",
				query:       `json_prolog(Json, json([foo-bar])).`,
				wantResult: []types.TermResults{{
					"Json": "'{\"foo\":\"bar\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with multiple attribute from prolog",
				query:       `json_prolog(Json, json([foo-bar,foobar-'bar foo'])).`,
				wantResult: []types.TermResults{{
					"Json": "'{\"foo\":\"bar\",\"foobar\":\"bar foo\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with attribute with a space into prolog",
				query:       `json_prolog(Json, json(['string with space'-bar])).`,
				wantResult: []types.TermResults{{
					"Json": "'{\"string with space\":\"bar\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "ensure determinism on object attribute key sorted alphabetically",
				query:       `json_prolog(Json, json([b-a,a-b])).`,
				wantResult: []types.TermResults{{
					"Json": "'{\"a\":\"b\",\"b\":\"a\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "invalid json term compound",
				query:       `json_prolog(Json, foo([a-b])).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("json_prolog/2: invalid functor foo"),
			},
			// ** Prolog -> JSON **
			// Number
			{
				description: "convert json number from prolog",
				query:       `json_prolog(Json, 10).`,
				wantResult: []types.TermResults{{
					"Json": "'10'",
				}},
				wantSuccess: true,
			},
			{
				description: "decimal number not compatible yet",
				query:       `json_prolog(Json, 10.4).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("json_prolog/2: could not convert %%!s(engine.Float=10.4) {engine.Float} to json"),
			},
			// ** Prolog -> Json **
			// Array
			{
				description: "convert json array from prolog",
				query:       `json_prolog(Json, [foo,bar]).`,
				wantResult: []types.TermResults{{
					"Json": "'[\"foo\",\"bar\"]'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json string array from prolog",
				query:       `json_prolog(Json, ['string with space',bar]).`,
				wantResult: []types.TermResults{{
					"Json": "'[\"string with space\",\"bar\"]'",
				}},
				wantSuccess: true,
			},
			// ** Prolog -> JSON **
			// Bool
			{
				description: "convert true boolean from prolog",
				query:       `json_prolog(Json, @(true)).`,
				wantResult: []types.TermResults{{
					"Json": "true",
				}},
				wantSuccess: true,
			},
			{
				description: "convert false boolean from prolog",
				query:       `json_prolog(Json, @(false)).`,
				wantResult: []types.TermResults{{
					"Json": "false",
				}},
				wantSuccess: true,
			},
			// ** Prolog -> Json **
			// Null
			{
				description: "convert json null value into prolog",
				query:       `json_prolog(Json, @(null)).`,
				wantResult: []types.TermResults{{
					"Json": "null",
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
						interpreter := testutil.NewLightInterpreterMust(ctx)
						interpreter.Register2(engine.NewAtom("json_prolog"), JsonProlog)

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

func TestJsonPrologWithMoreComplexStructBidirectional(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			json        string
			term        string
			wantError   error
			wantSuccess bool
		}{
			{
				json:        "'{\"foo\":\"bar\"}'",
				term:        "json([foo-bar])",
				wantSuccess: true,
			},
			{
				json:        "'{\"employee\":{\"age\":30,\"city\":\"New York\",\"name\":\"John\"}}'",
				term:        "json([employee-json([age-30,city-'New York',name-'John'])])",
				wantSuccess: true,
			},
			{
				json:        "'{\"cosmos\":[\"okp4\",{\"name\":\"localnet\"}]}'",
				term:        "json([cosmos-[okp4,json([name-localnet])]])",
				wantSuccess: true,
			},
			{
				json:        "'{\"object\":{\"array\":[1,2,3],\"arrayobject\":[{\"name\":\"toto\"},{\"name\":\"tata\"}],\"bool\":true,\"boolean\":false,\"null\":null}}'",
				term:        "json([object-json([array-[1,2,3],arrayobject-[json([name-toto]),json([name-tata])],bool- @(true),boolean- @(false),null- @(null)])])",
				wantSuccess: true,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("#%d : given the json: %s and the term %s", nc, tc.json, tc.term), func() {
				Convey("and a context", func() {
					db := tmdb.NewMemDB()
					stateStore := store.NewCommitMultiStore(db)
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewLightInterpreterMust(ctx)
						interpreter.Register2(engine.NewAtom("json_prolog"), JsonProlog)

						Convey("When the predicate `json_prolog` is called to convert json to prolog", func() {
							sols, err := interpreter.QueryContext(ctx, fmt.Sprintf("json_prolog(%s, Term).", tc.json))

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
											So(len(got), ShouldEqual, 1)
											for _, resultGot := range got {
												for _, termGot := range resultGot {
													So(testutil.ReindexUnknownVariables(termGot), ShouldEqual, tc.term)
												}
											}
										} else {
											So(len(got), ShouldEqual, 0)
										}
									}
								})
							})
						})

						Convey("When the predicate `json_prolog` is called to convert prolog to json", func() {
							sols, err := interpreter.QueryContext(ctx, fmt.Sprintf("json_prolog(Json, %s).", tc.term))

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
											So(len(got), ShouldEqual, 1)
											for _, resultGot := range got {
												for _, termGot := range resultGot {
													So(testutil.ReindexUnknownVariables(termGot), ShouldEqual, tc.json)
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
