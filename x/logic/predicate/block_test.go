package predicate

import (
	"fmt"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"
	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/v7/x/logic/testutil"
)

func TestBlock(t *testing.T) {
	cases := []struct {
		header      tmproto.Header
		implication string
		wantOk      bool
	}{
		{header: tmproto.Header{Height: 102}, implication: `block_height(102)`, wantOk: true},
		{header: tmproto.Header{Height: 905}, implication: `block_height(102)`, wantOk: false},
		{header: tmproto.Header{Height: 102}, implication: `block_height(X), X == 102`, wantOk: true},
		{header: tmproto.Header{Height: 102}, implication: `block_height(X), X == 905`, wantOk: false},
		{header: tmproto.Header{Time: time.Unix(1494505756, 0)}, implication: `block_time(1494505756)`, wantOk: true},
		{header: tmproto.Header{Time: time.Unix(1494505757, 0)}, implication: `block_time(1494505756)`, wantOk: false},
		{header: tmproto.Header{Time: time.Unix(1494505756, 0)}, implication: `block_time(X), X == 1494505756`, wantOk: true},
		{header: tmproto.Header{Time: time.Unix(1494505756, 0)}, implication: `block_time(X), X == 1494505757`, wantOk: false},
	}
	for _, tc := range cases {
		Convey(fmt.Sprintf("Given the clause body: %s", tc.implication), t, func() {
			Convey("Given a context", func() {
				db := dbm.NewMemDB()
				stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
				ctx := sdk.NewContext(stateStore, tc.header, false, log.NewNopLogger())

				Convey("and a vm", func() {
					interpreter := testutil.NewLightInterpreterMust(ctx)
					interpreter.Register1(engine.NewAtom("block_height"), BlockHeight)
					interpreter.Register1(engine.NewAtom("block_time"), BlockTime)
					testutil.CompileMust(ctx, interpreter, fmt.Sprintf("test :- %s.", tc.implication))

					Convey("When the predicate is called", func() {
						ok, err := interpreter.Arrive(engine.NewAtom("test"), []engine.Term{}, engine.Success, nil).Force(ctx)

						Convey("Then the result should be true and there should be no error", func() {
							So(err, ShouldBeNil)
							So(ok, ShouldEqual, tc.wantOk)
						})
					})
				})
			})
		})
	}
}
