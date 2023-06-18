package predicate

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"

	. "github.com/smartystreets/goconvey/convey"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestWrite(t *testing.T) {
	Convey("Given test cases", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cases := []struct {
			userOutputSize int
			query          string
			wantUserOutput string
		}{
			{userOutputSize: 1, query: "put_char('b').", wantUserOutput: "b"},
			{userOutputSize: 5, query: "put_char('a'), put_char('b').", wantUserOutput: "ab"},
			{userOutputSize: 5, query: "write('hello world'), put_char('!').", wantUserOutput: "orld!"},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the query #%d: %s", nc, tc.query), func() {
				Convey("and a mocked output stream", func() {
					buffer := util.NewBoundedBufferMust(tc.userOutputSize)

					Convey("and a context", func() {
						db := tmdb.NewMemDB()
						stateStore := store.NewCommitMultiStore(db)
						ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

						Convey("and a vm", func() {
							interpreter := testutil.NewComprehensiveInterpreterMust(ctx)
							interpreter.SetUserOutput(engine.NewOutputTextStream(buffer))

							Convey("When the predicate is called", func() {
								sols, err := interpreter.QueryContext(ctx, tc.query)

								Convey("Then the error should be nil", func() {
									So(err, ShouldBeNil)
									So(sols, ShouldNotBeNil)

									m := types.TermResults{}
									for sols.Next() {
										err := sols.Scan(m)
										So(err, ShouldBeNil)
									}
									So(sols.Err(), ShouldBeNil)

									Convey(fmt.Sprintf("and the user output should be: %s", tc.wantUserOutput), func() {
										So(buffer.String(), ShouldEqual, tc.wantUserOutput)
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
