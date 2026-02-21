package embedded

import (
	"io/fs"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/pathutil"
)

type vfs struct {
	fsys fs.FS
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates a read-only filesystem wrapper around an embedded filesystem.
func NewFS(fsys fs.FS) fs.ReadFileFS {
	return &vfs{fsys: fsys}
}

func (f *vfs) Open(name string) (fs.File, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	file, err := f.fsys.Open(subpath)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: pathutil.UnwrapPathError(err)}
	}

	return file, nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	subpath, err := pathutil.NormalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	content, err := fs.ReadFile(f.fsys, subpath)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: pathutil.UnwrapPathError(err)}
	}

	return content, nil
}
