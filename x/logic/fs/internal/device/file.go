package device

import (
	"errors"
	"io"
	"io/fs"
	"strings"
	"time"
)

// ErrInvalidRequest indicates that the device received no request payload
// before the first read commit.
var ErrInvalidRequest = errors.New("invalid_request")

// ErrResponseTooLarge indicates that the committed response exceeds
// the configured maximum size.
var ErrResponseTooLarge = errors.New("response_too_large")

// CommitFunc executes the device transaction once the first read occurs.
// The input request is an immutable copy of all bytes written before commit.
type CommitFunc func(request []byte) ([]byte, error)

// HalfDuplexConfig defines a half-duplex transactional protocol:
//   - Request phase: zero or more writes append request bytes.
//   - Commit phase: first read commits the request via Commit.
//   - Response phase: reads stream response bytes until EOF.
//   - Finalization: close invalidates subsequent reads/writes.
//
// State machine:
//   - writable (pre-commit)
//   - committed (readable, writes rejected)
//   - closed (all I/O rejected with fs.ErrClosed)
//
// By default the protocol requires at least one written byte before commit.
// Set AllowEmptyRequest=true to allow empty requests.
// If AllowEmptyRequest=false and EmptyRequestErr is nil, ErrInvalidRequest is used.
// If MaxRequestBytes > 0, writes exceeding that bound fail with fs.ErrPermission.
// If MaxResponseBytes > 0, committed responses larger than that fail with ErrResponseTooLarge.
type HalfDuplexConfig struct {
	Path              string
	ModTime           time.Time
	MaxRequestBytes   int
	MaxResponseBytes  int
	AllowEmptyRequest bool
	EmptyRequestErr   error
	Commit            CommitFunc
}

type halfDuplexFile struct {
	cfg HalfDuplexConfig

	request   []byte
	response  []byte
	readPos   int
	committed bool
	closed    bool
}

// New creates an fs.File implementing the half-duplex protocol.
func New(cfg HalfDuplexConfig) fs.File {
	if cfg.EmptyRequestErr == nil {
		cfg.EmptyRequestErr = ErrInvalidRequest
	}

	return &halfDuplexFile{cfg: cfg}
}

func (f *halfDuplexFile) Stat() (fs.FileInfo, error) {
	return fileInfo{name: baseName(f.cfg.Path), modTime: f.cfg.ModTime}, nil
}

func (f *halfDuplexFile) Read(p []byte) (int, error) {
	if f.closed {
		return 0, f.pathError("read", fs.ErrClosed)
	}

	if !f.committed {
		if !f.cfg.AllowEmptyRequest && len(f.request) == 0 {
			return 0, f.pathError("read", f.cfg.EmptyRequestErr)
		}

		response, err := f.cfg.Commit(append([]byte(nil), f.request...))
		if err != nil {
			return 0, f.pathError("read", err)
		}

		if f.cfg.MaxResponseBytes > 0 && len(response) > f.cfg.MaxResponseBytes {
			return 0, f.pathError("read", ErrResponseTooLarge)
		}

		f.response = response
		f.committed = true
	}

	if f.readPos >= len(f.response) {
		return 0, io.EOF
	}

	n := copy(p, f.response[f.readPos:])
	f.readPos += n

	return n, nil
}

func (f *halfDuplexFile) Write(p []byte) (int, error) {
	if f.closed {
		return 0, f.pathError("write", fs.ErrClosed)
	}

	if f.committed {
		return 0, f.pathError("write", fs.ErrPermission)
	}

	if f.cfg.MaxRequestBytes > 0 && len(f.request)+len(p) > f.cfg.MaxRequestBytes {
		return 0, f.pathError("write", fs.ErrPermission)
	}

	f.request = append(f.request, p...)
	return len(p), nil
}

func (f *halfDuplexFile) Close() error {
	f.closed = true
	return nil
}

type fileInfo struct {
	name    string
	modTime time.Time
}

func (fi fileInfo) Name() string       { return fi.name }
func (fi fileInfo) Size() int64        { return 0 }
func (fi fileInfo) Mode() fs.FileMode  { return 0o666 }
func (fi fileInfo) ModTime() time.Time { return fi.modTime }
func (fi fileInfo) IsDir() bool        { return false }
func (fi fileInfo) Sys() any           { return nil }

func baseName(path string) string {
	if idx := strings.LastIndexByte(path, '/'); idx >= 0 {
		return path[idx+1:]
	}

	return path
}

func (f *halfDuplexFile) pathError(op string, err error) error {
	if err == nil || errors.Is(err, io.EOF) {
		return err
	}

	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return err
	}

	return &fs.PathError{
		Op:   op,
		Path: f.cfg.Path,
		Err:  err,
	}
}
