package predicate

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func TestBlock(t *testing.T) {
	cases := []struct {
		header      tmproto.Header
		implication string
		ok          bool
	}{
		{header: tmproto.Header{Height: 102}, implication: `block_height(102)`, ok: true},
		{header: tmproto.Header{Height: 905}, implication: `block_height(102)`, ok: false},
		{header: tmproto.Header{Height: 102}, implication: `block_height(X), X == 102`, ok: true},
		{header: tmproto.Header{Height: 102}, implication: `block_height(X), X == 905`, ok: false},
		{header: tmproto.Header{Time: time.Unix(1494505756, 0)}, implication: `block_time(1494505756)`, ok: true},
		{header: tmproto.Header{Time: time.Unix(1494505757, 0)}, implication: `block_time(1494505756)`, ok: false},
		{header: tmproto.Header{Time: time.Unix(1494505756, 0)}, implication: `block_time(X), X == 1494505756`, ok: true},
		{header: tmproto.Header{Time: time.Unix(1494505756, 0)}, implication: `block_time(X), X == 1494505757`, ok: false},
	}
	for _, tc := range cases {
		convey.Convey(fmt.Sprintf("Given the clause body: %s", tc.implication), t, func() {
			convey.Convey("Given a context", func() {
				db := tmdb.NewMemDB()
				stateStore := store.NewCommitMultiStore(db)
				ctx := sdk.NewContext(stateStore, tc.header, false, log.NewNopLogger())

				convey.Convey("and a vm", func() {
					vm := testutil.NewVMMust(ctx)
					vm.Register1(engine.NewAtom("block_height"), BlockHeight(ctx))
					vm.Register1(engine.NewAtom("block_time"), BlockTime(ctx))
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
