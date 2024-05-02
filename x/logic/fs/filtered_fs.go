package fs

import (
	"io/fs"
	"net/url"

	"github.com/axone/axoned/v7/x/logic/util"
)

// FilteredFS is a wrapper around a fs.FS that filters out files that are not allowed to be read.
// This is used to prevent the interpreter from reading files, using protocols that are not allowed to be used
// by the interpreter on the blockchain.
// The whitelist and blacklist are mutually exclusive. If both are set, the blacklist will be ignored.
type FilteredFS struct {
	decorated fs.FS
	whitelist []*url.URL
	blacklist []*url.URL
}

var _ fs.FS = (*FilteredFS)(nil)

// NewFilteredFS returns a new FilteredFS object that will filter out files that are not allowed to be read
// according to the whitelist and blacklist parameters.
func NewFilteredFS(whitelist, blacklist []*url.URL, decorated fs.FS) *FilteredFS {
	return &FilteredFS{
		decorated: decorated,
		whitelist: whitelist,
		blacklist: blacklist,
	}
}

// Open opens the named file.
// The name parameter is a URL that will be parsed and checked against the whitelist and blacklist configured.
func (f *FilteredFS) Open(name string) (fs.File, error) {
	urlFile, err := url.Parse(name)
	if err != nil {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  err,
		}
	}

	if !util.WhitelistBlacklistMatches(f.whitelist, f.blacklist, util.URLMatches)(urlFile) {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  fs.ErrPermission,
		}
	}

	return f.decorated.Open(name)
}
