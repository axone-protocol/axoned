package devfile

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"strings"
	"time"
)

// ErrInvalidRequest indicates that the device received no request payload
// before the first read commit.
var ErrInvalidRequest = errors.New("invalid_request")

// ErrMissingCommit indicates that no commit function was provided.
var ErrMissingCommit = errors.New("missing_commit_function")

// CommitFunc executes the device transaction once the first read occurs.
// It reads the request from the provided Reader (bounded by maxRequestBytes)
// and writes the response to the provided Writer (bounded by maxResponseBytes).
type CommitFunc func(io.Reader, io.Writer) error

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
	path             string
	modTime          time.Time
	maxRequestBytes  int
	maxResponseBytes int
	commit           CommitFunc
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
// If exceeded, writes will fail with ErrWriteLimit.
func WithMaxRequestBytes(maxBytes int) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.maxRequestBytes = maxBytes
	}
}

// WithMaxResponseBytes sets the maximum number of bytes in the response.
// If exceeded, the commit will fail with ErrWriteLimit.
func WithMaxResponseBytes(maxBytes int) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.maxResponseBytes = maxBytes
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

	buffer    *boundedWriter
	commitErr error
	readPos   int
	committed bool
	closed    bool
}

// New creates an fs.File implementing the half-duplex protocol.
// Returns an error if the commit function is not provided.
func New(opts ...Option) (fs.File, error) {
	cfg := halfDuplexConfig{
		modTime: time.Now(),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if cfg.commit == nil {
		return nil, ErrMissingCommit
	}

	return &halfDuplexFile{
		cfg:    cfg,
		buffer: newBoundedWriter(cfg.maxRequestBytes),
	}, nil
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

	responseBytes := f.buffer.Bytes()
	if f.readPos >= len(responseBytes) {
		return 0, io.EOF
	}

	n := copy(p, responseBytes[f.readPos:])
	f.readPos += n

	return n, nil
}

func (f *halfDuplexFile) commit() {
	f.committed = true

	// Read request from buffer
	reader := bytes.NewReader(f.buffer.Bytes())

	// Reset buffer to reuse for response (half-duplex: write → commit → read)
	f.buffer.Reset()
	f.buffer.maxBytes = f.cfg.maxResponseBytes

	if err := f.cfg.commit(reader, f.buffer); err != nil {
		f.commitErr = err
	}
}

func (f *halfDuplexFile) Write(p []byte) (int, error) {
	if f.closed {
		return 0, f.pathError("write", fs.ErrClosed)
	}

	if f.committed {
		return 0, f.pathError("write", fs.ErrPermission)
	}

	n, err := f.buffer.Write(p)
	if err != nil {
		return n, f.pathError("write", err)
	}
	return n, nil
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
