//nolint:gocognit
package predicate

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func TestCryptoHash(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			program    string
			query      string
			wantResult []types.TermResults
			wantError  error
		}{
			{
				query:      `crypto_hash('foo', Hash).`,
				wantResult: []types.TermResults{{"Hash": "'2c26b46b68ffc68ff99b453c1d30413413422d706483bfa0f98a5e886266e7ae'"}},
			},
			{
				query:      `crypto_hash(Foo, Hash).`,
				wantResult: []types.TermResults{},
				wantError:  fmt.Errorf("crypto_hash/2: invalid data type: engine.Variable, should be Atom"),
			},
			{
				query:      `crypto_hash(foo, bar).`,
				wantResult: []types.TermResults{},
				wantError:  fmt.Errorf("crypto_hash/2: invalid hash type: engine.Atom, should be Variable"),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a context", func() {
					db := tmdb.NewMemDB()
					stateStore := store.NewCommitMultiStore(db)
					ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

					Convey("and a vm", func() {
						interpreter := testutil.NewInterpreterMust(ctx)
						interpreter.Register2(engine.NewAtom("crypto_hash"), CryptoHash)

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
