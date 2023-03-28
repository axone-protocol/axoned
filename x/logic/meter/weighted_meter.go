package meter

import (
	"math"

	"github.com/cosmos/cosmos-sdk/types"
)

// WeightedMeterDecorator is decorator that wraps a gas meter and adds a weight multiplier to the consumed gas.
type WeightedMeterDecorator struct {
	decorated types.GasMeter
	weight    uint64
}

// WithWeightedMeter returns a new WeightedMeterDecorator with the given weight.
func WithWeightedMeter(decorated types.GasMeter, weight uint64) types.GasMeter {
	return &WeightedMeterDecorator{
		decorated: decorated,
		weight:    weight,
	}
}

// GasConsumed returns the amount of gas consumed by the decorated gas meter.
func (m *WeightedMeterDecorator) GasConsumed() types.Gas {
	return m.decorated.GasConsumed()
}

// GasConsumedToLimit returns the amount of gas consumed by the decorated gas meter.
func (m *WeightedMeterDecorator) GasConsumedToLimit() types.Gas {
	return m.decorated.GasConsumedToLimit()
}

// GasRemaining returns the amount of gas remaining in the decorated gas meter.
func (m *WeightedMeterDecorator) GasRemaining() types.Gas {
	return m.decorated.GasRemaining()
}

// Limit returns the limit of the decorated gas meter.
func (m *WeightedMeterDecorator) Limit() types.Gas {
	return m.decorated.Limit()
}

// ConsumeGas consumes the given amount of gas from the decorated gas meter.
func (m *WeightedMeterDecorator) ConsumeGas(amount types.Gas, descriptor string) {
	consumed, overflow := multiplyUint64Overflow(m.weight, amount)
	if overflow {
		m.decorated.ConsumeGas(math.MaxUint64, descriptor)
	} else {
		m.decorated.ConsumeGas(consumed, descriptor)
	}
}

// RefundGas refunds the given amount of gas to the decorated gas meter.
func (m *WeightedMeterDecorator) RefundGas(amount types.Gas, descriptor string) {
	consumed, overflow := multiplyUint64Overflow(m.decorated.GasConsumed(), amount)
	if overflow {
		m.decorated.RefundGas(math.MaxUint64, descriptor)
	} else {
		m.decorated.RefundGas(consumed, descriptor)
	}
}

// IsPastLimit returns true if the decorated gas meter is past the limit.
func (m *WeightedMeterDecorator) IsPastLimit() bool {
	return m.decorated.IsPastLimit()
}

// IsOutOfGas returns true if the decorated gas meter is out of gas.
func (m *WeightedMeterDecorator) IsOutOfGas() bool {
	return m.decorated.IsOutOfGas()
}

// String returns the decorated gas meter's string representation.
func (m *WeightedMeterDecorator) String() string {
	return m.decorated.String()
}

// multiplyUint64Overflow returns the product of a and b and a boolean indicating whether the product overflows.
func multiplyUint64Overflow(a, b uint64) (uint64, bool) {
	if a == 0 || b == 0 {
		return 0, false
	}

	c := a * b
	if c/a != b || c/b != a {
		return 0, true
	}

	return c, false
}
