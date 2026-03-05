package wasm

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

const (
	queryPath       = "query"
	maxRequestBytes = 64 * 1024
)

type vfs struct {
	ctx        context.Context
	wasmKeeper types.WasmKeeper
}

var (
	_ fs.FS            = (*vfs)(nil)
	_ iface.OpenFileFS = (*vfs)(nil)

	errInvalidRequest  = errors.New("invalid_request")
	errWasmQueryFailed = errors.New("wasm_query_failed")
)

// NewFS creates the /v1/dev/wasm transactional device filesystem.
func NewFS(ctx context.Context, wasmKeeper types.WasmKeeper) fs.FS {
	return &vfs{ctx: ctx, wasmKeeper: wasmKeeper}
}

func (f *vfs) Open(name string) (fs.File, error) {
	return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
}

func (f *vfs) OpenFile(name string, flag int, _ fs.FileMode) (fs.File, error) {
	if flag != os.O_RDWR {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrPermission}
	}

	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	contractAddr, err := validateQueryPath(subpath)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return newDeviceFile(name, prolog.ResolveHeaderInfo(sdkCtx).Time, sdkCtx, f.wasmKeeper, contractAddr), nil
}

func validateQueryPath(subpath string) (sdk.AccAddress, error) {
	segments := strings.Split(subpath, "/")
	if len(segments) != 2 || segments[1] != queryPath {
		return nil, fs.ErrNotExist
	}

	address := segments[0]
	if address == "" {
		return nil, fs.ErrNotExist
	}

	addr, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return nil, fs.ErrNotExist
	}

	return addr, nil
}

type deviceFile struct {
	path string

	modTime    time.Time
	ctx        context.Context
	wasmKeeper types.WasmKeeper
	addr       sdk.AccAddress

	request   []byte
	response  []byte
	readPos   int
	committed bool
	closed    bool
}

func newDeviceFile(
	path string,
	modTime time.Time,
	ctx context.Context,
	wasmKeeper types.WasmKeeper,
	addr sdk.AccAddress,
) fs.File {
	return &deviceFile{
		path:       path,
		modTime:    modTime,
		ctx:        ctx,
		wasmKeeper: wasmKeeper,
		addr:       addr,
	}
}

func (f *deviceFile) Stat() (fs.FileInfo, error) {
	return fileInfo{name: baseName(f.path), modTime: f.modTime}, nil
}

func (f *deviceFile) Read(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}

	if !f.committed {
		if len(f.request) == 0 {
			return 0, errInvalidRequest
		}

		response, err := f.wasmKeeper.QuerySmart(f.ctx, f.addr, append([]byte(nil), f.request...))
		if err != nil {
			return 0, &fs.PathError{Op: "read", Path: f.path, Err: errWasmQueryFailed}
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

func (f *deviceFile) Write(p []byte) (int, error) {
	if f.closed {
		return 0, fs.ErrClosed
	}

	if f.committed {
		return 0, fs.ErrPermission
	}

	if len(f.request)+len(p) > maxRequestBytes {
		return 0, fs.ErrPermission
	}

	f.request = append(f.request, p...)
	return len(p), nil
}

func (f *deviceFile) Close() error {
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
