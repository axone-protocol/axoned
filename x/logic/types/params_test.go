package types_test

import (
	"fmt"
	"testing"

	"github.com/okp4/okp4d/x/logic/types"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/math"
)

func TestValidateParams(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			name      string
			params    types.Params
			expectErr bool
			err       error
		}{
			{
				name:      "validate default params",
				params:    types.DefaultParams(),
				expectErr: false,
				err:       nil,
			},
			{
				name: "validate set params",
				params: types.NewParams(
					types.NewInterpreter(
						types.WithBootstrap("bootstrap"),
						types.WithPredicatesBlacklist([]string{"halt/1"}),
						types.WithPredicatesWhitelist([]string{"source_file/1"}),
						types.WithVirtualFilesBlacklist([]string{"file1"}),
						types.WithVirtualFilesWhitelist([]string{"file2"}),
					),
					types.NewLimits(
						types.WithMaxGas(math.NewUint(1)),
						types.WithMaxSize(math.NewUint(2)),
						types.WithMaxResultCount(math.NewUint(3)),
						types.WithMaxUserOutputSize(math.NewUint(4)),
					),
				),
				expectErr: false,
				err:       nil,
			},
			{
				name: "validate invalid virtual files blacklist params",
				params: types.NewParams(
					types.NewInterpreter(
						types.WithVirtualFilesBlacklist([]string{"https://foo{bar/"}),
					),
					types.NewLimits(),
				),
				expectErr: true,
				err:       fmt.Errorf("invalid virtual file in blacklist: https://foo{bar/"),
			},
			{
				name: "validate invalid virtual files whitelist params",
				params: types.NewParams(
					types.NewInterpreter(
						types.WithVirtualFilesWhitelist([]string{"https://foo{bar/"}),
					),
					types.NewLimits(),
				),
				expectErr: true,
				err:       fmt.Errorf("invalid virtual file in whitelist: https://foo{bar/"),
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d: %v, with params: %v", nc, tc.name, tc.params), func() {
					Convey("when validate params", func() {
						err := tc.params.Validate()

						if tc.expectErr {
							Convey("then params validation expect error", func() {
								So(err, ShouldNotBeNil)
								So(err, ShouldResemble, tc.err)
							})
						} else {
							Convey("then error should be nil", func() {
								So(err, ShouldBeNil)
							})
						}
					})
				})
		}
	})
}
