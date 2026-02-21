package sys

import (
	"context"
	"fmt"
	"io/fs"
	"strconv"
	"strings"

	coreheader "cosmossdk.io/core/header"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
)

const (
	headerAtPath      = "header/@"
	headerHeightPath  = "header/height"
	headerHashPath    = "header/hash"
	headerTimePath    = "header/time"
	headerChainIDPath = "header/chain_id"
	headerAppHashPath = "header/app_hash"
)

type vfs struct {
	ctx context.Context
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates the /v1/sys snapshot filesystem.
func NewFS(ctx context.Context) fs.ReadFileFS {
	return &vfs{ctx: ctx}
}

func (f *vfs) Open(name string) (fs.File, error) {
	data, err := f.readFile("open", name)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return NewVirtualFile(name, data, prolog.ResolveHeaderInfo(sdkCtx).Time), nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	return f.readFile("readfile", name)
}

func (f *vfs) readFile(op, name string) ([]byte, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: op, Path: name, Err: err}
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	header := prolog.ResolveHeaderInfo(sdkCtx)

	content, err := renderFile(header, subpath)
	if err != nil {
		return nil, &fs.PathError{Op: op, Path: name, Err: err}
	}

	return content, nil
}

func renderFile(header coreheader.Info, subpath string) ([]byte, error) {
	switch subpath {
	case headerAtPath:
		return []byte(fmt.Sprintf("header{height:%d,hash:%s,time:%d,chain_id:%s,app_hash:%s}.\n",
			header.Height,
			formatByteList(header.Hash),
			header.Time.Unix(),
			quoteAtom(header.ChainID),
			formatByteList(header.AppHash),
		)), nil
	case headerHeightPath:
		return []byte(fmt.Sprintf("%d.\n", header.Height)), nil
	case headerHashPath:
		return []byte(formatByteList(header.Hash) + ".\n"), nil
	case headerTimePath:
		return []byte(fmt.Sprintf("%d.\n", header.Time.Unix())), nil
	case headerChainIDPath:
		return []byte(quoteAtom(header.ChainID) + ".\n"), nil
	case headerAppHashPath:
		return []byte(formatByteList(header.AppHash) + ".\n"), nil
	default:
		return nil, fs.ErrNotExist
	}
}

func formatByteList(b []byte) string {
	if len(b) == 0 {
		return "[]"
	}

	parts := make([]string, 0, len(b))
	for _, value := range b {
		parts = append(parts, strconv.FormatInt(int64(value), 10))
	}

	return "[" + strings.Join(parts, ",") + "]"
}

func quoteAtom(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "''") + "'"
}
