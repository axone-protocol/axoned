package wasm

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/device"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

const (
	queryPath        = "query"
	maxRequestBytes  = 64 * 1024
	maxResponseBytes = 64 * 1024
)

type vfs struct {
	ctx        context.Context
	wasmKeeper types.WasmKeeper
}

var (
	_ fs.FS            = (*vfs)(nil)
	_ iface.OpenFileFS = (*vfs)(nil)

	errInvalidRequest  = device.ErrInvalidRequest
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
	return device.New(device.HalfDuplexConfig{
		Path:              name,
		ModTime:           prolog.ResolveHeaderInfo(sdkCtx).Time,
		MaxRequestBytes:   maxRequestBytes,
		MaxResponseBytes:  maxResponseBytes,
		AllowEmptyRequest: false,
		Commit: func(request []byte) ([]byte, error) {
			response, err := f.wasmKeeper.QuerySmart(sdkCtx, contractAddr, request)
			if err != nil {
				return nil, &fs.PathError{Op: "read", Path: name, Err: errWasmQueryFailed}
			}

			return response, nil
		},
	}), nil
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
