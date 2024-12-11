package filtered

import (
	"io/fs"
	"net/url"

	"github.com/axone-protocol/axoned/v11/x/logic/util"
)

type vfs struct {
	fs        fs.FS
	whitelist []*url.URL
	blacklist []*url.URL
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates a new filtered filesystem that wraps the provided filesystem.
// The whitelist and blacklist are used to filter the paths that can be accessed.
func NewFS(underlyingFS fs.FS, whitelist, blacklist []*url.URL) fs.ReadFileFS {
	return &vfs{fs: underlyingFS, whitelist: whitelist, blacklist: blacklist}
}

func (f *vfs) Open(name string) (fs.File, error) {
	if err := f.accept("open", name); err != nil {
		return nil, err
	}
	return f.fs.Open(name)
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	if err := f.accept("open", name); err != nil {
		return nil, err
	}

	if vfs, ok := f.fs.(fs.ReadFileFS); ok {
		return vfs.ReadFile(name)
	}

	return nil, &fs.PathError{Op: "readfile", Path: name, Err: fs.ErrInvalid}
}

// validatePath checks if the provided path is a valid URL.
func (f *vfs) validatePath(name string) (*url.URL, error) {
	uri, err := url.Parse(name)
	if err != nil {
		return nil, err
	}

	return uri, nil
}

// accept checks if the provided path is allowed by the whitelist and blacklist.
// If the path is allowed, it returns nil; otherwise, it returns an error.
func (f *vfs) accept(op string, name string) error {
	uri, err := f.validatePath(name)
	if err != nil {
		return &fs.PathError{Op: op, Path: name, Err: fs.ErrInvalid}
	}

	if !util.WhitelistBlacklistMatches(f.whitelist, f.blacklist, util.URLMatches)(uri) {
		return &fs.PathError{
			Op:   op,
			Path: name,
			Err:  fs.ErrPermission,
		}
	}

	return nil
}
