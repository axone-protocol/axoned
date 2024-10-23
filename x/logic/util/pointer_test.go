package util

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
