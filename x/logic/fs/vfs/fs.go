package vfs

import (
	"errors"
	"io"
	"io/fs"
	"sort"
	"strings"

	"github.com/axone-protocol/axoned/v14/x/logic/fs/internal/iface"
)

// OpenFileFS is implemented by file systems that support opening files with flags.
type OpenFileFS = iface.OpenFileFS

// Mount defines a filesystem mounted at a canonical absolute prefix.
type Mount struct {
	Prefix string
	FS     fs.FS
}

// Router routes a path to a mounted filesystem and returns the subpath.
type Router interface {
	Route(path string) (mounted fs.FS, subpath string, err error)
}

// FileSystem is a policy-free VFS router (normalize + mount + route + dispatch).
type FileSystem struct {
	mounts []Mount
}

var (
	_ fs.FS         = (*FileSystem)(nil)
	_ fs.ReadFileFS = (*FileSystem)(nil)
	_ OpenFileFS    = (*FileSystem)(nil)
	_ Router        = (*FileSystem)(nil)
)

// New creates an empty VFS router.
func New() *FileSystem {
	return &FileSystem{}
}

// Mount mounts a filesystem at the given prefix.
func (v *FileSystem) Mount(prefix string, fsys fs.FS) error {
	if fsys == nil {
		return &fs.PathError{Op: "mount", Path: prefix, Err: fs.ErrInvalid}
	}

	canonicalPrefix, err := normalizePath(prefix)
	if err != nil {
		return &fs.PathError{Op: "mount", Path: prefix, Err: err}
	}

	for _, m := range v.mounts {
		if m.Prefix == canonicalPrefix {
			return &fs.PathError{Op: "mount", Path: prefix, Err: fs.ErrExist}
		}
	}

	v.mounts = append(v.mounts, Mount{
		Prefix: canonicalPrefix,
		FS:     fsys,
	})

	sort.Slice(v.mounts, func(i, j int) bool {
		li := len(v.mounts[i].Prefix)
		lj := len(v.mounts[j].Prefix)

		if li == lj {
			return v.mounts[i].Prefix < v.mounts[j].Prefix
		}

		return li > lj
	})

	return nil
}

// Route resolves the mounted filesystem and subpath using longest segment-safe prefix.
func (v *FileSystem) Route(path string) (mounted fs.FS, subpath string, err error) {
	canonicalPath, err := normalizePath(path)
	if err != nil {
		return nil, "", err
	}

	for _, m := range v.mounts {
		if !matchPrefix(canonicalPath, m.Prefix) {
			continue
		}

		return m.FS, toSubpath(canonicalPath, m.Prefix), nil
	}

	return nil, "", fs.ErrNotExist
}

// Open routes and dispatches to mounted FS.Open.
func (v *FileSystem) Open(path string) (fs.File, error) {
	mounted, subpath, err := v.Route(path)
	if err != nil {
		return nil, wrapPathError("open", path, err)
	}

	f, err := mounted.Open(subpath)
	if err != nil {
		return nil, wrapPathError("open", path, err)
	}

	return wrapFile(path, f), nil
}

// ReadFile routes and dispatches to mounted FS.ReadFile when supported.
func (v *FileSystem) ReadFile(path string) ([]byte, error) {
	mounted, subpath, err := v.Route(path)
	if err != nil {
		return nil, wrapPathError("open", path, err)
	}

	if rfs, ok := mounted.(fs.ReadFileFS); ok {
		content, err := rfs.ReadFile(subpath)
		if err != nil {
			return nil, wrapPathError("open", path, err)
		}
		return content, nil
	}

	f, err := mounted.Open(subpath)
	if err != nil {
		return nil, wrapPathError("open", path, err)
	}
	defer f.Close()

	content, err := io.ReadAll(f)
	if err != nil {
		return nil, wrapPathError("open", path, err)
	}

	return content, nil
}

// OpenFile routes and dispatches to mounted FS.OpenFile when supported.
func (v *FileSystem) OpenFile(path string, flag int, perm fs.FileMode) (fs.File, error) {
	mounted, subpath, err := v.Route(path)
	if err != nil {
		return nil, wrapPathError("open", path, err)
	}

	ofs, ok := mounted.(OpenFileFS)
	if !ok {
		return nil, wrapPathError("open", path, errors.Join(errors.ErrUnsupported, fs.ErrPermission))
	}

	f, err := ofs.OpenFile(subpath, flag, perm)
	if err != nil {
		return nil, wrapPathError("open", path, err)
	}

	return wrapFile(path, f), nil
}

func wrapFile(path string, file fs.File) fs.File {
	wrapped := &pathErrorFile{
		file: file,
		path: path,
	}

	if writer, ok := file.(interface{ Write([]byte) (int, error) }); ok {
		return &pathErrorWritableFile{
			pathErrorFile: wrapped,
			writer:        writer,
		}
	}

	return wrapped
}

type pathErrorFile struct {
	file fs.File
	path string
}

func (f *pathErrorFile) Stat() (fs.FileInfo, error) {
	info, err := f.file.Stat()
	return info, wrapPathError("stat", f.path, err)
}

func (f *pathErrorFile) Read(p []byte) (int, error) {
	n, err := f.file.Read(p)
	if err != nil {
		err = wrapPathError("read", f.path, err)
	}

	return n, err
}

func (f *pathErrorFile) Close() error {
	return wrapPathError("close", f.path, f.file.Close())
}

type pathErrorWritableFile struct {
	*pathErrorFile
	writer interface{ Write([]byte) (int, error) }
}

func (f *pathErrorWritableFile) Write(p []byte) (int, error) {
	n, err := f.writer.Write(p)
	if err != nil {
		err = wrapPathError("write", f.path, err)
	}

	return n, err
}

func wrapPathError(op, path string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, io.EOF) {
		return err
	}

	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		err = pathErr.Err
	}

	return &fs.PathError{
		Op:   op,
		Path: path,
		Err:  err,
	}
}

func matchPrefix(path, prefix string) bool {
	if prefix == "/" {
		return true
	}

	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func toSubpath(path, prefix string) string {
	if prefix == "/" {
		if path == "/" {
			return "."
		}

		return strings.TrimPrefix(path, "/")
	}

	if path == prefix {
		return "."
	}

	return strings.TrimPrefix(path, prefix+"/")
}

func normalizePath(path string) (string, error) {
	parts := strings.Split(path, "/")
	normalized := make([]string, 0, len(parts))

	for _, part := range parts {
		switch part {
		case "", ".":
			continue
		case "..":
			if len(normalized) == 0 {
				return "", fs.ErrPermission
			}

			normalized = normalized[:len(normalized)-1]
		default:
			normalized = append(normalized, part)
		}
	}

	if len(normalized) == 0 {
		return "/", nil
	}

	return "/" + strings.Join(normalized, "/"), nil
}
