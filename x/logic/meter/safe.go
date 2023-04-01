package meter

import (
	"runtime"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// safeMeterDecorater is a wrapper around sdk.GasMeter that provides go-routine-safe access to the underlying gas meter.
// This is needed because the interpreter is uses multiple go-routines, and the gas meter is shared between multiple
// goroutines.
type safeMeterDecorater struct {
	gasMeter sdk.GasMeter
	mutex    sync.RWMutex
}

// WithSafeMeter returns a new instance of sdk.GasMeter that is go-routine-safe.
func WithSafeMeter(gasMeter sdk.GasMeter) sdk.GasMeter {
	return &safeMeterDecorater{
		gasMeter: gasMeter,
	}
}

// ConsumeGas consumes the given amount of gas from the decorated gas meter.
func (m *safeMeterDecorater) ConsumeGas(amount uint64, descriptor string) {
	m.mutex.Lock()
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(sdk.ErrorOutOfGas); ok {
				// Since predicate is called into a goroutine, when out of gas is thrown, the main caller
				// (grpc: https://github.com/okp4/okp4d/blob/main/x/logic/keeper/grpc_query_ask.go#L25-L36, or querier)
				// cannot recover ErrOutOfGas. To avoid the chain panicking, we need to exit without panic.
				// Goexit runs all deferred calls before terminating the goroutine. Because Goexit
				// is not a panic, any recover calls in those deferred functions will return nil.
				// This is a temporary solution before implementing a context cancellation.
				runtime.Goexit()
			}
			panic(r)
		}
	}()
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
