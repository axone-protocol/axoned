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

func TestURIEncoded(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program     string
			query       string
			wantResult  []types.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				query:       `uri_encoded(hey, foo, Decoded).`,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(uri_component,hey),uri_encoded/3)"),
			},
			{
				query:       `uri_encoded(path, Decoded, foo).`,
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Decoded": "foo",
				}},
			},
			{
				query:       `uri_encoded(path, Decoded, 'foo%20bar').`,
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Decoded": "'foo bar'",
				}},
			},
			{
				query:       `uri_encoded(path, foo, Encoded).`,
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Encoded": "foo",
				}},
			},
			{
				query:       `uri_encoded(query_value, 'foo bar', Encoded).`,
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Encoded": "'foo%20bar'",
				}},
			},
			{
				query:       "uri_encoded(query_value, ' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~', Encoded).",
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Encoded": "'%20!%22%23$%25%26\\'()*%2B,-./0123456789%3A%3B%3C%3D%3E?@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~'",
				}},
			},
			{
				query:       "uri_encoded(path, ' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~', Encoded).",
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Encoded": "'%20!%22%23$%25&\\'()*+,-./0123456789%3A;%3C=%3E%3F@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~'",
				}},
			},
			{
				query:       "uri_encoded(segment, ' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~', Encoded).",
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Encoded": "'%20!%22%23$%25&\\'()*+,-.%2F0123456789%3A;%3C=%3E%3F@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~'",
				}},
			},
			{
				query:       "uri_encoded(fragment, ' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~', Encoded).",
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Encoded": "'%20!%22%23$%25&\\'()*+,-./0123456789:;%3C=%3E?@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~'",
				}},
			},
			{
				query:       "uri_encoded(query_value, Decoded, '%20!%22%23$%25%26\\'()*%2B,-./0123456789%3A%3B%3C%3D%3E?@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~+').",
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Decoded": "' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~+'",
				}},
			},
			{
				query:       "uri_encoded(path, Decoded, '%20!%22%23$%25&\\'()*+,-./0123456789%3A;%3C=%3E%3F@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~').",
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Decoded": "' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~'",
				}},
			},
			{
				query:       "uri_encoded(segment, Decoded, '%20!%22%23$%25&\\'()*+,-.%2F0123456789%3A;%3C=%3E%3F@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~').",
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Decoded": "' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~'",
				}},
			},
			{
				query:       "uri_encoded(fragment, Decoded, '%20!%22%23$%25&\\'()*+,-./0123456789:;%3C=%3E?@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~').",
				wantSuccess: true,
				wantResult: []types.TermResults{{
					"Decoded": "' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~'",
				}},
			},
			{
				query:       "uri_encoded(query_value, ' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~+', '%20!%22%23$%25%26\\'()*%2B,-./0123456789%3A%3B%3C%3D%3E?@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~%2B').",
				wantSuccess: true,
				wantResult:  []types.TermResults{{}},
			},
			{
				query:       "uri_encoded(path, ' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~', '%20!%22%23$%25&\\'()*+,-./0123456789%3A;%3C=%3E%3F@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~').",
				wantSuccess: true,
				wantResult:  []types.TermResults{{}},
			},
			{
				query:       "uri_encoded(segment, ' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~', '%20!%22%23$%25&\\'()*+,-.%2F0123456789%3A;%3C=%3E%3F@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~').",
				wantSuccess: true,
				wantResult:  []types.TermResults{{}},
			},
			{
				query:       "uri_encoded(fragment, ' !\"#$%&\\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\\\]^_`abcdefghijklmnopqrstuvwxyz{|}~', '%20!%22%23$%25&\\'()*+,-./0123456789:;%3C=%3E?@ABCDEFGHIJKLMNOPQRSTUVWXYZ%5B%5C%5D%5E_%60abcdefghijklmnopqrstuvwxyz%7B%7C%7D~').",
				wantSuccess: true,
				wantResult:  []types.TermResults{{}},
			},
			{
				query:       "uri_encoded(fragment, 'foo bar', 'bar%20foo').",
				wantSuccess: false,
			},
			{
				query:       "uri_encoded(Var, 'foo bar', 'bar%20foo').",
				wantSuccess: false,
				wantError:   fmt.Errorf("error(instantiation_error,uri_encoded/3)"),
			},
			{
				query:       "uri_encoded(path, compound(2), 'bar%20foo').",
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(text,compound(2)),uri_encoded/3)"),
			},
			{
				query:       "uri_encoded(path, 'foo', compound(2)).",
				wantSuccess: false,
				wantResult:  []types.TermResults{{}},
			},
			{
				query:       "uri_encoded(path, X, compound(2)).",
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(text,compound(2)),uri_encoded/3)"),
			},
			{
				query:       "uri_encoded(path, Decoded, 'bar%%3foo').",
				wantSuccess: false,
				wantError: fmt.Errorf("error(domain_error(encoding(uri),bar%%%%3foo),[%s],uri_encoded/3)",
					strings.Join(strings.Split("invalid URL escape \"%%3\"", ""), ",")),
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
						interpreter.Register3(engine.NewAtom("uri_encoded"), URIEncoded)

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
