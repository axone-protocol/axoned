package types_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

func TestValidateParams(t *testing.T) {
	Convey("Given test cases", t, func() {
		cases := []struct {
			name   string
			params types.Params
		}{
			{
				name:   "default params",
				params: types.DefaultParams(),
			},
			{
				name: "custom params",
				params: types.NewParams(
					types.NewLimits(
						types.WithMaxSize(2),
						types.WithMaxResultCount(3),
						types.WithMaxUserOutputSize(4),
						types.WithMaxVariables(5),
					),
					types.GasPolicy{
						ComputeCoeff: 2,
						MemoryCoeff:  3,
						UnifyCoeff:   4,
						IoCoeff:      5,
					},
				),
			},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given test case #%d: %s", nc, tc.name), func() {
				Convey("when validate params", func() {
					err := tc.params.Validate()

					Convey("then error should be nil", func() {
						So(err, ShouldBeNil)
					})
				})
			})
		}

		// Caps must match x/logic/types/params.go (validateLimits).
		const (
			capBytes        uint64 = 512 << 20
			capResultCount  uint64 = 1_000_000
			capMaxVariables uint64 = 10_000_000
		)

		invalidCases := []struct {
			name   string
			params types.Params
			substr string
		}{
			{
				name: "max_size above cap",
				params: types.NewParams(
					types.NewLimits(types.WithMaxSize(capBytes+1)),
					types.DefaultGasPolicy(),
				),
				substr: "max_size",
			},
			{
				name: "max_result_count above cap",
				params: types.NewParams(
					types.NewLimits(types.WithMaxResultCount(capResultCount+1)),
					types.DefaultGasPolicy(),
				),
				substr: "max_result_count",
			},
			{
				name: "max_user_output_size above cap",
				params: types.NewParams(
					types.NewLimits(types.WithMaxUserOutputSize(capBytes+1)),
					types.DefaultGasPolicy(),
				),
				substr: "max_user_output_size",
			},
			{
				name: "max_variables above cap",
				params: types.NewParams(
					types.NewLimits(types.WithMaxVariables(capMaxVariables+1)),
					types.DefaultGasPolicy(),
				),
				substr: "max_variables",
			},
		}

		for nc, tc := range invalidCases {
			Convey(fmt.Sprintf("Invalid case #%d: %s", nc, tc.name), func() {
				err := tc.params.Validate()
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, tc.substr)
			})
		}
	})
}
