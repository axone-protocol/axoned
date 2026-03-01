package types_test

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/axone-protocol/axoned/v14/x/logic/types"
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
						WeightingFactor:      2,
						DefaultPredicateCost: 1,
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
	})
}
