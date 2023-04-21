package fs

import (
	goctx "context"
	"io/fs"
)

// VirtualFS is the custom virtual file system used into the blockchain.
// It will hold a list of handler that can resolve file URI and return the corresponding binary file.
type VirtualFS struct {
	ctx    goctx.Context
	router Router
}

var _ fs.FS = (*VirtualFS)(nil)

// NewVirtualFS return a new VirtualFS object that will handle all virtual file on the interpreter.
// File can be provided from different sources like CosmWasm cw-storage smart contract.
func NewVirtualFS(ctx goctx.Context, handlers []URIHandler) *VirtualFS {
	router := NewRouter()
	for _, handler := range handlers {
		router.RegisterHandler(handler)
	}
	return &VirtualFS{
		ctx:    ctx,
		router: router,
	}
}

// Open opens the named file.
//
// When Open returns an error, it should be of type *PathError
// with the Op field set to "open", the Path field set to name,
// and the Err field describing the problem.
//
// Open should reject attempts to open names that do not satisfy
// ValidPath(name), returning a *PathError with Err set to
// ErrInvalid or ErrNotExist.
func (f *VirtualFS) Open(name string) (fs.File, error) {
	data, err := f.router.Open(f.ctx, name)
	if err != nil {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  err,
		}
	}
	return data, nil
}
