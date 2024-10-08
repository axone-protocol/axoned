package types

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/math"
)

func toLegacyDec(i int64) *math.LegacyDec {
	ii := math.LegacyNewDec(i)
	return &ii
}

func Test_validateBounds(t *testing.T) {
	Convey("Given test cases", t, func() {
		cases := []struct {
			name               string
			minBound, maxBound interface{}
			expectErr          bool
			err                error
		}{
			{
				name:      "validate bounds types",
				minBound:  toLegacyDec(2),
				maxBound:  toLegacyDec(99999999),
				expectErr: false,
				err:       nil,
			},
			{
				name:      "validate invalid min type",
				minBound:  2,
				maxBound:  toLegacyDec(99999999),
				expectErr: true,
				err: func() error {
					return fmt.Errorf("invalid parameter type: %T", 2)
				}(),
			},
			{
				name:      "validate invalid max type",
				minBound:  toLegacyDec(0),
				maxBound:  999999.0,
				expectErr: true,
				err: func() error {
					return fmt.Errorf("invalid parameter type: %T", 999999.0)
				}(),
			},
			{
				name:      "validate non-negative bounds",
				minBound:  toLegacyDec(-2),
				maxBound:  toLegacyDec(99999999),
				expectErr: true,
				err: func() error {
					return fmt.Errorf("inflation bound cannot be negative: %s", toLegacyDec(-2))
				}(),
			},
			{
				name:      "validate non-negative bounds",
				minBound:  toLegacyDec(99999999),
				maxBound:  toLegacyDec(1),
				expectErr: true,
				err: func() error {
					return fmt.Errorf("inflation min cannot be greater than inflation max")
				}(),
			},
		}
		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d: %v, with params: %v, %v", nc, tc.name, tc.minBound, tc.maxBound), func() {
					Convey("when validate bounds", func() {
						err := validateBounds(tc.minBound, tc.maxBound)

						if tc.expectErr {
							Convey("then bounds validation expect error", func() {
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
