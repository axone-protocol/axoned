package meter

import (
	"fmt"
	"math"
	"testing"

	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/axone-protocol/axoned/v9/x/logic/testutil"
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

func TestWeightedMeter(t *testing.T) {
	Convey("Under a mocked environment", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		// for _, _ = range cases {
		Convey("with a context", func() {
			mockGasMeter := testutil.NewMockGasMeter(ctrl)

			Convey("and a weighted gas meter", func() {
				sut := WithWeightedMeter(mockGasMeter, 2)

				Convey("then we should be able to get the gas consumed", func() {
					mockGasMeter.EXPECT().GasConsumed().Return(uint64(100)).Times(1)
					mockGasMeter.EXPECT().GasConsumedToLimit().Return(uint64(200)).Times(1)
					mockGasMeter.EXPECT().GasRemaining().Return(uint64(300)).Times(1)
					mockGasMeter.EXPECT().Limit().Return(uint64(400)).Times(1)
					mockGasMeter.EXPECT().IsPastLimit().Return(false).Times(1)
					mockGasMeter.EXPECT().IsOutOfGas().Return(false).Times(1)

					So(uint64(100), ShouldEqual, sut.GasConsumed())
					So(uint64(200), ShouldEqual, sut.GasConsumedToLimit())
					So(uint64(300), ShouldEqual, sut.GasRemaining())
					So(uint64(400), ShouldEqual, sut.Limit())
					So(sut.IsPastLimit(), ShouldBeFalse)
					So(sut.IsOutOfGas(), ShouldBeFalse)
				})

				Convey("then we should be able to consume gas without overflow", func() {
					mockGasMeter.EXPECT().ConsumeGas(uint64(200), "mock").Times(1)

					sut.ConsumeGas(100, "mock")
				})

				Convey("then we should be able to consume gas with overflow", func() {
					mockGasMeter.EXPECT().ConsumeGas(uint64(math.MaxUint64), "mock").Times(1)

					sut.ConsumeGas(uint64(math.MaxUint64)>>1+1, "mock")
				})
			})
		})
	})
}
