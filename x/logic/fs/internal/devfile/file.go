package devfile

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

// ErrMissingCommit indicates that no commit function was provided.
var ErrMissingCommit = errors.New("missing_commit_function")

// CommitFunc executes the device transaction once the first read occurs.
// The input request is an immutable copy of all bytes written before commit.
type CommitFunc func(request []byte) ([]byte, error)

// halfDuplexConfig defines a half-duplex transactional protocol:
//   - Request phase: zero or more writes append request bytes.
//   - Commit phase: first read commits the request via Commit.
//   - Response phase: reads stream response bytes until EOF.
//   - Finalization: close invalidates subsequent reads/writes.
//
// State machine:
//   - writable (pre-commit)
//   - committed (readable, writes rejected)
//   - closed (all I/O rejected with fs.ErrClosed)
type halfDuplexConfig struct {
	path              string
	modTime           time.Time
	maxRequestBytes   int
	maxResponseBytes  int
	allowEmptyRequest bool
	emptyRequestErr   error
	commit            CommitFunc
}

// Option is a functional option for configuring a half-duplex file.
type Option func(*halfDuplexConfig)

// WithPath sets the file path.
func WithPath(path string) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.path = path
	}
}

// WithModTime sets the modification time.
func WithModTime(modTime time.Time) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.modTime = modTime
	}
}

// WithMaxRequestBytes sets the maximum number of bytes that can be written before commit.
// If exceeded, writes will fail with fs.ErrPermission.
func WithMaxRequestBytes(maxBytes int) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.maxRequestBytes = maxBytes
	}
}

// WithMaxResponseBytes sets the maximum number of bytes in the response.
// If exceeded, the commit will fail with ErrResponseTooLarge.
func WithMaxResponseBytes(maxBytes int) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.maxResponseBytes = maxBytes
	}
}

// WithAllowEmptyRequest allows committing without writing any bytes.
// By default, at least one byte must be written before the first read.
func WithAllowEmptyRequest(allow bool) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.allowEmptyRequest = allow
	}
}

// WithEmptyRequestError sets the error returned when reading before writing any bytes.
// Only effective when AllowEmptyRequest is false (the default).
// If not set, ErrInvalidRequest is used.
func WithEmptyRequestError(err error) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.emptyRequestErr = err
	}
}

// WithCommit sets the commit function that processes the request and returns a response.
// This option is required.
func WithCommit(commit CommitFunc) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.commit = commit
	}
}

type halfDuplexFile struct {
	cfg halfDuplexConfig

	request   []byte
	response  []byte
	commitErr error
	readPos   int
	committed bool
	closed    bool
}

// New creates an fs.File implementing the half-duplex protocol.
// Returns an error if the commit function is not provided.
func New(opts ...Option) (fs.File, error) {
	cfg := halfDuplexConfig{
		modTime:         time.Now(),
		emptyRequestErr: ErrInvalidRequest,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.commit == nil {
		return nil, ErrMissingCommit
	}

	return &halfDuplexFile{cfg: cfg}, nil
}

func (f *halfDuplexFile) Stat() (fs.FileInfo, error) {
	return fileInfo{name: baseName(f.cfg.path), modTime: f.cfg.modTime}, nil
}

func (f *halfDuplexFile) Read(p []byte) (int, error) {
	if f.closed {
		return 0, f.pathError("read", fs.ErrClosed)
	}

	if !f.committed {
		f.commit()
	}

	if f.commitErr != nil {
		return 0, f.pathError("read", f.commitErr)
	}

	if f.readPos >= len(f.response) {
		return 0, io.EOF
	}

	n := copy(p, f.response[f.readPos:])
	f.readPos += n

	return n, nil
}

func (f *halfDuplexFile) commit() {
	f.committed = true

	if !f.cfg.allowEmptyRequest && len(f.request) == 0 {
		f.commitErr = f.cfg.emptyRequestErr
		return
	}

	response, err := f.cfg.commit(append([]byte(nil), f.request...))
	switch {
	case err != nil:
		f.commitErr = err
	case f.cfg.maxResponseBytes > 0 && len(response) > f.cfg.maxResponseBytes:
		f.commitErr = ErrResponseTooLarge
	default:
		f.response = response
	}
}

func (f *halfDuplexFile) Write(p []byte) (int, error) {
	if f.closed {
		return 0, f.pathError("write", fs.ErrClosed)
	}

	if f.committed {
		return 0, f.pathError("write", fs.ErrPermission)
	}

	if f.cfg.maxRequestBytes > 0 && len(f.request)+len(p) > f.cfg.maxRequestBytes {
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
		Path: f.cfg.path,
		Err:  err,
	}
}
