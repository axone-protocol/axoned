//nolint:gocognit
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
)

func TestDID(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program    string
			query      string
			wantResult []testutil.TermResults
			wantError  error
		}{
			{
				query:      `did_components('did:example:123456',did_components(X,Y,_,_,_)).`,
				wantResult: []testutil.TermResults{{"X": "example", "Y": "'123456'"}},
			},
			{
				query:      `did_components('did:example:123456',did_components(X,Y,Z,_,_)).`,
				wantResult: []testutil.TermResults{{"X": "example", "Y": "'123456'", "Z": "_1"}},
			},
			{
				query:      `did_components('did:example:123456/path', X).`,
				wantResult: []testutil.TermResults{{"X": "did_components(example,'123456',path,_1,_2)"}},
			},
			{
				query:      `did_components('did:example:123456?versionId=1', X).`,
				wantResult: []testutil.TermResults{{"X": "did_components(example,'123456',_1,'versionId=1',_2)"}},
			},
			{
				query:      `did_components('did:example:123456/path%20with/space', X).`,
				wantResult: []testutil.TermResults{{"X": "did_components(example,'123456','path%20with/space',_1,_2)"}},
			},
			{
				query:      `did_components(X,did_components(example,'123456',_,'versionId=1',_)).`,
				wantResult: []testutil.TermResults{{"X": "'did:example:123456?versionId=1'"}},
			},
			{
				query:      `did_components(X,did_components(example,'123456','/foo/bar','versionId=1','test')).`,
				wantResult: []testutil.TermResults{{"X": "'did:example:123456/foo/bar?versionId=1#test'"}},
			},
			{
				query:      `did_components(X,did_components(example,'123456','path%20with/space',_,test)).`,
				wantResult: []testutil.TermResults{{"X": "'did:example:123456/path%20with/space#test'"}},
			},
			{
				query:      `did_components(X,did_components(example,'123456','/foo/bar','version%20Id=1','test')).`,
				wantResult: []testutil.TermResults{{"X": "'did:example:123456/foo/bar?version%20Id=1#test'"}},
			},
			{
				query:      `did_components(X,Y).`,
				wantResult: []testutil.TermResults{},
				wantError:  fmt.Errorf("error(instantiation_error,did_components/2)"),
			},
			{
				query:      `did_components('foo',X).`,
				wantResult: []testutil.TermResults{},
				wantError: fmt.Errorf("error(domain_error(encoding(did),foo),[%s],did_components/2)",
					strings.Join(strings.Split("invalid DID", ""), ",")),
			},
			{
				query:      `did_components(123,X).`,
				wantResult: []testutil.TermResults{},
				wantError:  fmt.Errorf("error(type_error(atom,123),did_components/2)"),
			},
			{
				query:      `did_components(X, 123).`,
				wantResult: []testutil.TermResults{},
				wantError:  fmt.Errorf("error(type_error(did_components,123),did_components/2)"),
			},
			{
				query:      `did_components(X,foo('bar')).`,
				wantResult: []testutil.TermResults{},
				wantError:  fmt.Errorf("error(domain_error(did_components,foo(bar)),did_components/2)"),
			},
			{
				query:      `did_components(X,did_components('bar')).`,
				wantResult: []testutil.TermResults{},
				wantError:  fmt.Errorf("error(domain_error(did_components,did_components(bar)),did_components/2)"),
			},
			{
				query:      `did_components(X,did_components(example,'123456','path with/space',5,test)).`,
				wantResult: []testutil.TermResults{},
				wantError:  fmt.Errorf("error(type_error(atom,5),did_components/2)"),
			},
			{
				query:      `did_components('did:example:123456',foo(X)).`,
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
						interpreter.Register2(engine.NewAtom("did_components"), DIDComponents)

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
										So(len(got), ShouldEqual, len(tc.wantResult))
										for iGot, resultGot := range got {
											for varGot, termGot := range resultGot {
												So(testutil.ReindexUnknownVariables(termGot), ShouldEqual, tc.wantResult[iGot][varGot])
											}
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
