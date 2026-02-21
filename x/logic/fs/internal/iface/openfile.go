package iface

import "io/fs"

// OpenFileFS is implemented by file systems that support opening files with flags.
type OpenFileFS interface {
	fs.FS
	OpenFile(name string, flag int, perm fs.FileMode) (fs.File, error)
}
