package devfile

import "io"

// WriteHook is a function type that can be used to hook into write operations.
type WriteHook func(n int)

// hookedWriter is an io.Writer that invokes a WriteHook after each successful Write.
type hookedWriter struct {
	writer io.Writer
	hook   WriteHook
}

// Write writes data to the underlying writer and invokes the hook with the number of bytes written.
func (w hookedWriter) Write(p []byte) (int, error) {
	n, err := w.writer.Write(p)
	if n > 0 && w.hook != nil {
		w.hook(n)
	}

	return n, err
}
