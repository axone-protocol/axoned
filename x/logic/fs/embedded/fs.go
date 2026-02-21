package embedded

import (
	"errors"
	"io/fs"
	"path"
	"strings"
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
	subpath, err := normalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	file, err := f.fsys.Open(subpath)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: unwrapPathError(err)}
	}

	return file, nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	subpath, err := normalizeSubpath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: err}
	}

	content, err := fs.ReadFile(f.fsys, subpath)
	if err != nil {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: unwrapPathError(err)}
	}

	return content, nil
}

func normalizeSubpath(name string) (string, error) {
	trimmed := strings.TrimPrefix(name, "/")
	if trimmed == "" || trimmed == "." {
		return ".", nil
	}

	for _, segment := range strings.Split(trimmed, "/") {
		if segment == ".." {
			return "", fs.ErrPermission
		}
	}

	cleaned := path.Clean(trimmed)
	if cleaned == "." {
		return ".", nil
	}

	if !fs.ValidPath(cleaned) {
		return "", fs.ErrInvalid
	}

	return cleaned, nil
}

func unwrapPathError(err error) error {
	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return pathErr.Err
	}

	return err
}
