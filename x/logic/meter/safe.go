package meter

import (
	"sync"

	storetypes "cosmossdk.io/store/types"
)

// safeMeterDecorater is a wrapper around storetypes.GasMeter that provides go-routine-safe access to the underlying gas meter.
// This is needed because the interpreter uses multiple go-routines, and the gas meter is shared between multiple
// goroutines.
type safeMeterDecorater struct {
	gasMeter storetypes.GasMeter
	mutex    sync.RWMutex
}

// WithSafeMeter returns a new instance of storetypes.GasMeter that is go-routine-safe.
func WithSafeMeter(gasMeter storetypes.GasMeter) storetypes.GasMeter {
	return &safeMeterDecorater{
		gasMeter: gasMeter,
	}
}

// ConsumeGas consumes the given amount of gas from the decorated gas meter.
func (m *safeMeterDecorater) ConsumeGas(amount uint64, descriptor string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.gasMeter.ConsumeGas(amount, descriptor)
}

// GasConsumed returns the amount of gas consumed by the decorated gas meter.
func (m *safeMeterDecorater) GasConsumed() uint64 {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.GasConsumed()
}

// GasConsumedToLimit returns the amount of gas consumed by the decorated gas meter.
func (m *safeMeterDecorater) GasConsumedToLimit() uint64 {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.GasConsumedToLimit()
}

// IsPastLimit returns true if the gas limit has been reached.
func (m *safeMeterDecorater) IsPastLimit() bool {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.IsPastLimit()
}

// Limit returns the gas limit of the decorated gas meter.
func (m *safeMeterDecorater) Limit() uint64 {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.Limit()
}

// RefundGas refunds the given amount of gas to the decorated gas meter.
func (m *safeMeterDecorater) RefundGas(amount uint64, descriptor string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.gasMeter.RefundGas(amount, descriptor)
}

// GasRemaining returns the amount of gas remaining in the decorated gas meter.
func (m *safeMeterDecorater) GasRemaining() uint64 {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.GasRemaining()
}

// String returns a string representation of the decorated gas meter.
func (m *safeMeterDecorater) String() string {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.String()
}

// IsOutOfGas returns true if the gas limit has been reached.
func (m *safeMeterDecorater) IsOutOfGas() bool {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.IsOutOfGas()
}
