package source

import (
	"context"
	"io/fs"

	"github.com/axone-protocol/prolog/v3/engine"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/pathutil"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/prologterm"
	"github.com/axone-protocol/axoned/v15/x/logic/fs/internal/virtualfile"
	"github.com/axone-protocol/axoned/v15/x/logic/prolog"
	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

const filesPath = "files"

type FilesProvider func() []string

type vfs struct {
	ctx context.Context
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates a read-only snapshot filesystem for currently loaded Prolog sources.
func NewFS(ctx context.Context) fs.ReadFileFS {
	return &vfs{ctx: ctx}
}

func (f *vfs) Open(name string) (fs.File, error) {
	sdkCtx := sdk.UnwrapSDKContext(f.ctx)

	data, err := f.readFile("open", name)
	if err != nil {
		return nil, err
	}

	return virtualfile.New(name, data, prolog.ResolveHeaderInfo(sdkCtx).Time), nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	return f.readFile("open", name)
}

func (f *vfs) readFile(op, name string) ([]byte, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: op, Path: name, Err: err}
	}

	content, err := f.renderFile(subpath)
	if err != nil {
		return nil, &fs.PathError{Op: op, Path: name, Err: err}
	}

	return content, nil
}

func (f *vfs) renderFile(subpath string) ([]byte, error) {
	switch subpath {
	case filesPath:
		return prologterm.Render(sourceFilesTerm(sourceFiles(f.ctx)), true)
	default:
		return nil, fs.ErrNotExist
	}
}

func sourceFiles(ctx context.Context) []string {
	provider, _ := ctx.Value(types.SourceFilesProviderContextKey).(FilesProvider)
	if provider == nil {
		return nil
	}

	return provider()
}

func sourceFilesTerm(files []string) engine.Term {
	terms := make([]engine.Term, 0, len(files))
	for _, file := range files {
		terms = append(terms, engine.NewAtom(file))
	}

	return engine.List(terms...)
}
