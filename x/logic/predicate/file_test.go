//nolint:gocognit
package predicate

import (
	goctx "context"
	"fmt"
	"net/url"
	"testing"
	"time"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/fs"
	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSourceFile(t *testing.T) {
	Convey("Given test cases", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			interpreter func(ctx goctx.Context) (i *prolog.Interpreter)
			query       string
			wantResult  []types.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "source_file(file).",
				wantSuccess: false,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "consult(file1), consult(file2), source_file(file1).",
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "consult(file1), consult(file2), consult(file3), source_file(file2).",
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "consult(file1), consult(file2), source_file(file3).",
				wantSuccess: false,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "source_file(X).",
				wantSuccess: false,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "consult(file1), consult(file2), source_file(X).",
				wantResult:  []types.TermResults{{"X": "file1"}, {"X": "file2"}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "consult(file2), consult(file3), consult(file1), source_file(X).",
				wantResult:  []types.TermResults{{"X": "file1"}, {"X": "file2"}, {"X": "file3"}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "source_file(foo(bar)).",
				wantResult:  []types.TermResults{},
				wantError:   fmt.Errorf("source_file/1: cannot unify file with *engine.compound"),
			},

			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "source_files([file]).",
				wantSuccess: false,
			},
			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "consult(file1), consult(file2), source_files([file1, file2]).",
				wantResult:  []types.TermResults{{}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "consult(file1), consult(file2), source_files([file1, file2, file3]).",
				wantSuccess: false,
			},
			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "source_files(X).",
				wantSuccess: false,
			},
			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "consult(file2), consult(file1), source_files(X).",
				wantResult:  []types.TermResults{{"X": "[file1,file2]"}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "consult(file2), consult(file1), source_files([file1, X]).",
				wantResult:  []types.TermResults{{"X": "file2"}},
				wantSuccess: true,
			},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a mocked file system", func() {
					uri, _ := url.Parse("file://dump.pl")
					mockedFS := testutil.NewMockFS(ctrl)
					mockedFS.EXPECT().Open(gomock.Any()).AnyTimes().Return(fs.NewVirtualFile(
						[]byte("dumb(dumber)."),
						uri,
						time.Now(),
					), nil)

					Convey("and a context", func() {
						db := tmdb.NewMemDB()
						stateStore := store.NewCommitMultiStore(db)
						ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

						Convey("and a vm", func() {
							interpreter := tc.interpreter(ctx)
							interpreter.FS = mockedFS
							interpreter.Register1(engine.NewAtom("source_file"), SourceFile)

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
			})
		}
	})
}
