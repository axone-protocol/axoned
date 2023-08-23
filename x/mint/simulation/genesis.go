package simulation

// DONTCOVER

import (
	"encoding/json"
	"fmt"
	"math/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/okp4/okp4d/x/mint/types"
)

// Simulation parameter constants.
const (
	Inflation           = "inflation"
	InflationCoef       = "inflation_coef"
	BoundingAdjustment  = "bounding_adjustment"
	TargetBoundingRatio = "target_bounding_ratio"
)

// GenInflation randomized Inflation.
func GenInflation(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenInflationCoefMax randomized AnnualReductionFactor.
func GenInflationCoefMax(_ *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(73, 3)
}

// GenBoundingAdjustmentMax randomized AnnualReductionFactor.
func GenBoundingAdjustmentMax(_ *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(25, 1)
}

// GenTargetBoundingRatioMax randomized AnnualReductionFactor.
func GenTargetBoundingRatioMax(_ *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(66, 2)
}

// RandomizedGenState generates a random GenesisState for mint.
//
//nolint:forbidigo
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	// params

	var inflationCoef sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, InflationCoef, &inflationCoef, simState.Rand,
		func(r *rand.Rand) { inflationCoef = GenInflationCoefMax(r) },
	)
	var targetBoundingRatio sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, TargetBoundingRatio, &targetBoundingRatio, simState.Rand,
		func(r *rand.Rand) { targetBoundingRatio = GenTargetBoundingRatioMax(r) },
	)

	var boundingAdjustment sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, BoundingAdjustment, &boundingAdjustment, simState.Rand,
		func(r *rand.Rand) { boundingAdjustment = GenBoundingAdjustmentMax(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)
	params := types.NewParams(mintDenom, inflationCoef, boundingAdjustment, targetBoundingRatio, blocksPerYear)
	annualProvision := inflation.MulInt(simState.InitialStake)

	minter := types.InitialMinter(inflation)
	minter.AnnualProvisions = annualProvision

	mintGenesis := types.NewGenesisState(minter, params)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
