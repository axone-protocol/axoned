package wasm

import (
	"context"
	"errors"
	"io"
	"io/fs"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/devfile"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

var errWasmQueryFailed = errors.New("wasm_query_failed")

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
)

// NewFS creates the dev/wasm transactional device filesystem.
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
	return devfile.New(
		devfile.WithPath(name),
		devfile.WithModTime(prolog.ResolveHeaderInfo(sdkCtx).Time),
		devfile.WithMaxRequestBytes(maxRequestBytes),
		devfile.WithMaxResponseBytes(maxResponseBytes),
		devfile.WithCommit(func(r io.Reader, w io.Writer) error {
			request, err := io.ReadAll(r)
			if err != nil {
				return err
			}

			if len(request) == 0 {
				return devfile.ErrInvalidRequest
			}

			response, err := f.wasmKeeper.QuerySmart(sdkCtx, contractAddr, request)
			if err != nil {
				return errWasmQueryFailed
			}

			_, err = w.Write(response)
			return err
		}),
	)
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
