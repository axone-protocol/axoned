package predicate

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/engine"
	dbm "github.com/cosmos/cosmos-db"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v10/x/logic/testutil"
	"github.com/axone-protocol/axoned/v10/x/logic/util"
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
						db := dbm.NewMemDB()
						stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
						ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())

						Convey("and a vm", func() {
							interpreter := testutil.NewComprehensiveInterpreterMust(ctx)
							interpreter.SetUserOutput(engine.NewOutputTextStream(buffer))

							Convey("When the predicate is called", func() {
								sols, err := interpreter.QueryContext(ctx, tc.query)

								Convey("Then the error should be nil", func() {
									So(err, ShouldBeNil)
									So(sols, ShouldNotBeNil)

									m := testutil.TermResults{}
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
