package predicate_test

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"time"
)

var errInvalidRequest = errors.New("invalid_request")

type echoDeviceFS struct{}

func newEchoDeviceFS() fs.FS {
	return &echoDeviceFS{}
}

func (f *echoDeviceFS) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
}

func (f *echoDeviceFS) OpenFile(name string, flag int, _ fs.FileMode) (fs.File, error) {
	if name != "." {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}

	if flag != os.O_RDWR {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
	}

	return &echoDeviceFile{name: "echo"}, nil
}

type echoDeviceFile struct {
	name string

	request   []byte
	response  []byte
	readPos   int
	committed bool
	closed    bool
}

func (f *echoDeviceFile) Stat() (fs.FileInfo, error) {
	return echoFileInfo{name: f.name}, nil
}

func (f *echoDeviceFile) Close() error {
	f.closed = true
	return nil
}

func (f *echoDeviceFile) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}

	if !f.committed {
		if len(f.request) == 0 {
			return 0, errInvalidRequest
		}

		f.response = append([]byte(nil), f.request...)
		f.committed = true
	}

	if f.readPos >= len(f.response) {
		return 0, io.EOF
	}

	n := copy(p, f.response[f.readPos:])
	f.readPos += n

	return n, nil
}

func (f *echoDeviceFile) Write(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}

	if f.committed {
		return 0, fs.ErrPermission
	}

	f.request = append(f.request, p...)
	return len(p), nil
}

type echoFileInfo struct {
	name string
}

func (fi echoFileInfo) Name() string       { return fi.name }
func (fi echoFileInfo) Size() int64        { return 0 }
func (fi echoFileInfo) Mode() fs.FileMode  { return 0o666 }
func (fi echoFileInfo) ModTime() time.Time { return time.Unix(0, 0).UTC() }
func (fi echoFileInfo) IsDir() bool        { return false }
func (fi echoFileInfo) Sys() any           { return nil }
