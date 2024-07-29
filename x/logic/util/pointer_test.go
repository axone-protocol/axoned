package util

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	sdkmath "cosmossdk.io/math"
)

func TestDerefOrDefault(t *testing.T) {
	Convey("Given a pointer to an int and a default int value", t, func() {
		x := 5
		ptr := &x
		defaultValue := 10

		Convey("When the pointer is not nil", func() {
			result := DerefOrDefault(ptr, defaultValue)

			Convey("The result should be the value pointed to by the pointer", func() {
				So(result, ShouldEqual, x)
			})
		})

		Convey("When the pointer is nil", func() {
			result := DerefOrDefault(nil, defaultValue)

			Convey("The result should be the default value", func() {
				So(result, ShouldEqual, defaultValue)
			})
		})
	})
}

func TestNonZeroOrDefault(t *testing.T) {
	Convey("Given a value", t, func() {
		cases := []struct {
			v            any
			defaultValue any
			expected     any
		}{
			{nil, 0, 0},
			{0, 10, 10},
			{1, 0, 1},
			{"", "default", "default"},
			{"hello", "default", "hello"},
		}
		for _, tc := range cases {
			Convey(fmt.Sprintf("When the value is %v", tc.v), func() {
				Convey(fmt.Sprintf("Then the default value %v is returned", tc.defaultValue), func() {
					So(NonZeroOrDefault(tc.v, tc.defaultValue), ShouldEqual, tc.expected)
				})
			})
		}
	})
}

func TestNonZeroOrDefaultUInt(t *testing.T) {
	Convey("Given a value", t, func() {
		cases := []struct {
			v            *sdkmath.Uint
			defaultValue sdkmath.Uint
			expected     sdkmath.Uint
		}{
			{nil, sdkmath.ZeroUint(), sdkmath.ZeroUint()},
			{
				func() *sdkmath.Uint { u := sdkmath.ZeroUint(); return &u }(),
				sdkmath.NewUint(10),
				sdkmath.NewUint(10),
			},
			{
				func() *sdkmath.Uint { u := sdkmath.NewUint(1); return &u }(),
				sdkmath.ZeroUint(),
				sdkmath.NewUint(1),
			},
		}
		for _, tc := range cases {
			Convey(fmt.Sprintf("When the value is %v", tc.v), func() {
				Convey(fmt.Sprintf("Then the default value %v is returned", tc.defaultValue), func() {
					So(NonZeroOrDefaultUInt(tc.v, tc.defaultValue), ShouldEqual, tc.expected)
				})
			})
		}
	})
}
