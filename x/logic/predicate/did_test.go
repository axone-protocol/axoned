//nolint:gocognit
package predicate

import (
	"fmt"
	"strings"
	"testing"

	"cosmossdk.io/store/metrics"

	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
)

func TestDID(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program    string
			query      string
			wantResult []types.TermResults
			wantError  error
		}{
			{
				query:      `did_components('did:example:123456',did(X,Y,_,_,_)).`,
				wantResult: []types.TermResults{{"X": "example", "Y": "'123456'"}},
			},
			{
				query:      `did_components('did:example:123456',did(X,Y,Z,_,_)).`,
				wantResult: []types.TermResults{{"X": "example", "Y": "'123456'", "Z": "''"}},
			},
			{
				query:      `did_components('did:example:123456/path', X).`,
				wantResult: []types.TermResults{{"X": "did(example,'123456',path,'','')"}},
			},
			{
				query:      `did_components('did:example:123456?versionId=1', X).`,
				wantResult: []types.TermResults{{"X": "did(example,'123456','','versionId=1','')"}},
			},
			{
				query:      `did_components('did:example:123456/path%20with/space', X).`,
				wantResult: []types.TermResults{{"X": "did(example,'123456','path with/space','','')"}},
			},
			{
				query:      `did_components(X,did(example,'123456',_,'versionId=1',_)).`,
				wantResult: []types.TermResults{{"X": "'did:example:123456?versionId=1'"}},
			},
			{
				query:      `did_components(X,did(example,'123456','/foo/bar','versionId=1','test')).`,
				wantResult: []types.TermResults{{"X": "'did:example:123456/foo/bar?versionId=1#test'"}},
			},
			{
				query:      `did_components(X,did(example,'123456','path with/space',_,test)).`,
				wantResult: []types.TermResults{{"X": "'did:example:123456/path%20with/space#test'"}},
			},
			{
				query:      `did_components(X,Y).`,
				wantResult: []types.TermResults{},
				wantError:  fmt.Errorf("error(instantiation_error,did_components/2)"),
			},
			{
				query:      `did_components('foo',X).`,
				wantResult: []types.TermResults{},
				wantError: fmt.Errorf("error(domain_error(encoding(did),foo),[%s],did_components/2)",
					strings.Join(strings.Split("invalid DID: input length is less than 7", ""), ",")),
			},
			{
				query:      `did_components(123,X).`,
				wantResult: []types.TermResults{},
				wantError:  fmt.Errorf("error(type_error(atom,123),did_components/2)"),
			},
			{
				query:      `did_components(X, 123).`,
				wantResult: []types.TermResults{},
				wantError:  fmt.Errorf("error(type_error(did,123),did_components/2)"),
			},
			{
				query:      `did_components(X,foo('bar')).`,
				wantResult: []types.TermResults{},
				wantError:  fmt.Errorf("error(domain_error(did,foo(bar)),did_components/2)"),
			},
			{
				query:      `did_components(X,did('bar')).`,
				wantResult: []types.TermResults{},
				wantError:  fmt.Errorf("error(domain_error(did,did(bar)),did_components/2)"),
			},
			{
				query:      `did_components(X,did(example,'123456','path with/space',5,test)).`,
				wantResult: []types.TermResults{},
				wantError:  fmt.Errorf("error(type_error(atom,5),did_components/2)"),
			},
			{
				query:      `did_components('did:example:123456',foo(X)).`,
				wantResult: []types.TermResults{},
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
						So(err, ShouldEqual, nil)

						Convey("When the predicate is called", func() {
							sols, err := interpreter.QueryContext(ctx, tc.query)

							Convey("Then the error should be nil", func() {
								So(err, ShouldEqual, nil)
								So(sols, ShouldNotBeNil)

								Convey("and the bindings should be as expected", func() {
									var got []types.TermResults
									for sols.Next() {
										m := types.TermResults{}
										err := sols.Scan(m)
										So(err, ShouldEqual, nil)

										got = append(got, m)
									}
									if tc.wantError != nil {
										So(sols.Err(), ShouldNotEqual, nil)
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
