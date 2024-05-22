package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/axone-protocol/axoned/v8/x/mint/types"
)

// Simulation parameter constants.
const (
	Inflation     = "inflation"
	InflationCoef = "inflation_coef"
)

// GenInflation randomized Inflation.
func GenInflation(r *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationCoefMax randomized AnnualReductionFactor.
func GenInflationCoefMax(_ *rand.Rand) math.LegacyDec {
	return math.LegacyNewDecWithPrec(73, 3)
}

// RandomizedGenState generates a random GenesisState for mint.
//
//nolint:forbidigo
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation math.LegacyDec
	simState.AppParams.GetOrGenerate(
		Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	// params

	var inflationCoef math.LegacyDec
	simState.AppParams.GetOrGenerate(
		InflationCoef, &inflationCoef, simState.Rand,
		func(r *rand.Rand) { inflationCoef = GenInflationCoefMax(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)
	params := types.NewParams(mintDenom, inflationCoef, blocksPerYear)
	annualProvision := inflation.MulInt(simState.InitialStake)

	minter := types.NewMinterWithInitialInflation(inflation)
	minter.AnnualProvisions = annualProvision

	mintGenesis := types.NewGenesisState(minter, params)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
