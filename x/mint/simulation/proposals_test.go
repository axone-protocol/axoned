package simulation_test

import (
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/axone-protocol/axoned/v14/x/mint/simulation"
	"github.com/axone-protocol/axoned/v14/x/mint/types"
)

func TestProposalMsgs(t *testing.T) {
	Convey("Given a deterministic simulation context", t, func() {
		s := rand.NewSource(1)
		r := rand.New(s) //nolint:gosec

		ctx := sdk.NewContext(nil, tmproto.Header{}, true, nil)
		accounts := simtypes.RandomAccounts(r, 3)

		Convey("When proposal messages are generated", func() {
			weightedProposalMsgs := simulation.ProposalMsgs()

			Convey("Then exactly one weighted proposal is returned", func() {
				So(len(weightedProposalMsgs), ShouldEqual, 1)
			})

			Convey("And the weighted proposal exposes expected metadata", func() {
				So(len(weightedProposalMsgs), ShouldEqual, 1)
				if len(weightedProposalMsgs) != 1 {
					return
				}

				w0 := weightedProposalMsgs[0]
				So(w0.AppParamsKey(), ShouldEqual, simulation.OpWeightMsgUpdateParams)
				So(w0.DefaultWeight(), ShouldEqual, simulation.DefaultWeightMsgUpdateParams)
			})

			Convey("When the message simulator function is executed", func() {
				So(len(weightedProposalMsgs), ShouldEqual, 1)
				if len(weightedProposalMsgs) != 1 {
					return
				}

				w0 := weightedProposalMsgs[0]
				msg := w0.MsgSimulatorFn()(r, ctx, accounts)
				msgUpdateParams, ok := msg.(*types.MsgUpdateParams)

				Convey("Then it returns a MsgUpdateParams with expected fields", func() {
					So(ok, ShouldBeTrue)
					So(msgUpdateParams.Authority, ShouldEqual, sdk.AccAddress(address.Module("gov")).String())
					So(msgUpdateParams.Params.BlocksPerYear, ShouldEqual, uint64(122877))
					So(msgUpdateParams.Params.InflationCoef, ShouldResemble, math.LegacyNewDecWithPrec(95, 2))
					So(msgUpdateParams.Params.MintDenom, ShouldEqual, "eAerqyNEUz")
				})
			})
		})
	})
}
