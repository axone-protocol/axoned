package streamingfile

import (
	"errors"
	"io"
	"io/fs"
	"strings"
	"time"
)

type (
	Next[T any]   func() (T, bool, error)
	Stop          func() error
	Render[T any] func(T) ([]byte, error)
)

type OpenCursor[T any] func() (Next[T], Stop, error)

// File represents a streaming file.
type File[T any] struct {
	name    string
	modTime time.Time

	open   OpenCursor[T]
	render Render[T]

	next Next[T]
	stop Stop

	pending []byte
	err     error
	eof     bool
	closed  bool
}

var _ fs.File = (*File[byte])(nil)

func New[T any](name string, modTime time.Time, open OpenCursor[T], render Render[T]) *File[T] {
	return &File[T]{
		name:    name,
		modTime: modTime,
		open:    open,
		render:  render,
	}
}

func (f *File[T]) isInitialized() bool { return f.next != nil }

func (f *File[T]) Stat() (fs.FileInfo, error) {
	return fileInfo{name: base(f.name), modTime: f.modTime}, nil
}

func (f *File[T]) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}
	if f.err != nil {
		return 0, f.err
	}

	if !f.isInitialized() {
		next, stop, err := f.open()
		if err != nil {
			f.err = err
			return 0, err
		}
		f.next, f.stop = next, stop
	}

	for {
		if len(f.pending) > 0 {
			n := copy(p, f.pending)
			f.pending = f.pending[n:]
			return n, nil
		}

		if f.eof {
			return 0, io.EOF
		}

		item, ok, err := f.next()
		if err != nil {
			return f.failWithCleanup(err)
		}
		if !ok {
			f.eof = true
			if err := f.cleanup(); err != nil {
				f.err = err
				return 0, err
			}
			return 0, io.EOF
		}

		b, err := f.render(item)
		if err != nil {
			return f.failWithCleanup(err)
		}

		f.pending = b
	}
}

func (f *File[T]) failWithCleanup(err error) (int, error) {
	if cleanupErr := f.cleanup(); cleanupErr != nil {
		err = errors.Join(err, cleanupErr)
	}
	f.err = err
	return 0, err
}

func (f *File[T]) cleanup() error {
	if f.isInitialized() && f.stop != nil {
		err := f.stop()
		f.stop = nil
		return err
	}
	return nil
}

func (f *File[T]) Close() error {
	if f.closed {
		return nil
	}
	f.closed = true
	return f.cleanup()
}

type fileInfo struct {
	name    string
	modTime time.Time
}

func (fi fileInfo) Name() string       { return fi.name }
func (fi fileInfo) Size() int64        { return 0 }
func (fi fileInfo) Mode() fs.FileMode  { return 0o444 }
func (fi fileInfo) ModTime() time.Time { return fi.modTime }
func (fi fileInfo) IsDir() bool        { return false }
func (fi fileInfo) Sys() any           { return nil }

func base(p string) string {
	if i := strings.LastIndexByte(p, '/'); i >= 0 {
		return p[i+1:]
	}
	return p
}
