//nolint:gocognit,lll,nestif
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

	"github.com/axone-protocol/axoned/v9/x/logic/testutil"
)

func TestJsonProlog(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			description string
			program     string
			query       string
			wantResult  []testutil.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				description: "two variable",
				query:       `json_prolog(Json, Term).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(instantiation_error,json_prolog/2)"),
			},
			{
				description: "two variable",
				query:       `json_prolog(ooo(r), Term).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(atom,ooo(r)),json_prolog/2)"),
			},

			// ** JSON -> Prolog **
			// String
			{
				description: "convert direct string (valid json) into prolog",
				query:       `json_prolog('"foo"', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "foo",
				}},
				wantSuccess: true,
			},
			{
				description: "convert direct string with space (valid json) into prolog",
				query:       `json_prolog('"a string with space"', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "'a string with space'",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Object
			{
				description: "convert json object into prolog",
				query:       `json_prolog('{"foo": "bar"}', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "json([foo-bar])",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with multiple attribute into prolog",
				query:       `json_prolog('{"foo": "bar", "foobar": "bar foo"}', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "json([foo-bar,foobar-'bar foo'])",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with attribute with a space into prolog",
				query:       `json_prolog('{"string with space": "bar"}', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "json(['string with space'-bar])",
				}},
				wantSuccess: true,
			},
			{
				description: "ensure determinism on object attribute key sorted alphabetically",
				query:       `json_prolog('{"b": "a", "a": "b"}', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "json([a-b,b-a])",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Number
			{
				description: "convert json number into prolog",
				query:       `json_prolog('10', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "10",
				}},
				wantSuccess: true,
			},
			{
				description: "convert large json number into prolog",
				query:       `json_prolog('100000000000000000000', Term).`,
				wantSuccess: false,
				wantError: fmt.Errorf("error(domain_error(encoding(json),100000000000000000000),[%s],json_prolog/2)",
					strings.Join(strings.Split("could not convert number '100000000000000000000' into integer term, overflow", ""), ",")),
			},
			{
				description: "decimal number not compatible yet",
				query:       `json_prolog('10.4', Term).`,
				wantSuccess: false,
				wantError: fmt.Errorf("error(domain_error(encoding(json),10.4),[%s],json_prolog/2)",
					strings.Join(strings.Split("could not convert number '10.4' into integer term, decimal number is not handled yet", ""), ",")),
			},
			// ** JSON -> Prolog **
			// Bool
			{
				description: "convert json true boolean into prolog",
				query:       `json_prolog('true', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "@(true)",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json false boolean into prolog",
				query:       `json_prolog('false', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "@(false)",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Null
			{
				description: "convert json null value into prolog",
				query:       `json_prolog('null', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "@(null)",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Array
			{
				description: "convert empty json array into prolog",
				query:       `json_prolog('[]', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "@([])",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json array into prolog",
				query:       `json_prolog('["foo", "bar"]', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "[foo,bar]",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json array with null element into prolog",
				query:       `json_prolog('[null]', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "[@(null)]",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json string array into prolog",
				query:       `json_prolog('["string with space", "bar"]', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "['string with space',bar]",
				}},
				wantSuccess: true,
			},

			// ** Prolog -> JSON **
			// String
			{
				description: "convert string term to json",
				query:       `json_prolog(Json, 'foo').`,
				wantResult: []testutil.TermResults{{
					"Json": "'\"foo\"'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert atom term to json",
				query:       `json_prolog(Json, foo).`,
				wantResult: []testutil.TermResults{{
					"Json": "'\"foo\"'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert string with space to json",
				query:       `json_prolog(Json, 'foo bar').`,
				wantResult: []testutil.TermResults{{
					"Json": "'\"foo bar\"'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert string with space to json",
				query:       `json_prolog(Json, 'foo bar').`,
				wantResult: []testutil.TermResults{{
					"Json": "'\"foo bar\"'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert empty-list atom term to json",
				query:       `json_prolog(Json, []).`,
				wantResult: []testutil.TermResults{{
					"Json": "'\"[]\"'",
				}},
				wantSuccess: true,
			},
			// ** Prolog -> JSON **
			// Object
			{
				description: "convert json object from prolog",
				query:       `json_prolog(Json, json([foo-bar])).`,
				wantResult: []testutil.TermResults{{
					"Json": "'{\"foo\":\"bar\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with multiple attribute from prolog",
				query:       `json_prolog(Json, json([foo-bar,foobar-'bar foo'])).`,
				wantResult: []testutil.TermResults{{
					"Json": "'{\"foo\":\"bar\",\"foobar\":\"bar foo\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with attribute with a space into prolog",
				query:       `json_prolog(Json, json(['string with space'-bar])).`,
				wantResult: []testutil.TermResults{{
					"Json": "'{\"string with space\":\"bar\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "ensure determinism on object attribute key sorted alphabetically",
				query:       `json_prolog(Json, json([b-a,a-b])).`,
				wantResult: []testutil.TermResults{{
					"Json": "'{\"a\":\"b\",\"b\":\"a\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "invalid json term compound",
				query:       `json_prolog(Json, foo([a-b])).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(json,foo([-(a,b)])),json_prolog/2)"),
			},
			{
				description: "convert json term object from prolog with error inside",
				query:       `json_prolog(Json, ['string with space',json('toto')]).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(list,toto),json_prolog/2)"),
			},
			{
				description: "convert json term object from prolog with error inside another object",
				query:       `json_prolog(Json, ['string with space',json([key-json(error)])]).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(list,error),json_prolog/2)"),
			},
			// ** Prolog -> JSON **
			// Number
			{
				description: "convert json number from prolog",
				query:       `json_prolog(Json, 10).`,
				wantResult: []testutil.TermResults{{
					"Json": "'10'",
				}},
				wantSuccess: true,
			},
			{
				description: "decimal number not compatible yet",
				query:       `json_prolog(Json, 10.4).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(json,10.4),json_prolog/2)"),
			},
			// ** Prolog -> Json **
			// Array
			{
				description: "convert empty json array from prolog",
				query:       `json_prolog(Json, @([])).`,
				wantResult: []testutil.TermResults{{
					"Json": "[]",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json array from prolog",
				query:       `json_prolog(Json, [foo,bar]).`,
				wantResult: []testutil.TermResults{{
					"Json": "'[\"foo\",\"bar\"]'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json array with null element from prolog",
				query:       `json_prolog(Json, [@(null)]).`,
				wantResult: []testutil.TermResults{{
					"Json": "'[null]'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json string array from prolog",
				query:       `json_prolog(Json, ['string with space',bar]).`,
				wantResult: []testutil.TermResults{{
					"Json": "'[\"string with space\",\"bar\"]'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json string array from prolog with error inside",
				query:       `json_prolog(Json, ['string with space',hey('toto')]).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(json,hey(toto)),json_prolog/2)"),
			},
			// ** Prolog -> JSON **
			// Bool
			{
				description: "convert true boolean from prolog",
				query:       `json_prolog(Json, @(true)).`,
				wantResult: []testutil.TermResults{{
					"Json": "true",
				}},
				wantSuccess: true,
			},
			{
				description: "convert false boolean from prolog",
				query:       `json_prolog(Json, @(false)).`,
				wantResult: []testutil.TermResults{{
					"Json": "false",
				}},
				wantSuccess: true,
			},
			// ** Prolog -> Json **
			// Null
			{
				description: "convert json null value into prolog",
				query:       `json_prolog(Json, @(null)).`,
				wantResult: []testutil.TermResults{{
					"Json": "null",
				}},
				wantSuccess: true,
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
						interpreter.Register2(engine.NewAtom("json_prolog"), JSONProlog)

						err := interpreter.Compile(ctx, tc.program)
						So(err, ShouldBeNil)

						Convey("When the predicate is called", func() {
							sols, err := interpreter.QueryContext(ctx, tc.query)

							Convey("Then the error should be nil", func() {
								So(err, ShouldBeNil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									var got []testutil.TermResults
									for sols.Next() {
										m := testutil.TermResults{}
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
				json:        "'{\"foo\":\"null\"}'",
				term:        "json([foo-null])",
				wantSuccess: true,
			},
			{
				json:        "'{\"foo\":null}'",
				term:        "json([foo- @(null)])",
				wantSuccess: true,
			},
			{
				json:        "'{\"employee\":{\"age\":30,\"city\":\"New York\",\"name\":\"John\"}}'",
				term:        "json([employee-json([age-30,city-'New York',name-'John'])])",
				wantSuccess: true,
			},
			{
				json:        "'{\"cosmos\":[\"axone\",{\"name\":\"localnet\"}]}'",
				term:        "json([cosmos-[axone,json([name-localnet])]])",
				wantSuccess: true,
			},
			{
				json:        "'{\"object\":{\"array\":[1,2,3],\"arrayobject\":[{\"name\":\"toto\"},{\"name\":\"tata\"}],\"bool\":true,\"boolean\":false,\"null\":null}}'",
				term:        "json([object-json([array-[1,2,3],arrayobject-[json([name-toto]),json([name-tata])],bool- @(true),boolean- @(false),null- @(null)])])",
				wantSuccess: true,
			},
			{
				json:        "'{\"foo\":\"bar\"}'",
				term:        "json([a-b])",
				wantSuccess: false,
			},
			{
				json:        `'{"key1":null,"key2":[],"key3":{"nestedKey1":null,"nestedKey2":[],"nestedKey3":["a",null,null]}}'`,
				term:        `json([key1- @(null),key2- @([]),key3-json([nestedKey1- @(null),nestedKey2- @([]),nestedKey3-[a,@(null),@(null)]])])`,
				wantSuccess: true,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("#%d : given the json: %s and the term %s", nc, tc.json, tc.term), func() {
				Convey("and a context", func() {
					db := dbm.NewMemDB()
					stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewLightInterpreterMust(ctx)
						interpreter.Register2(engine.NewAtom("json_prolog"), JSONProlog)

						if tc.wantSuccess {
							Convey("When the predicate `json_prolog` is called to convert json to prolog", func() {
								sols, err := interpreter.QueryContext(ctx, fmt.Sprintf("json_prolog(%s, Term).", tc.json))

								Convey("Then the error should be nil", func() {
									So(err, ShouldBeNil)
									So(sols, ShouldNotBeNil)

									Convey("and the bindings should be as expected", func() {
										var got []testutil.TermResults
										for sols.Next() {
											m := testutil.TermResults{}
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
														reindexedTerm := fmt.Sprintf("%v", testutil.ReindexUnknownVariables(termGot))
														So(reindexedTerm, ShouldEqual, tc.term)
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
										var got []testutil.TermResults
										for sols.Next() {
											m := testutil.TermResults{}
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
														reindexedTerm := fmt.Sprintf("%v", testutil.ReindexUnknownVariables(termGot))
														So(reindexedTerm, ShouldEqual, tc.json)
													}
												}
											} else {
												So(len(got), ShouldEqual, 0)
											}
										}
									})
								})
							})
						}

						Convey("When the predicate `json_prolog` is called to check prolog matching json", func() {
							sols, err := interpreter.QueryContext(ctx, fmt.Sprintf("json_prolog(%s, %s).", tc.json, tc.term))

							Convey("Then the error should be nil", func() {
								So(err, ShouldBeNil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									So(sols.Next(), ShouldEqual, tc.wantSuccess)

									if tc.wantError != nil {
										So(sols.Err(), ShouldNotBeNil)
										So(sols.Err().Error(), ShouldEqual, tc.wantError.Error())
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
