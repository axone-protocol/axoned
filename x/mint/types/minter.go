package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation, annual
// provisions values and annual reduction factor.
func NewMinter(inflation, annualProvisions sdk.Dec, targetSupply math.Int) Minter {
	return Minter{
		Inflation:        inflation,
		AnnualProvisions: annualProvisions,
		TargetSupply:     targetSupply,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value and annual reduction factor.
func InitialMinter(inflation sdk.Dec, targetSupply math.Int) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
		targetSupply,
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 15%.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		sdk.NewDecWithPrec(15, 2),
		math.NewInt(230000000000000),
	)
}

// validate minter.
func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			minter.Inflation.String())
	}
	return nil
}

// NextInflation return the new inflation rate for the next year
// Get the current inflation and multiply by (1 - annual reduction factor).
func (m Minter) NextInflation(params Params) sdk.Dec {
	return m.Inflation.Mul(sdk.OneDec().Sub(params.AnnualReductionFactor))
}

// NextAnnualProvisions returns the annual provisions based on current total
// supply and inflation rate.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply math.Int) sdk.Dec {
	return m.Inflation.MulInt(totalSupply)
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params, totalSupply math.Int) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerYear)))

	// Fixe rounding by limiting to the target supply at the end of the year block.
	futureSupply := totalSupply.Add(provisionAmt.TruncateInt())
	if futureSupply.GT(m.TargetSupply) {
		// In case of a rounding is not precise enough, truncating int of provisionAmt could return Zero
		// To avoid negative coin if provisionAmt is equal to Zero, return minimum Zero or more coin.
		return sdk.NewCoin(params.MintDenom, sdk.MaxInt(m.TargetSupply.Sub(totalSupply), sdk.ZeroInt()))
	}

	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
