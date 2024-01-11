//nolint:gosec
package types

import (
	"fmt"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestNextInflation(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			name                     string
			inflationRatio           math.LegacyDec
			bondedRatio              math.LegacyDec
			totalSupply              math.Int
			expectedInflation        math.LegacyDec
			expectedAnnualProvisions math.LegacyDec
			expectedErr              error
		}{
			{
				name:                     "inflation ratio is 0",
				inflationRatio:           sdk.NewDec(0),
				bondedRatio:              math.LegacyNewDecWithPrec(20, 2),
				totalSupply:              math.NewInt(1000),
				expectedInflation:        sdk.NewDec(0),
				expectedAnnualProvisions: sdk.NewDec(0),
			},
			{
				name:                     "inflation ratio is 0.03",
				inflationRatio:           math.LegacyNewDecWithPrec(3, 2),
				bondedRatio:              math.LegacyNewDecWithPrec(2, 1),
				totalSupply:              math.NewInt(1000),
				expectedInflation:        math.LegacyNewDecWithPrec(15, 2),
				expectedAnnualProvisions: sdk.NewDec(150),
			},
			{
				name:           "bonded ratio is 0",
				inflationRatio: math.LegacyNewDecWithPrec(3, 2),
				bondedRatio:    sdk.NewDec(0),
				totalSupply:    math.NewInt(1000),
				expectedErr:    fmt.Errorf("bonded ratio is zero"),
			},
			{
				name:           "negative inflation ratio",
				inflationRatio: math.LegacyNewDecWithPrec(3, 2),
				bondedRatio:    math.LegacyNewDecWithPrec(-2, 1),
				totalSupply:    math.NewInt(1000),
				expectedErr:    fmt.Errorf("mint parameter Inflation should be positive, is -0.150000000000000000"),
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d: %v", nc, tc.name), func() {
					Convey("when calling NewMinterWithInflationCoef function", func() {
						minter, err := NewMinterWithInflationCoef(tc.inflationRatio, tc.bondedRatio, tc.totalSupply)
						if tc.expectedErr != nil {
							Convey("then an error should occur", func() {
								So(err, ShouldNotBeNil)
								So(err.Error(), ShouldEqual, tc.expectedErr.Error())
							})
						} else {
							Convey("then minter values should be as expected", func() {
								So(err, ShouldBeNil)
								So(minter.Inflation.String(), ShouldEqual, tc.expectedInflation.String())
								So(minter.AnnualProvisions.String(), ShouldEqual, tc.expectedAnnualProvisions.String())
							})
						}
					})
				})
		}
	})
}

// Benchmarking :)
// previously using math.Int operations:
// BenchmarkBlockProvision-4 5000000 220 ns/op
//
// using math.LegacyDec operations: (current implementation)
// BenchmarkBlockProvision-4 3000000 429 ns/op.
func BenchmarkBlockProvision(b *testing.B) {
	b.ReportAllocs()
	minter := NewMinterWithInitialInflation(math.LegacyNewDecWithPrec(1, 1))
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

	params := DefaultParams()
	bondedRatio := math.LegacyNewDecWithPrec(66, 2)
	totalSupply := sdk.NewInt(100000000000000)

	// run the NextInflationRate function b.N times
	for n := 0; n < b.N; n++ {
		_, err := NewMinterWithInflationCoef(params.InflationCoef, bondedRatio, totalSupply)
		if err != nil {
			panic(err)
		}
	}
}
