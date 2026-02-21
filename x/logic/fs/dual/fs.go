package dual

import (
	"errors"
	"io"
	"io/fs"
	"net/url"
)

type openFileFS interface {
	fs.FS
	OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error)
}

type vfs struct {
	pathFS   fs.FS
	legacyFS fs.FS
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
	_ openFileFS    = (*vfs)(nil)
)

// NewFS creates a dual-stack filesystem:
//   - URI source-sinks are dispatched to legacyFS.
//   - non-URI source-sinks (path-like) are dispatched to pathFS.
//
// This allows us to support both the new path-based VFS and the legacy URI-based FS
// in parallel, without forcing a full migration of all file accesses to the new path-based VFS at once.
// The dual-stack approach is a temporary solution to facilitate the transition to the new path-based VFS,
// and should be removed once all file accesses have been migrated to the new path-based VFS.
func NewFS(pathFS fs.FS, legacyFS fs.FS) fs.FS {
	return &vfs{
		pathFS:   pathFS,
		legacyFS: legacyFS,
	}
}

func (f *vfs) Open(name string) (fs.File, error) {
	if isURI(name) {
		return f.legacyFS.Open(name)
	}

	return f.pathFS.Open(name)
}

func (f *vfs) OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error) {
	selected := f.pathFS
	if isURI(name) {
		selected = f.legacyFS
	}

	if ofs, ok := selected.(openFileFS); ok {
		return ofs.OpenFile(name, flag, perm)
	}

	return nil, &fs.PathError{
		Op:   "open",
		Path: name,
		Err:  errors.Join(errors.ErrUnsupported, fs.ErrPermission),
	}
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	selected := f.pathFS
	if isURI(name) {
		selected = f.legacyFS
	}

	if rfs, ok := selected.(fs.ReadFileFS); ok {
		return rfs.ReadFile(name)
	}

	file, err := selected.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func isURI(name string) bool {
	u, err := url.Parse(name)
	return err == nil && u.Scheme != ""
}
