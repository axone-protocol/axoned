package predicate

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v2/engine"
	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v12/x/logic/testutil"
)

func TestChainID(t *testing.T) {
	cases := []struct {
		header      tmproto.Header
		implication string
		wantOk      bool
	}{
		{header: tmproto.Header{ChainID: "axone-nemeton-1"}, implication: `chain_id('axone-nemeton-1')`, wantOk: true},
		{header: tmproto.Header{ChainID: "axone-nemeton-1"}, implication: `chain_id('akashnet-2')`, wantOk: false},
		{header: tmproto.Header{ChainID: "axone-nemeton-1"}, implication: `chain_id(X), X == 'axone-nemeton-1'`, wantOk: true},
		{header: tmproto.Header{ChainID: "axone-nemeton-1"}, implication: `chain_id(X), X == "axone-nemeton-1"`, wantOk: false},
		{header: tmproto.Header{ChainID: "axone-nemeton-1"}, implication: `chain_id(X), X == 'akashnet-2'`, wantOk: false},
	}
	for _, tc := range cases {
		Convey(fmt.Sprintf("Given the clause body: %s", tc.implication), t, func() {
			Convey("Given a context", func() {
				db := dbm.NewMemDB()
				stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
				ctx := sdk.NewContext(stateStore, tc.header, false, log.NewNopLogger())

				Convey("and an interpreter", func() {
					interpreter := testutil.NewLightInterpreterMust(ctx)
					interpreter.Register1(engine.NewAtom("chain_id"), ChainID)
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
