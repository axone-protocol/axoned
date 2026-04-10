package meter

import (
	"math"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/axone-protocol/axoned/v15/x/logic/testutil"
)

func TestVMMeter(t *testing.T) {
	Convey("Given a VM meter backed by an SDK gas meter", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGasMeter := testutil.NewMockGasMeter(ctrl)
		sut := NewVMMeter(mockGasMeter, 2, 3, 4)

		Convey("it uses the compute coefficient for compute kinds", func() {
			mockGasMeter.EXPECT().ConsumeGas(uint64(10), "Instruction").Times(1)
			sut(engine.MeterInstruction, 5)

			mockGasMeter.EXPECT().ConsumeGas(uint64(12), "ArithNode").Times(1)
			sut(engine.MeterArithNode, 6)

			mockGasMeter.EXPECT().ConsumeGas(uint64(14), "CompareStep").Times(1)
			sut(engine.MeterCompareStep, 7)
		})

		Convey("it uses the memory coefficient for memory kinds", func() {
			mockGasMeter.EXPECT().ConsumeGas(uint64(15), "CopyNode").Times(1)
			sut(engine.MeterCopyNode, 5)

			mockGasMeter.EXPECT().ConsumeGas(uint64(18), "ListCell").Times(1)
			sut(engine.MeterListCell, 6)
		})

		Convey("it uses the unify coefficient for unify kinds", func() {
			mockGasMeter.EXPECT().ConsumeGas(uint64(20), "UnifyStep").Times(1)
			sut(engine.MeterUnifyStep, 5)
		})
	})

	Convey("Given a VM meter with zero-valued coefficients", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGasMeter := testutil.NewMockGasMeter(ctrl)
		sut := NewVMMeter(mockGasMeter, 0, 0, 0)

		Convey("it defaults all coefficients to one", func() {
			mockGasMeter.EXPECT().ConsumeGas(uint64(5), "Instruction").Times(1)
			sut(engine.MeterInstruction, 5)

			mockGasMeter.EXPECT().ConsumeGas(uint64(6), "ListCell").Times(1)
			sut(engine.MeterListCell, 6)

			mockGasMeter.EXPECT().ConsumeGas(uint64(7), "UnifyStep").Times(1)
			sut(engine.MeterUnifyStep, 7)
		})
	})

	Convey("Given a VM meter multiplication overflow", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockGasMeter := testutil.NewMockGasMeter(ctrl)
		sut := NewVMMeter(mockGasMeter, math.MaxUint64, 1, 1)

		Convey("it saturates the gas charge", func() {
			mockGasMeter.EXPECT().ConsumeGas(uint64(math.MaxUint64), "Instruction").Times(1)
			sut(engine.MeterInstruction, 2)
		})
	})

	Convey("Given a VM meter hitting the SDK gas limit", t, func() {
		gasMeter := storetypes.NewGasMeter(10)
		sut := NewVMMeter(gasMeter, 2, 1, 1)

		Convey("it converts the out-of-gas panic into a Prolog resource_error", func() {
			formal := sut(engine.MeterInstruction, 6)

			So(matchesTerm(formal, engine.NewAtom("resource_error").Apply(engine.NewAtom("instruction"))), ShouldBeTrue)
		})
	})

	Convey("Given a VM meter saturation on a non-empty gas meter", t, func() {
		gasMeter := storetypes.NewGasMeter(math.MaxUint64)
		gasMeter.ConsumeGas(1, "seed")
		sut := NewVMMeter(gasMeter, math.MaxUint64, 1, 1)

		Convey("it converts the gas overflow panic into a Prolog resource_error", func() {
			formal := sut(engine.MeterInstruction, 2)

			So(matchesTerm(formal, engine.NewAtom("resource_error").Apply(engine.NewAtom("instruction"))), ShouldBeTrue)
		})
	})
}

func matchesTerm(actual, expected engine.Term) bool {
	_, ok := ((*engine.Env)(nil)).Unify(actual, expected)

	return ok
}
