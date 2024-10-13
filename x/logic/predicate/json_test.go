//nolint:gocognit,lll,nestif
package predicate

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/axone-protocol/prolog/engine"
	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v10/x/logic/prolog"
	"github.com/axone-protocol/axoned/v10/x/logic/testutil"
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
				description: "incorrect 1st argument",
				query:       `json_prolog(ooo(r), Term).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(text,ooo(r)),json_prolog/2)"),
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
					"Term": "json([foo=bar])",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object (given as string) into prolog",
				query:       `json_prolog("{\"foo\": \"bar\"}", Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "json([foo=bar])",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with multiple attribute into prolog",
				query:       `json_prolog('{"foo": "bar", "foobar": "bar foo"}', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "json([foo=bar,foobar='bar foo'])",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with attribute with a space into prolog",
				query:       `json_prolog('{"string with space": "bar"}', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "json(['string with space'=bar])",
				}},
				wantSuccess: true,
			},
			{
				description: "ensure prolog encoded json follows same order as json",
				query:       `json_prolog('{"b": "a", "a": "b"}', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "json([b=a,a=b])",
				}},
				wantSuccess: true,
			},
			// ** JSON -> Prolog **
			// Number
			{
				description: "convert json 0 number into prolog",
				query:       `json_prolog('0', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "0.0",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json 10 number into prolog",
				query:       `json_prolog('10', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "10.0",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json -10.9 number into prolog",
				query:       `json_prolog('-10.9', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "-10.9",
				}},
				wantSuccess: true,
			},
			{
				description: "convert large json number into prolog",
				query:       `json_prolog('100000000000000000000', Term).`,
				wantResult: []testutil.TermResults{{
					"Term": "100000000000000000000.0",
				}},
				wantSuccess: true,
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
					"Term": "[]",
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
			// ** JSON -> Prolog **
			// Pathological
			{
				description: "convert an object with an invalid key type (numeric) to Prolog",
				query:       `json_prolog('{5:"bar"}', Term).`,
				wantError:   fmt.Errorf("error(syntax_error(json(malformed_json(1))),json_prolog/2)"),
				wantSuccess: false,
			},
			{
				description: "convert incorrect json into prolog",
				query:       `json_prolog('@wtf!', Term).`,
				wantError:   fmt.Errorf("error(syntax_error(json(malformed_json(1))),json_prolog/2)"),
				wantSuccess: false,
			},
			{
				description: "convert large json number with ridonculous exponent into prolog",
				query:       `json_prolog('1E30923434', Term).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(syntax_error(json(malformed_json(11,number 1E30923434))),json_prolog/2)"),
			},
			{
				description: "convert unfinished json into prolog",
				query:       `json_prolog('{"foo": ', Term).`,
				wantError:   fmt.Errorf("error(syntax_error(json(eof)),json_prolog/2)"),
				wantSuccess: false,
			},
			{
				description: "check json array is well formed",
				query:       `json_prolog('[&', Term).`,
				wantError:   fmt.Errorf("error(syntax_error(json(malformed_json(1))),json_prolog/2)"),
				wantSuccess: false,
			},
			{
				description: "check json object is well formed (1)",
				query:       `json_prolog('{"foo": "bar"}{"foo": "bar"}', Term).`,
				wantError:   fmt.Errorf("error(syntax_error(json(malformed_json(15))),json_prolog/2)"),
				wantSuccess: false,
			},
			{
				description: "check json object is well formed (2)",
				query:       `json_prolog('{&', Term).`,
				wantError:   fmt.Errorf("error(syntax_error(json(malformed_json(1))),json_prolog/2)"),
				wantSuccess: false,
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
				description: "convert utf8 string to json",
				query:       `json_prolog(Json, json([foo='今日は'])).`,
				wantResult: []testutil.TermResults{{
					"Json": `'{"foo":"今日は"}'`,
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
			// ** Prolog -> JSON **
			// Object
			{
				description: "convert json object from prolog",
				query:       `json_prolog(Json, json([foo=bar])).`,
				wantResult: []testutil.TermResults{{
					"Json": "'{\"foo\":\"bar\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with multiple attribute from prolog",
				query:       `json_prolog(Json, json([foo=bar,foobar='bar foo'])).`,
				wantResult: []testutil.TermResults{{
					"Json": "'{\"foo\":\"bar\",\"foobar\":\"bar foo\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert json object with attribute with a space into prolog",
				query:       `json_prolog(Json, json(['string with space'=bar])).`,
				wantResult: []testutil.TermResults{{
					"Json": "'{\"string with space\":\"bar\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "ensure json follows same order as prolog encoded",
				query:       `json_prolog(Json, json([b=a,a=b])).`,
				wantResult: []testutil.TermResults{{
					"Json": "'{\"b\":\"a\",\"a\":\"b\"}'",
				}},
				wantSuccess: true,
			},
			{
				description: "invalid json term compound",
				query:       `json_prolog(Json, foo([a=b])).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(json,foo([=(a,b)])),json_prolog/2)"),
			},
			{
				description: "convert json term object from prolog with error inside",
				query:       `json_prolog(Json, ['string with space',json('toto')]).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(list,toto),json_prolog/2)"),
			},
			{
				description: "convert json term object from prolog with error inside another object",
				query:       `json_prolog(Json, ['string with space',json([key=json(error)])]).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(list,error),json_prolog/2)"),
			},
			{
				description: "convert json term object which incorrectly defines key/value pair",
				query:       `json_prolog(Json, json([not_a_key_value(key,value)])).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(key_value,not_a_key_value(key,value)),json_prolog/2)"),
			},
			{
				description: "convert json term object which uses a non atom for key",
				query:       `json_prolog(Json, json([=(42,value)])).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(atom,42),json_prolog/2)"),
			},
			{
				description: "convert json term object with arity > 2",
				query:       `json_prolog(Json, json(a,b,c)).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(json,json(a,b,c)),json_prolog/2)"),
			},
			// ** Prolog -> JSON **
			// Number
			{
				description: "convert prolog 0 number",
				query:       `json_prolog(Json, 0).`,
				wantResult: []testutil.TermResults{{
					"Json": "'0'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert prolog array of numbers",
				query:       `json_prolog(Json, [1, 2, 3]).`,
				wantResult: []testutil.TermResults{{
					"Json": "'[1,2,3]'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert prolog 10 number",
				query:       `json_prolog(Json, 10).`,
				wantResult: []testutil.TermResults{{
					"Json": "'10'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert prolog decimal 10.4 number",
				query:       `json_prolog(Json, 10.4).`,
				wantResult: []testutil.TermResults{{
					"Json": "'10.4'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert prolog decimal -10.4 number",
				query:       `json_prolog(Json, -10.4).`,
				wantResult: []testutil.TermResults{{
					"Json": "'-10.4'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert big prolog decimal",
				query:       `json_prolog(Json, 100000000000000000000.0).`,
				wantResult: []testutil.TermResults{{
					"Json": "'100000000000000000000'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert prolog decimal with exponent",
				query:       `json_prolog(Json, 1.0E99).`,
				wantResult: []testutil.TermResults{{
					"Json": "'1e+99'",
				}},
				wantSuccess: true,
			},
			{
				description: "convert prolog decimal with ridonculous exponent",
				query:       `json_prolog(Json, 1.8e308).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(domain_error(json_number,1.8e+308),json_prolog/2)"),
			},
			// ** Prolog -> Json **
			// Array
			{
				description: "convert empty json array from prolog",
				query:       `json_prolog(Json, []).`,
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
			// ** JSON <-> Prolog **
			{
				description: "ensure unification doesn't depend on formatting",
				query:       `json_prolog('{\n\t"foo": "bar"\n}', json( [ foo  =  bar ] )).`,
				wantResult:  []testutil.TermResults{{}},
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
						interpreter.Register2(engine.NewAtom("json_read"), JSONRead)

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
				term:        "json([foo=bar])",
				wantSuccess: true,
			},
			{
				json:        "'{\"foo\":\"null\"}'",
				term:        "json([foo=null])",
				wantSuccess: true,
			},
			{
				json:        "'{\"foo\":null}'",
				term:        "json([foo= @(null)])",
				wantSuccess: true,
			},
			{
				json:        "'{\"employee\":{\"age\":30,\"city\":\"New York\",\"name\":\"John\"}}'",
				term:        "json([employee=json([age=30.0,city='New York',name='John'])])",
				wantSuccess: true,
			},
			{
				json:        "'{\"cosmos\":[\"axone\",{\"name\":\"localnet\"}]}'",
				term:        "json([cosmos=[axone,json([name=localnet])]])",
				wantSuccess: true,
			},
			{
				json:        "'{\"object\":{\"array\":[1,2,3],\"arrayobject\":[{\"name\":\"toto\"},{\"name\":\"tata\"}],\"bool\":true,\"boolean\":false,\"null\":null}}'",
				term:        "json([object=json([array=[1.0,2.0,3.0],arrayobject=[json([name=toto]),json([name=tata])],bool= @(true),boolean= @(false),null= @(null)])])",
				wantSuccess: true,
			},
			{
				json:        "'{\"foo\":\"bar\"}'",
				term:        "json([a=b])",
				wantSuccess: false,
			},
			{
				json:        `'{"key1":null,"key2":[],"key3":{"nestedKey1":null,"nestedKey2":[],"nestedKey3":["a",null,null]}}'`,
				term:        `json([key1= @(null),key2=[],key3=json([nestedKey1= @(null),nestedKey2=[],nestedKey3=[a,@(null),@(null)]])])`,
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

func TestJSONRead(t *testing.T) {
	Convey("Given a test cases for JSONRead", t, func() {
		cases := []struct {
			stream         func() engine.Term
			wantSuccess    bool
			wantErrorMatch string
		}{
			{
				stream:         func() engine.Term { return engine.NewInputBinaryStream(strings.NewReader("{}")) },
				wantErrorMatch: `error\(permission_error\(input,text_stream,<stream>\(0x[[:xdigit:]]+\)\),root\)`,
			},
			{
				stream: func() engine.Term {
					var buf bytes.Buffer
					return engine.NewOutputTextStream(&buf)
				},
				wantErrorMatch: `error\(permission_error\(input,stream,<stream>\(0x[[:xdigit:]]+\)\),root\)`,
			},
			{
				stream: func() engine.Term {
					return engine.NewAtom("not a stream")
				},
				wantErrorMatch: `error\(type_error\(stream,not a stream\),root\)`,
			},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given a stream (#%d)", nc), func() {
				stream := tc.stream()
				defer func() {
					if s, ok := stream.(*engine.Stream); ok {
						_ = s.Close()
					}
				}()

				Convey("and an interpreter", func() {
					ctx := context.Background()
					interpreter := testutil.NewLightInterpreterMust(ctx)
					env := engine.NewEnv()

					Convey("When the predicate JSONRead is called", func() {
						got := engine.NewVariable()
						ok, err := JSONRead(&interpreter.VM, stream, got, engine.Success, env).Force(ctx)

						Convey("Then the result should be as expected", func() {
							So(ok, ShouldEqual, tc.wantSuccess)
							if tc.wantErrorMatch != "" {
								So(err, ShouldNotBeNil)
								SoMsg(fmt.Sprintf("%s ~ %s", err, tc.wantErrorMatch),
									regexp.MustCompile(tc.wantErrorMatch).FindStringIndex(err.Error()), ShouldNotBeNil)
							}
						})
					})
				})
			})
		}
	})
}

func TestJSONWrite(t *testing.T) {
	Convey("Given a test cases for JSONWrite", t, func() {
		cases := []struct {
			stream         func() engine.Term
			wantSuccess    bool
			wantErrorMatch string
		}{
			{
				stream: func() engine.Term {
					var buf bytes.Buffer
					return engine.NewOutputBinaryStream(&buf)
				},
				wantErrorMatch: `error\(permission_error\(output,text_stream,<stream>\(0x[[:xdigit:]]+\)\),root\)`,
			},
			{
				stream: func() engine.Term {
					return engine.NewInputBinaryStream(strings.NewReader("{}"))
				},
				wantErrorMatch: `error\(permission_error\(output,stream,<stream>\(0x[[:xdigit:]]+\)\),root\)`,
			},
			{
				stream: func() engine.Term {
					return engine.NewAtom("not a stream")
				},
				wantErrorMatch: `error\(type_error\(stream,not a stream\),root\)`,
			},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given a stream (#%d)", nc), func() {
				stream := tc.stream()
				defer func() {
					if s, ok := stream.(*engine.Stream); ok {
						_ = s.Close()
					}
				}()

				Convey("and an interpreter", func() {
					ctx := context.Background()
					interpreter := testutil.NewLightInterpreterMust(ctx)
					env := engine.NewEnv()

					Convey("When the predicate JSONWrite is called", func() {
						term := prolog.AtomJSON.Apply(
							engine.List(prolog.AtomKeyValue.Apply(prolog.StringToAtom("key"), prolog.StringToAtom("value"))))
						ok, err := JSONWrite(&interpreter.VM, stream, term, engine.Success, env).Force(ctx)

						Convey("Then the result should be as expected", func() {
							So(ok, ShouldEqual, tc.wantSuccess)
							if tc.wantErrorMatch != "" {
								So(err, ShouldNotBeNil)
								SoMsg(fmt.Sprintf("%s ~ %s", err, tc.wantErrorMatch),
									regexp.MustCompile(tc.wantErrorMatch).FindStringIndex(err.Error()), ShouldNotBeNil)
							}
						})
					})
				})
			})
		}
	})
}
