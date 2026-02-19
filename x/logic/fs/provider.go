package fs

import (
	goctx "context"
	"io/fs"
)

// Provider is a function that returns a filesystem and an error.
// It is used to provide the filesystem to the logic module when executing logic.
type Provider = func(ctx goctx.Context) (fs.FS, error)
