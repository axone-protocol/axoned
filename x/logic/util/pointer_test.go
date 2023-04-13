package util

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDerefOrDefault(t *testing.T) {
	Convey("Given a pointer to an int and a default int values", t, func() {
		x := 5
		ptr := &x
		defaultValue := 10

		Convey("When the pointer is not nil", func() {
			result := DerefOrDefault(ptr, defaultValue)

			Convey("The result should be the values pointed to by the pointer", func() {
				So(result, ShouldEqual, x)
			})
		})

		Convey("When the pointer is nil", func() {
			result := DerefOrDefault(nil, defaultValue)

			Convey("The result should be the default values", func() {
				So(result, ShouldEqual, defaultValue)
			})
		})
	})
}

func TestNonZeroOrDefault(t *testing.T) {
	Convey("Given a values", t, func() {
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
			Convey(fmt.Sprintf("When the values is %v", tc.v), func() {
				Convey(fmt.Sprintf("Then the default values %v is returned", tc.defaultValue), func() {
					So(NonZeroOrDefault(tc.v, tc.defaultValue), ShouldEqual, tc.expected)
				})
			})
		}
	})
}
