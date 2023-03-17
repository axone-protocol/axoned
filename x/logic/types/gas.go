package types

import (
	"runtime"
	"sync"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.GasMeter = (*safeGasMeter)(nil)

// safeGasMeter is a wrapper around sdk.GasMeter that provides go-routine-safe access to the underlying gas meter.
// This is needed because the interpreter is uses multiple go-routines, and the gas meter is shared between multiple
// goroutines.
type safeGasMeter struct {
	gasMeter sdk.GasMeter
	mutex    sync.RWMutex
}

func (m *safeGasMeter) ConsumeGas(amount uint64, descriptor string) {
	m.mutex.Lock()
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(sdk.ErrorOutOfGas); ok {
				runtime.Goexit()
			}
			panic(r)
		}
	}()
	defer m.mutex.Unlock()

	m.gasMeter.ConsumeGas(amount, descriptor)
}

func (m *safeGasMeter) GasConsumed() uint64 {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.GasConsumed()
}

func (m *safeGasMeter) GasConsumedToLimit() uint64 {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.GasConsumedToLimit()
}

func (m *safeGasMeter) IsPastLimit() bool {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.IsPastLimit()
}

func (m *safeGasMeter) Limit() uint64 {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.Limit()
}

func (m *safeGasMeter) RefundGas(amount uint64, descriptor string) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.gasMeter.RefundGas(amount, descriptor)
}

func (m *safeGasMeter) GasRemaining() uint64 {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.GasRemaining()
}

func (m *safeGasMeter) String() string {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.String()
}

func (m *safeGasMeter) IsOutOfGas() bool {
	m.mutex.RLocker().Lock()
	defer m.mutex.RLocker().Unlock()

	return m.gasMeter.IsOutOfGas()
}

// NewSafeGasMeter returns a new instance of sdk.GasMeter that is go-routine-safe.
func NewSafeGasMeter(gasMeter sdk.GasMeter) sdk.GasMeter {
	return &safeGasMeter{
		gasMeter: gasMeter,
	}
}
