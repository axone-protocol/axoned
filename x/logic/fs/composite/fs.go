package composite

import (
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"sort"

	"golang.org/x/exp/maps"
)

type CopositeFS interface {
	fs.FS

	// Mount mounts a filesystem to the given mount point.
	// The mount point is the scheme of the URI.
	Mount(mountPoint string, fs fs.FS)
	// ListMounts returns a list of all mount points in sorted order.
	ListMounts() []string
}

type vfs struct {
	mounted map[string]fs.FS
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

func NewFS() CopositeFS {
	return &vfs{mounted: make(map[string]fs.FS)}
}

func (f *vfs) Open(name string) (fs.File, error) {
	uri, err := f.validatePath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	vfs, err := f.resolve(uri)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrNotExist}
	}
	return vfs.Open(name)
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	uri, err := f.validatePath(name)
	if err != nil {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: fs.ErrInvalid}
	}

	vfs, err := f.resolve(uri)
	if err != nil {
		return nil, &fs.PathError{Op: "readfile", Path: name, Err: fs.ErrNotExist}
	}

	if vfs, ok := vfs.(fs.ReadFileFS); ok {
		return vfs.ReadFile(name)
	}

	file, err := vfs.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (f *vfs) Mount(mountPoint string, fs fs.FS) {
	f.mounted[mountPoint] = fs
}

func (f *vfs) ListMounts() []string {
	mounts := maps.Keys(f.mounted)
	sort.Strings(mounts)

	return mounts
}

// validatePath checks if the provided path is a valid URL and returns its URI.
func (f *vfs) validatePath(name string) (*url.URL, error) {
	uri, err := url.Parse(name)
	if err != nil {
		return nil, err
	}
	if uri.Scheme == "" {
		return nil, fmt.Errorf("missing scheme in path: %s", name)
	}
	return uri, nil
}

// resolve returns the filesystem mounted at the given URI scheme.
func (f *vfs) resolve(uri *url.URL) (fs.FS, error) {
	vfs, ok := f.mounted[uri.Scheme]
	if !ok {
		return nil, fmt.Errorf("no filesystem mounted at: %s", uri.Scheme)
	}
	return vfs, nil
}
