package fs

import (
	goctx "context"
	"io/fs"
)

// Provider is a function that returns a filesystem.
type Provider = func(ctx goctx.Context) fs.FS
