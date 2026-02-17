//nolint:gocognit
package predicate

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"
	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/testutil"
)

func TestTermToAtom(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			query       string
			wantResult  []testutil.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				query:     "term_to_atom(X, Y).",
				wantError: fmt.Errorf("error(instantiation_error,term_to_atom/2)"),
			},
			{
				query:       "term_to_atom(foo, X).",
				wantResult:  []testutil.TermResults{{"X": "foo"}},
				wantSuccess: true,
			},
			{
				query:       "term_to_atom(42, X).",
				wantResult:  []testutil.TermResults{{"X": "'42'"}},
				wantSuccess: true,
			},
			{
				query:       "term_to_atom(3.14159, X).",
				wantResult:  []testutil.TermResults{{"X": "'3.14159'"}},
				wantSuccess: true,
			},
			{
				query:       "term_to_atom(-0.5, X).",
				wantResult:  []testutil.TermResults{{"X": "'-0.5'"}},
				wantSuccess: true,
			},
			{
				query:       "term_to_atom(-1.3E-14, X).",
				wantResult:  []testutil.TermResults{{"X": "'-1.3e-14'"}},
				wantSuccess: true,
			},
			{
				query:       "term_to_atom(\"hello, world\", X).",
				wantResult:  []testutil.TermResults{{"X": `'[h,e,l,l,o,\',\',\' \',w,o,r,l,d]'`}},
				wantSuccess: true,
			},
			{
				query:       "term_to_atom(X, foo).",
				wantResult:  []testutil.TermResults{{"X": "foo"}},
				wantSuccess: true,
			},
			{
				query:     "term_to_atom(X, 42).",
				wantError: fmt.Errorf("error(type_error(atom,42),term_to_atom/2)"),
			},
			{
				query:       `term_to_atom(X, '"foo"').`,
				wantResult:  []testutil.TermResults{{"X": "[f,o,o]"}},
				wantSuccess: true,
			},
			{
				query:     "term_to_atom(X, '1/2').",
				wantError: fmt.Errorf("error(syntax_error(unexpected token: graphic(/)),term_to_atom/2)"),
			},
			{
				query:      "term_to_atom(\"foo\", foo).",
				wantResult: []testutil.TermResults{},
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
						interpreter.Register2(engine.NewAtom("term_to_atom"), TermToAtom)

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

func TestAtomicListConcat(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			query       string
			wantResult  []testutil.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				query:     "atomic_list_concat(X, Y).",
				wantError: fmt.Errorf("error(instantiation_error,atomic_list_concat/2)"),
			},
			{
				query:     "atomic_list_concat([1,2], Y, Z).",
				wantError: fmt.Errorf("error(instantiation_error,atomic_list_concat/3)"),
			},
			{
				query:       "atomic_list_concat([], X).",
				wantResult:  []testutil.TermResults{{"X": "''"}},
				wantSuccess: true,
			},
			{
				query:       "atomic_list_concat([], '', X).",
				wantResult:  []testutil.TermResults{{"X": "''"}},
				wantSuccess: true,
			},
			{
				query:       "atomic_list_concat([a, 42], X).",
				wantResult:  []testutil.TermResults{{"X": "a42"}},
				wantSuccess: true,
			},
			{
				query:       "atomic_list_concat([a, '=', 42], X).",
				wantResult:  []testutil.TermResults{{"X": "'a=42'"}},
				wantSuccess: true,
			},
			{
				query:       "atomic_list_concat([a, 42], '=', X).",
				wantResult:  []testutil.TermResults{{"X": "'a=42'"}},
				wantSuccess: true,
			},
			{
				query:       "atomic_list_concat([a, '=', 42], ' ', X).",
				wantResult:  []testutil.TermResults{{"X": "'a = 42'"}},
				wantSuccess: true,
			},
			{
				query:       "atomic_list_concat([a,b], X).",
				wantResult:  []testutil.TermResults{{"X": "ab"}},
				wantSuccess: true,
			},
			{
				query:       `atomic_list_concat(["a","b"], X).`,
				wantResult:  []testutil.TermResults{{"X": "'[a][b]'"}},
				wantSuccess: true,
			},
			{
				query:     `atomic_list_concat([a,_X,c], X).`,
				wantError: fmt.Errorf("error(instantiation_error,atomic_list_concat/2)"),
			},
			{
				query:     `atomic_list_concat(foo, X).`,
				wantError: fmt.Errorf("error(type_error(list,foo),atomic_list_concat/2)"),
			},
			{
				query: "atomic_list_concat([a,b,c], foo).",
			},
			{
				query:       "atomic_list_concat([a,b,c], abc).",
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
						interpreter.Register2(engine.NewAtom("atomic_list_concat"), AtomicListConcat2)
						interpreter.Register3(engine.NewAtom("atomic_list_concat"), AtomicListConcat3)

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
