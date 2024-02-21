//nolint:gocognit
package predicate

import (
	goctx "context"
	"fmt"
	fs2 "io/fs"
	"net/url"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/golang/mock/gomock"
	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/logic/fs"
	"github.com/okp4/okp4d/x/logic/testutil"
)

func TestSourceFile(t *testing.T) {
	Convey("Given test cases", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			interpreter func(ctx goctx.Context) (i *prolog.Interpreter)
			query       string
			wantResult  []testutil.TermResults
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
				wantResult:  []testutil.TermResults{{}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "consult(file1), consult(file2), consult(file3), source_file(file2).",
				wantResult:  []testutil.TermResults{{}},
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
				wantResult:  []testutil.TermResults{{"X": "file1"}, {"X": "file2"}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "consult(file2), consult(file3), consult(file1), source_file(X).",
				wantResult:  []testutil.TermResults{{"X": "file1"}, {"X": "file2"}, {"X": "file3"}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewLightInterpreterMust,
				query:       "source_file(foo(bar)).",
				wantResult:  []testutil.TermResults{},
				wantError:   fmt.Errorf("error(type_error(atom,foo(bar)),source_file/1)"),
			},

			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "source_files([file]).",
				wantSuccess: false,
			},
			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "consult(file1), consult(file2), source_files([file1, file2]).",
				wantResult:  []testutil.TermResults{{}},
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
				wantResult:  []testutil.TermResults{{"X": "[file1,file2]"}},
				wantSuccess: true,
			},
			{
				interpreter: testutil.NewComprehensiveInterpreterMust,
				query:       "consult(file2), consult(file1), source_files([file1, X]).",
				wantResult:  []testutil.TermResults{{"X": "file2"}},
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
						db := dbm.NewMemDB()
						stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
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

func TestOpen(t *testing.T) {
	Convey("Given a test cases", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			files       map[string][]byte
			program     string
			query       string
			wantResult  []testutil.TermResults
			wantError   error
			wantSuccess bool
		}{
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program: "get_first_char(C) :- open(file, read, Stream, _), get_char(Stream, C).",
				query:   `get_first_char(C).`,
				wantResult: []testutil.TermResults{{
					"C": "d",
				}},
				wantSuccess: true,
			},
			{
				files: map[string][]byte{
					"file": []byte("Hey"),
				},
				program: "get_first_char(C) :- open(file, read, Stream, []), get_char(Stream, C).",
				query:   `get_first_char(C).`,
				wantResult: []testutil.TermResults{{
					"C": "'H'",
				}},
				wantSuccess: true,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(File, write, Stream, _), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(instantiation_error,open/4)"),
				wantSuccess: false,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(34, write, Stream, _), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(type_error(atom,34),open/4)"),
				wantSuccess: false,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(file, write, stream, _), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(instantiation_error,open/4)"),
				wantSuccess: false,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(file, 45, Stream, _), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(type_error(io_mode,45),open/4)"),
				wantSuccess: false,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(file, foo, Stream, _), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(type_error(io_mode,foo),open/4)"),
				wantSuccess: false,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(my_file, write, Stream, _), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(permission_error(input,stream,my_file),open/4)"),
				wantSuccess: false,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(my_file, append, Stream, _), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(permission_error(input,stream,my_file),open/4)"),
				wantSuccess: false,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(file2, read, Stream, _), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(existence_error(source_sink,file2),open/4)"),
				wantSuccess: false,
			},
			{
				files: map[string][]byte{
					"file": []byte("dumb(dumber)."),
				},
				program:     "get_first_char(C) :- open(file, read, Stream, [option1]), get_char(Stream, C).",
				query:       `get_first_char(C).`,
				wantError:   fmt.Errorf("error(domain_error(empty_list,[option1]),open/4)"),
				wantSuccess: false,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a mocked file system", func() {
					uri, _ := url.Parse("file://dump.pl")
					mockedFS := testutil.NewMockFS(ctrl)
					mockedFS.EXPECT().Open(gomock.Any()).AnyTimes().DoAndReturn(func(name string) (fs.VirtualFile, error) {
						for key, bytes := range tc.files {
							if key == name {
								return fs.NewVirtualFile(bytes, uri, time.Now()), nil
							}
						}
						return fs.VirtualFile{}, &fs2.PathError{
							Op:   "read",
							Path: "file2",
							Err:  fmt.Errorf("path not found"),
						}
					})
					Convey("and a context", func() {
						db := dbm.NewMemDB()
						stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
						ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

						Convey("and a vm", func() {
							interpreter := testutil.NewComprehensiveInterpreterMust(ctx)
							interpreter.FS = mockedFS
							interpreter.Register4(engine.NewAtom("open"), Open)

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
