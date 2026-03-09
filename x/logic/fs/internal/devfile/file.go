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

// TransferDirection identifies which side of the half-duplex transport carried bytes.
type TransferDirection uint8

const (
	// TransferRequest denotes bytes written by the caller into the request buffer.
	TransferRequest TransferDirection = iota
	// TransferResponse denotes bytes produced by the device into the response buffer.
	TransferResponse
)

// TransferHook observes bytes successfully transferred through the device transport.
//
// The hook is invoked after a write succeeds, with the actual byte count returned by
// the underlying writer, so it reflects the real transferred volume even on partial
// writes. The direction indicates whether the bytes were appended to the request
// buffer or buffered as device response data during commit.
type TransferHook func(dir TransferDirection, n int)

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
	onTransfer       TransferHook
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

// WithTransferHook sets a hook invoked after bytes are transferred through the device transport.
func WithTransferHook(hook TransferHook) Option {
	return func(cfg *halfDuplexConfig) {
		cfg.onTransfer = hook
	}
}

type halfDuplexFile struct {
	cfg             halfDuplexConfig
	transportBuffer *boundedWriter
	w               *hookedWriter
	commitErr       error
	readPos         int
	committed       bool
	closed          bool
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

	transportBuffer := newBoundedWriter(cfg.maxRequestBytes)
	f := &halfDuplexFile{
		cfg:             cfg,
		transportBuffer: transportBuffer,
	}
	f.installHookedWriter(TransferRequest)

	return f, nil
}

func (f *halfDuplexFile) Stat() (fs.FileInfo, error) {
	return fileInfo{name: baseName(f.cfg.path), modTime: f.cfg.modTime}, nil
}

func (f *halfDuplexFile) Read(p []byte) (int, error) {
	if err := f.ensureReadable(); err != nil {
		return 0, err
	}

	return f.readResponse(p)
}

func (f *halfDuplexFile) commit() {
	f.committed = true

	request := f.transportBuffer.Bytes()
	reader := bytes.NewReader(request)
	f.beginResponsePhase()

	if err := f.cfg.commit(reader, f.w); err != nil {
		f.commitErr = err
	}
}

func (f *halfDuplexFile) ensureReadable() error {
	if f.closed {
		return f.pathError("read", fs.ErrClosed)
	}

	if !f.committed {
		f.commit()
	}

	if f.commitErr != nil {
		return f.pathError("read", f.commitErr)
	}

	return nil
}

func (f *halfDuplexFile) ensureWritable() error {
	if f.closed {
		return f.pathError("write", fs.ErrClosed)
	}

	if f.committed {
		return f.pathError("write", fs.ErrPermission)
	}

	return nil
}

func (f *halfDuplexFile) beginResponsePhase() {
	f.transportBuffer.Reset()
	f.transportBuffer.maxBytes = f.cfg.maxResponseBytes
	f.installHookedWriter(TransferResponse)
}

func (f *halfDuplexFile) readResponse(p []byte) (int, error) {
	responseBytes := f.transportBuffer.Bytes()
	if f.readPos >= len(responseBytes) {
		return 0, io.EOF
	}

	n := copy(p, responseBytes[f.readPos:])
	f.readPos += n

	return n, nil
}

func (f *halfDuplexFile) installHookedWriter(dir TransferDirection) {
	f.w = &hookedWriter{
		writer: f.transportBuffer,
		hook: func(n int) {
			if f.cfg.onTransfer != nil {
				f.cfg.onTransfer(dir, n)
			}
		},
	}
}

func (f *halfDuplexFile) Write(p []byte) (int, error) {
	if err := f.ensureWritable(); err != nil {
		return 0, err
	}

	n, err := f.w.Write(p)
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
