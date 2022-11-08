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
	Inflation             = "inflation"
	AnnualReductionFactor = "annual_reduction_factor"
)

// GenInflation randomized Inflation.
func GenInflation(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(int64(r.Intn(99)), 2)
}

// GenAnnualReductionFactor randomized AnnualReductionFactor.
func GenAnnualReductionFactorMax(r *rand.Rand) sdk.Dec {
	return sdk.NewDecWithPrec(20, 2)
}

// RandomizedGenState generates a random GenesisState for mint.
func RandomizedGenState(simState *module.SimulationState) {
	// minter
	var inflation sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, Inflation, &inflation, simState.Rand,
		func(r *rand.Rand) { inflation = GenInflation(r) },
	)

	// params

	var annualReductionFactor sdk.Dec
	simState.AppParams.GetOrGenerate(
		simState.Cdc, AnnualReductionFactor, &annualReductionFactor, simState.Rand,
		func(r *rand.Rand) { annualReductionFactor = GenAnnualReductionFactorMax(r) },
	)

	mintDenom := sdk.DefaultBondDenom
	blocksPerYear := uint64(60 * 60 * 8766 / 5)
	params := types.NewParams(mintDenom, annualReductionFactor, blocksPerYear)
	annualProvision := inflation.MulInt(simState.InitialStake)
	targetSupply := simState.InitialStake.Add(annualProvision.TruncateInt())

	minter := types.InitialMinter(inflation, targetSupply)
	minter.AnnualProvisions = annualProvision

	mintGenesis := types.NewGenesisState(minter, params)

	bz, err := json.MarshalIndent(&mintGenesis, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated minting parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(mintGenesis)
}
