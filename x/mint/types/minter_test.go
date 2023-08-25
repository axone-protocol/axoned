//nolint:gosec
package types

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNextInflation(t *testing.T) {
	minter := DefaultInitialMinter()
	params := DefaultParams()

	tests := []struct {
		boundedRatio, expInflation sdk.Dec
	}{
		// With a bounded ratio of 66 %, next inflation should be 10.95%
		{sdk.NewDecWithPrec(66, 2), sdk.NewDecWithPrec(1095, 4)},
		// With a bounded ratio of 0 %, next inflation should be 18.25%
		{sdk.NewDecWithPrec(0, 2), sdk.NewDecWithPrec(1825, 4)},
		// With a bounded ratio of 100 %, next inflation should be 7.18%
		{sdk.NewDecWithPrec(1, 0), sdk.NewDecWithPrec(71893939393939394, 18)},
	}
	for i, tc := range tests {
		inflation := minter.NextInflation(params, tc.boundedRatio)

		require.True(t, inflation.Equal(tc.expInflation),
			"Test Index: %v\nInflation:  %v\nExpected: %v\n", i, inflation, tc.expInflation)
	}
}

//nolint:lll
func TestBlockProvision(t *testing.T) {
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()

	secondsPerYear := int64(60 * 60 * 8766)
	blockInterval := int64(5) // there is 1 block each 5 second approximately

	tests := []struct {
		annualProvisions sdk.Dec
		expProvisions    int64
	}{
		{sdk.NewDec(secondsPerYear / blockInterval), 1},
		{sdk.NewDec(secondsPerYear/blockInterval + 1), 1},
		{sdk.NewDec((secondsPerYear / blockInterval) * 2), 2},
		{sdk.NewDec((secondsPerYear / blockInterval) / 2), 0},
		{sdk.NewDec((secondsPerYear / blockInterval) * 20), 20},
	}
	for i, tc := range tests {
		minter.AnnualProvisions = tc.annualProvisions
		provisions := minter.BlockProvision(params)

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
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()

	s1 := rand.NewSource(100)
	r1 := rand.New(s1)
	minter.AnnualProvisions = sdk.NewDec(r1.Int63n(1000000))

	// run the BlockProvision function b.N times
	for n := 0; n < b.N; n++ {
		minter.BlockProvision(params)
	}
}

// Next inflation benchmarking
// BenchmarkNextInflation-4 1000000 1828 ns/op.
func BenchmarkNextInflation(b *testing.B) {
	b.ReportAllocs()
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()

	// run the NextInflationRate function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextInflation(params, sdk.NewDecWithPrec(66, 2))
	}
}

// // Next annual provisions benchmarking
// // BenchmarkNextAnnualProvisions-4 5000000 251 ns/op.
func BenchmarkNextAnnualProvisions(b *testing.B) {
	b.ReportAllocs()
	minter := InitialMinter(sdk.NewDecWithPrec(1, 1))
	params := DefaultParams()
	totalSupply := sdk.NewInt(100000000000000)

	// run the NextAnnualProvisions function b.N times
	for n := 0; n < b.N; n++ {
		minter.NextAnnualProvisions(params, totalSupply)
	}
}
