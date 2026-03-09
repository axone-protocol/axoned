package devfile

import (
	"bytes"
	"errors"
)

// ErrWriteLimit is returned when a write would exceed the configured limit.
var ErrWriteLimit = errors.New("write limit exceeded")

// defaultPreallocSize is the initial buffer capacity for bounded writers.
// Set to 16KB as a reasonable size for blockchain request/response workloads.
const defaultPreallocSize = 16 * 1024

// boundedWriter wraps a buffer and enforces a maximum byte limit.
// Once the limit is reached, further writes return ErrWriteLimit.
type boundedWriter struct {
	buf      bytes.Buffer
	maxBytes int
}

// newBoundedWriter creates a bounded writer with the specified limit.
// If maxBytes is 0 or negative, no limit is enforced.
// The buffer is pre-allocated with a reasonable capacity to reduce reallocations.
func newBoundedWriter(maxBytes int) *boundedWriter {
	w := &boundedWriter{
		maxBytes: maxBytes,
	}
	if maxBytes > 0 {
		// Pre-allocate with a reasonable size (defaultPreallocSize or maxBytes, whichever is smaller)
		capacity := min(defaultPreallocSize, maxBytes)
		w.buf.Grow(capacity)
	}
	return w
}

// Write appends data to the buffer.
//
// Writes are all-or-nothing: if writing p would exceed the configured limit,
// no bytes are written and ErrWriteLimit is returned.
func (w *boundedWriter) Write(p []byte) (int, error) {
	if w.maxBytes > 0 && w.buf.Len()+len(p) > w.maxBytes {
		return 0, ErrWriteLimit
	}
	return w.buf.Write(p)
}

// Bytes returns the accumulated bytes.
func (w *boundedWriter) Bytes() []byte {
	return w.buf.Bytes()
}

// Len returns the number of bytes written so far.
func (w *boundedWriter) Len() int {
	return w.buf.Len()
}

// Reset clears the buffer for reuse.
func (w *boundedWriter) Reset() {
	w.buf.Reset()
}
