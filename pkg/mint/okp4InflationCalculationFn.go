package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/okp4/okp4d/x/mint/types"
)

var initialInflation = sdk.NewDecWithPrec(75, 3)

// Okp4InflationCalculationFn is the function used to calculate the inflation for the OKP4 network.
// Inflation is calculated in absolute terms, without taking into account previous inflation, on the basis of the current
// block height, knowing the total number of blocks for one year, using the formula:
//
// Inflation for year X = 0.15 * inflationRateChange^(X-1)
//
// See: https://docs.okp4.network/docs/whitepaper/tokenomics#staking-rewards
func Okp4InflationCalculationFn(ctx sdk.Context, minter minttypes.Minter, params minttypes.Params, _ sdk.Dec) sdk.Dec {
	year := uint64(ctx.BlockHeight()) / params.BlocksPerYear

	return initialInflation.Mul(params.InflationRateChange.Power(year))
}
