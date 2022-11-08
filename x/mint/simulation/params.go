package simulation

// DONTCOVER

import (
	"fmt"
	"math/rand"

	"github.com/cosmos/cosmos-sdk/x/simulation"

	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/okp4/okp4d/x/mint/types"
)

const (
	KeyAnnualReductionFactor = "AnnualReductionFactor"
)

// ParamChanges defines the parameters that can be modified by param change proposals
// on the simulation.
func ParamChanges(r *rand.Rand) []simtypes.ParamChange {
	return []simtypes.ParamChange{
		simulation.NewSimParamChange(types.ModuleName, KeyAnnualReductionFactor,
			func(r *rand.Rand) string {
				return fmt.Sprintf("\"%s\"", GenAnnualReductionFactorMax(r))
			},
		),
	}
}
