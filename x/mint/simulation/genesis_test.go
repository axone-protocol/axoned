//nolint:gosec
package simulation_test

import (
	"encoding/json"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/math"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"

	"github.com/axone-protocol/axoned/v14/x/mint/simulation"
	"github.com/axone-protocol/axoned/v14/x/mint/types"
)

// TestRandomizedGenState tests the normal scenario of applying RandomizedGenState.
// Abnormal scenarii are not tested here.
func TestRandomizedGenState(t *testing.T) {
	Convey("Given a fully initialized simulation state", t, func() {
		interfaceRegistry := codectypes.NewInterfaceRegistry()
		cdc := codec.NewProtoCodec(interfaceRegistry)

		s := rand.NewSource(1)
		r := rand.New(s)

		simState := module.SimulationState{
			AppParams:    make(simtypes.AppParams),
			Cdc:          cdc,
			Rand:         r,
			NumBonded:    3,
			Accounts:     simtypes.RandomAccounts(r, 3),
			InitialStake: math.NewInt(1000),
			GenState:     make(map[string]json.RawMessage),
		}

		simulation.RandomizedGenState(&simState)

		var mintGenesis types.GenesisState
		simState.Cdc.MustUnmarshalJSON(simState.GenState[types.ModuleName], &mintGenesis)

		inflationCoef := math.LegacyNewDecWithPrec(3, 2)
		bondedRatio := math.LegacyNewDecWithPrec(2, 1)
		minter, _ := types.NewMinterWithInflationCoef(inflationCoef, bondedRatio, nil, nil, simState.InitialStake)

		So(mintGenesis.Params.BlocksPerYear, ShouldEqual, uint64(6311520))
		So(mintGenesis.Params.InflationCoef.String(), ShouldEqual, "0.073000000000000000")
		So(mintGenesis.Params.MintDenom, ShouldEqual, "stake")
		So(mintGenesis.Minter.BlockProvision(mintGenesis.Params).String(), ShouldEqual, "0stake")
		So(mintGenesis.Minter.Inflation.String(), ShouldEqual, "0.170000000000000000")
		So(mintGenesis.Minter.AnnualProvisions.String(), ShouldEqual, "170.000000000000000000")
		So(minter.Inflation.String(), ShouldEqual, "0.150000000000000000")
		So(minter.AnnualProvisions.String(), ShouldEqual, "150.000000000000000000")
	})
}

// TestRandomizedGenState1 tests abnormal scenarios of applying RandomizedGenState.
func TestRandomizedGenState1(t *testing.T) {
	Convey("Given incompletely initialized simulation states", t, func() {
		interfaceRegistry := codectypes.NewInterfaceRegistry()
		cdc := codec.NewProtoCodec(interfaceRegistry)

		s := rand.NewSource(1)
		r := rand.New(s)

		tests := []struct {
			name     string
			simState module.SimulationState
		}{
			{
				name:     "empty simulation state panics",
				simState: module.SimulationState{},
			},
			{
				name: "missing GenState panics",
				simState: module.SimulationState{
					AppParams: make(simtypes.AppParams),
					Cdc:       cdc,
					Rand:      r,
				},
			},
		}

		for _, tc := range tests {
			Convey(tc.name, func() {
				So(func() { simulation.RandomizedGenState(&tc.simState) }, ShouldPanic)
			})
		}
	})
}
