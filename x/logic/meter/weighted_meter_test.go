package meter

import (
	"fmt"
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiplyUint64Overflow(t *testing.T) {
	Convey("Given test cases", t, func() {
		testCases := []struct {
			a            uint64
			b            uint64
			wantResult   uint64
			wantOverflow bool
		}{
			{0, 0, 0, false},
			{0, 1, 0, false},
			{1, 0, 0, false},
			{1, 1, 1, false},
			{1, math.MaxUint64, math.MaxUint64, false},
			{math.MaxUint64, 1, math.MaxUint64, false},
			{2, math.MaxUint64, 0, true},
			{math.MaxUint64, 2, 0, true},
			{math.MaxUint64, math.MaxUint64, 0, true},
		}

		for _, tc := range testCases {
			Convey(fmt.Sprintf("When calling the multiplyUint64Overflow(%d, %d)", tc.a, tc.b), func() {
				actual, overflow := multiplyUint64Overflow(tc.a, tc.b)

				Convey("Then we should get ()", func() {
					So(overflow, ShouldEqual, tc.wantOverflow)
					So(actual, ShouldEqual, tc.wantResult)
				})
			})
		}
	})
}
