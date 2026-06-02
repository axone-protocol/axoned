package cryptofs

import (
	"context"
	"io"
	"io/fs"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/devfile"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/iface"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
	"github.com/axone-protocol/axoned/v15/x/logic/util"
)

const (
	maxRequestBytes  = 256 * 1024
	maxResponseBytes = 64
)

type vfs struct {
	ctx context.Context
}

var (
	_ fs.FS            = (*vfs)(nil)
	_ iface.OpenFileFS = (*vfs)(nil)
)

// NewFS creates a transactional device filesystem for crypto utilities.
func NewFS(ctx context.Context) fs.FS {
	return &vfs{ctx: ctx}
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

	algorithm, err := validateHashPath(subpath)
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

			response, err := util.Hash(algorithm, request)
			if err != nil {
				return err
			}

			_, err = w.Write(response)
			return err
		}),
	)
}

func validateHashPath(subpath string) (util.HashAlg, error) {
	if strings.Contains(subpath, "/") {
		return util.HashAlg(0), fs.ErrNotExist
	}

	algorithm, err := util.ParseHashAlg(subpath)
	if err != nil {
		return util.HashAlg(0), fs.ErrNotExist
	}

	return algorithm, nil
}
