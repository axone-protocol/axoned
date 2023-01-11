package context

import (
	"context"
	"sync/atomic"
)

// IncrementCountByFunc is a function that increments the count value by the provided delta.
type IncrementCountByFunc func(delta uint64) uint64

func (f IncrementCountByFunc) By(delta uint64) func() uint64 {
	return func() uint64 {
		return f(delta)
	}
}

// WithLimit returns a copy of the parent context with a limit adjusted to the provided value, and a function that
// increments the count value. When the count value exceeds the limit, the returned context's Done channel is closed.
// The returned context's Done channel is also closed when the returned cancel function is called, or when the parent
// context's Done channel is closed, whichever happens first.
// Canceling this context releases resources associated with it, so code should call cancel as soon as the operations
// running in this Context complete.
func WithLimit(parent context.Context, limit uint64) (context.Context, IncrementCountByFunc) {
	var counter uint64
	ctx, cancel := context.WithCancel(parent)

	return ctx, func(delta uint64) uint64 {
		newValue := atomic.AddUint64(&counter, delta)
		if newValue > limit {
			cancel()
		}
		return newValue
	}
}
