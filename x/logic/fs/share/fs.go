package share

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/virtualfile"
	logictypes "github.com/axone-protocol/axoned/v14/x/logic/types"
)

const programFileExt = ".pl"

var errVFSUnavailable = errors.New("vfs_unavailable")

type programKeeper interface {
	GetStoredProgram(ctx sdk.Context, programID []byte) (logictypes.StoredProgram, bool, error)
	GetProgramPublication(ctx sdk.Context, publisher, programID []byte) (logictypes.ProgramPublication, bool, error)
}

type vfs struct {
	ctx    context.Context
	keeper programKeeper
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates a read-only filesystem backed by stored program artifacts.
func NewFS(ctx context.Context, keeper programKeeper) fs.ReadFileFS {
	return &vfs{
		ctx:    ctx,
		keeper: keeper,
	}
}

func (f *vfs) Open(name string) (fs.File, error) {
	content, modTime, err := f.readFile("open", name)
	if err != nil {
		return nil, err
	}

	return virtualfile.New(name, content, modTime), nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	content, _, err := f.readFile("open", name)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (f *vfs) readFile(op, name string) ([]byte, time.Time, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, time.Time{}, &fs.PathError{Op: op, Path: name, Err: err}
	}

	publisher, programID, err := decodePublicationPath(subpath)
	if err != nil {
		return nil, time.Time{}, &fs.PathError{Op: op, Path: name, Err: err}
	}

	if f.keeper == nil {
		return nil, time.Time{}, &fs.PathError{Op: op, Path: name, Err: errVFSUnavailable}
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	publication, found, err := f.keeper.GetProgramPublication(sdkCtx, publisher, programID)
	if err != nil {
		return nil, time.Time{}, &fs.PathError{Op: op, Path: name, Err: err}
	}
	if !found {
		return nil, time.Time{}, &fs.PathError{Op: op, Path: name, Err: fs.ErrNotExist}
	}

	program, found, err := f.keeper.GetStoredProgram(sdkCtx, programID)
	if err != nil {
		return nil, time.Time{}, &fs.PathError{Op: op, Path: name, Err: err}
	}
	if !found {
		return nil, time.Time{}, &fs.PathError{Op: op, Path: name, Err: fs.ErrNotExist}
	}

	modTime := time.Unix(publication.GetPublishedAt(), 0).UTC()

	return []byte(program.GetSource()), modTime, nil
}

func decodePublicationPath(subpath string) ([]byte, []byte, error) {
	if subpath == "." {
		return nil, nil, fs.ErrNotExist
	}

	userText, rest, found := strings.Cut(subpath, "/")
	if !found || userText == "" {
		return nil, nil, fs.ErrNotExist
	}

	programsRoot, programName, found := strings.Cut(rest, "/")
	if !found || programsRoot != "programs" || programName == "" || strings.Contains(programName, "/") {
		return nil, nil, fs.ErrNotExist
	}

	publisher, err := sdk.AccAddressFromBech32(userText)
	if err != nil {
		return nil, nil, fs.ErrNotExist
	}
	if !strings.HasSuffix(programName, programFileExt) {
		return nil, nil, fs.ErrNotExist
	}

	programIDHex := strings.TrimSuffix(programName, programFileExt)
	if len(programIDHex) != sha256.Size*2 {
		return nil, nil, fs.ErrNotExist
	}

	programID, err := hex.DecodeString(programIDHex)
	if err != nil {
		return nil, nil, fs.ErrNotExist
	}

	return publisher, programID, nil
}
