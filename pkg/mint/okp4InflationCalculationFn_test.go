package mint

import (
	"fmt"
	"testing"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/okp4/okp4d/x/mint/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOkp4InflationCalculationFn(t *testing.T) {
	Convey("Considering the Okp4InflationCalculationFn", t, func(c C) {
		type args struct {
			blockHeight         uint64
			blocksPerYear       uint64
			inflationRateChange float64
		}
		tests := []struct {
			name string
			args args
			want sdk.Dec
		}{
			{
				name: "Inflation for the first block of the first year",
				args: args{
					blockHeight:         0,
					blocksPerYear:       10,
					inflationRateChange: .8,
				}, want: sdk.NewDecWithPrec(75, 3),
			},
			{
				name: "Inflation for the last block of the first year",
				args: args{
					blockHeight:         9,
					blocksPerYear:       10,
					inflationRateChange: .8,
				}, want: sdk.NewDecWithPrec(75, 3),
			},
			{
				name: "Inflation for the first block of the second year",
				args: args{
					blockHeight:         10,
					blocksPerYear:       10,
					inflationRateChange: .8,
				}, want: sdk.NewDecWithPrec(6, 2),
			},
			{
				name: "Inflation for the second block of the third year",
				args: args{
					blockHeight:         21,
					blocksPerYear:       10,
					inflationRateChange: .8,
				}, want: sdk.MustNewDecFromStr("0.048"),
			},
			{
				name: "Inflation for a block in the 16th year",
				args: args{
					blockHeight:         87899401,
					blocksPerYear:       5256000,
					inflationRateChange: .8,
				}, want: sdk.MustNewDecFromStr("0.002111062325329920"),
			},
		}

		for i, tt := range tests {
			Convey(fmt.Sprintf("Given the context for test case %d (%s)", i, tt.name), func() {
				mkey := sdk.NewKVStoreKey(fmt.Sprintf("test-%d", i))
				tkey := sdk.NewTransientStoreKey(fmt.Sprintf("transient_test_%d", i))
				ctx := testutil.DefaultContext(mkey, tkey).WithBlockHeight(int64(tt.args.blockHeight))
				minter := minttypes.DefaultInitialMinter()
				params := minttypes.DefaultParams()
				params.MintDenom = "uknow"
				params.BlocksPerYear = tt.args.blocksPerYear
				params.InflationRateChange = sdk.MustNewDecFromStr(fmt.Sprintf("%f", tt.args.inflationRateChange))

				Convey("When calling testcase", func() {
					got := Okp4InflationCalculationFn(ctx, minter, params, sdk.ZeroDec())

					Convey(fmt.Sprintf("Then result should be '%d'", tt.want), func() {
						So(got, ShouldEqual, tt.want)
					})
				})
			})
		}
	})
}
