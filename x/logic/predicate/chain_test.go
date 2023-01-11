package predicate

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func TestChainID(t *testing.T) {
	cases := []struct {
		header      tmproto.Header
		implication string
		ok          bool
	}{
		{header: tmproto.Header{ChainID: "okp4-nemeton-1"}, implication: `chain_id("okp4-nemeton-1")`, ok: true},
		{header: tmproto.Header{ChainID: "okp4-nemeton-1"}, implication: `chain_id("akashnet-2")`, ok: false},
		{header: tmproto.Header{ChainID: "okp4-nemeton-1"}, implication: `chain_id(X), X == "okp4-nemeton-1"`, ok: true},
		{header: tmproto.Header{ChainID: "okp4-nemeton-1"}, implication: `chain_id(X), X == "akashnet-2"`, ok: false},
	}
	for _, tc := range cases {
		convey.Convey(fmt.Sprintf("Given the clause body: %s", tc.implication), t, func() {
			convey.Convey("Given a context", func() {
				db := tmdb.NewMemDB()
				stateStore := store.NewCommitMultiStore(db)
				ctx := sdk.NewContext(stateStore, tc.header, false, log.NewNopLogger())

				convey.Convey("and a vm", func() {
					vm := testutil.NewVMMust(ctx)
					vm.Register1(engine.NewAtom("chain_id"), ChainID(ctx))
					testutil.CompileMust(ctx, vm, fmt.Sprintf("test :- %s.", tc.implication))

					convey.Convey("When the predicate is called", func() {
						ok, err := vm.Arrive(engine.NewAtom("test"), []engine.Term{}, engine.Success, nil).Force(ctx)

						convey.Convey("Then the result should be true and there should be no error", func() {
							convey.So(err, convey.ShouldBeNil)
							convey.So(ok, convey.ShouldEqual, tc.ok)
						})
					})
				})
			})
		})
	}
}
