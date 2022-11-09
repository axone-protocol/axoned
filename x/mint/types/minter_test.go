// nolint: gosec
package types

import (
	"math/rand"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNextInflation(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()

	tests := []struct {
		setInflation, expInflation sdk.Dec
	}{
		// With annual reduction factor of 20 % (defined in params), next infaltion should be 12%
		{sdk.NewDecWithPrec(15, 2), sdk.NewDecWithPrec(12, 2)},
		// With annual reduction factor of 20 % (defined in params), next infaltion should be 12%
		{sdk.NewDecWithPrec(12, 2), sdk.NewDecWithPrec(96, 3)},
	}
	for i, tc := range tests {
		minter.Inflation = tc.setInflation

		inflation := minter.NextInflation(params)

		require.True(t, inflation.Equal(tc.expInflation),
			"Test Index: %v\nInflation:  %v\nExpected: %v\n", i, inflation, tc.expInflation)
	}
}

func TestBlockProvision(t *testing.T) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1), math.NewInt(1))
	params := DefaultParams()

	secondsPerYear := int64(60 * 60 * 8766)
	blockInterval := int64(5) // there is 1 block each 5 second approximately

	tests := []struct {
		annualProvisions int64
		expProvisions    int64
		totalSupply      math.Int
	}{
		{secondsPerYear / blockInterval, 1, sdk.NewInt(1)},
		{secondsPerYear/blockInterval + 1, 1, math.NewInt(1)},
		{(secondsPerYear / 5) * 2, 2, math.NewInt(1)},
		{(secondsPerYear / 5) / 2, 0, math.NewInt(1)},
	}
	for i, tc := range tests {
		minter.AnnualProvisions = sdk.NewDec(tc.annualProvisions)
		minter.TargetSupply = tc.totalSupply.Add(minter.AnnualProvisions.TruncateInt())

		provisions := minter.BlockProvision(params, tc.totalSupply)

		expProvisions := sdk.NewCoin(params.MintDenom,
			sdk.NewInt(tc.expProvisions))

		require.True(t, expProvisions.IsEqual(provisions),
			"test: %v\n\tExp: %v\n\tGot: %v\n",
			i, tc.expProvisions, provisions)
	}
}

// Benchmarking :)
// previously using math.Int operations:
// BenchmarkBlockProvision-4 5000000 220 ns/op
//
// using sdk.Dec operations: (current implementation)
// BenchmarkBlockProvision-4 3000000 429 ns/op.
func BenchmarkBlockProvision(b *testing.B) {
	b.ReportAllocs()
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1), math.NewInt(1))
	params := DefaultParams()

	s1 := rand.NewSource(100)
	r1 := rand.New(s1)
	minter.AnnualProvisions = sdk.NewDec(r1.Int63n(1000000))

	// run the BlockProvision function b.N times
	for n := 0; n < b.N; n++ {
		minter.BlockProvision(params, math.NewInt(1))
	}
}

// Next inflation benchmarking
// BenchmarkNextInflation-4 1000000 1828 ns/op.
func BenchmarkNextInflation(b *testing.B) {
	b.ReportAllocs()
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1), math.NewInt(1))
	params := DefaultParams()

	// run the NextInflationRate function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextInflation(params)
	}
}

// // Next annual provisions benchmarking
// // BenchmarkNextAnnualProvisions-4 5000000 251 ns/op.
func BenchmarkNextAnnualProvisions(b *testing.B) {
	b.ReportAllocs()
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1), math.NewInt(1))
	params := DefaultParams()
	totalSupply := sdk.NewInt(100000000000000)

	// run the NextAnnualProvisions function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextAnnualProvisions(params, totalSupply)
	}
}
