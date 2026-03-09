package meter

import (
	"fmt"
	"math"

	"github.com/axone-protocol/prolog/v3/engine"

	storetypes "cosmossdk.io/store/types"
)

const defaultCoeff = uint64(1)

// NewVMMeter returns a Prolog VM meter backed by an SDK gas meter.
func NewVMMeter(gasMeter storetypes.GasMeter, computeCoeff, memoryCoeff, unifyCoeff uint64) engine.MeterFunc {
	computeCoeff = nonZeroOrOne(computeCoeff)
	memoryCoeff = nonZeroOrOne(memoryCoeff)
	unifyCoeff = nonZeroOrOne(unifyCoeff)

	return func(kind engine.MeterKind, units uint64) (formal engine.Term) {
		coeff, descriptor, resource := coeffForKind(kind, computeCoeff, memoryCoeff, unifyCoeff)
		consumed, overflow := multiplyUint64Overflow(coeff, units)
		defer func() {
			if r := recover(); r != nil {
				switch r.(type) {
				case storetypes.ErrorOutOfGas, storetypes.ErrorGasOverflow:
					formal = engine.NewAtom("resource_error").Apply(engine.NewAtom(resource))
					return
				}
				panic(r)
			}
		}()

		if overflow {
			gasMeter.ConsumeGas(math.MaxUint64, descriptor)
			return formal
		}

		gasMeter.ConsumeGas(consumed, descriptor)
		return formal
	}
}

func coeffForKind(kind engine.MeterKind, computeCoeff, memoryCoeff, unifyCoeff uint64) (uint64, string, string) {
	switch kind {
	case engine.MeterInstruction:
		return computeCoeff, "Instruction", "instruction"
	case engine.MeterArithNode:
		return computeCoeff, "ArithNode", "arith_node"
	case engine.MeterCompareStep:
		return computeCoeff, "CompareStep", "compare_step"
	case engine.MeterCopyNode:
		return memoryCoeff, "CopyNode", "copy_node"
	case engine.MeterListCell:
		return memoryCoeff, "ListCell", "list_cell"
	case engine.MeterUnifyStep:
		return unifyCoeff, "UnifyStep", "unify_step"
	default:
		panic(fmt.Sprintf("unsupported prolog meter kind: %d", kind))
	}
}

func nonZeroOrOne(v uint64) uint64 {
	if v == 0 {
		return defaultCoeff
	}

	return v
}
