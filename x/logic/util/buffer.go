package util

import (
	"fmt"
	"io"
)

var (
	_ = io.Writer(&BoundedBuffer{})
)

var (
	ErrInvalidSize = func(size int) error {
		return fmt.Errorf("invalid buffer size %d", size)
	}
)

// BoundedBuffer is a fixed size buffer that overwrites older data when the buffer is full.
type BoundedBuffer struct {
	buf       []byte
	offset    int
	size      int
	overflown bool
}

// NewBoundedBuffer creates a new BoundedBuffer with the given fixed size.
// If size is 0, the buffer will be disabled.
func NewBoundedBuffer(size int) (*BoundedBuffer, error) {
	if size < 0 {
		return nil, ErrInvalidSize(size)
	}

	return &BoundedBuffer{
		buf:  make([]byte, size),
		size: size,
	}, nil
}

// NewBoundedBufferMust is like NewBoundedBuffer but panics if an error occurs.
func NewBoundedBufferMust(size int) *BoundedBuffer {
	b, err := NewBoundedBuffer(size)
	if err != nil {
		panic(err)
	}

	return b
}

// Write implements io.Writer.
func (b *BoundedBuffer) Write(p []byte) (n int, err error) {
	size := len(p)

	if b.size == 0 {
		return size, nil
	}

	for i := 0; i < size; i++ {
		b.buf[b.offset] = p[i]
		b.offset = (b.offset + 1) % b.size
		b.overflown = b.overflown || b.offset == 0
	}

	return size, nil
}

// String returns the contents of the buffer as a string.
func (b *BoundedBuffer) String() string {
	if b.overflown {
		r := string(b.buf[b.offset:]) + string(b.buf[:b.offset])

		return r
	}

	return string(b.buf[:b.offset])
}
