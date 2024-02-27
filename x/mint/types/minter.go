package types

import (
	"fmt"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation, annual
// provisions values.
func NewMinter(inflation, annualProvisions math.LegacyDec) Minter {
	return Minter{
		Inflation:        inflation,
		AnnualProvisions: annualProvisions,
	}
}

// NewMinterWithInitialInflation returns an initial Minter object with a given inflation value and zero annual provisions.
func NewMinterWithInitialInflation(inflation math.LegacyDec) Minter {
	return NewMinter(
		inflation,
		math.LegacyNewDec(0),
	)
}

// NewMinterWithInflationCoef returns a new Minter with updated inflation and annual provisions values.
func NewMinterWithInflationCoef(
	inflationCoef, bondedRatio math.LegacyDec,
	minBound, maxBound *math.LegacyDec,
	totalSupply math.Int,
) (Minter, error) {
	inflationRate, err := inflationRate(inflationCoef, bondedRatio, minBound, maxBound)
	if err != nil {
		return Minter{}, err
	}
	minter := NewMinter(inflationRate, inflationRate.MulInt(totalSupply))

	return minter, minter.Validate()
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 0%.
func DefaultInitialMinter() Minter {
	return NewMinterWithInitialInflation(
		math.LegacyNewDec(0),
	)
}

// Validate validates the mint parameters.
func (m Minter) Validate() error {
	if m.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			m.Inflation.String())
	}
	return nil
}

// inflationRate returns the inflation rate computed from the current bonded ratio
// and the inflation parameter.
func inflationRate(inflationCoef, bondedRatio math.LegacyDec, minBound, maxBound *math.LegacyDec) (math.LegacyDec, error) {
	if bondedRatio.IsZero() {
		return math.LegacyZeroDec(), ErrBondedRatioIsZero
	}

	rate := inflationCoef.Quo(bondedRatio)
	if minBound != nil {
		rate = math.LegacyMaxDec(rate, *minBound)
	}
	if maxBound != nil {
		rate = math.LegacyMinDec(rate, *maxBound)
	}

	return rate, nil
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(math.NewIntFromUint64(params.BlocksPerYear))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
