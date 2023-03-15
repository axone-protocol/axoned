package fs

import (
	goctx "context"
	"io/fs"
)

// FileSystem is the custom virtual file system used into the blockchain.
// It will hold a list of handler that can resolve file URI and return the corresponding binary file.
type FileSystem struct {
	ctx    goctx.Context
	router Router
}

// New return a new FileSystem object that will handle all virtual file on the interpreter.
// File can be provided from different sources like CosmWasm cw-storage smart contract.
func New(ctx goctx.Context, handlers []URIHandler) FileSystem {
	router := NewRouter()
	for _, handler := range handlers {
		router.RegisterHandler(handler)
	}
	return FileSystem{
		ctx:    ctx,
		router: router,
	}
}

// Open will read the entire file from ReadFile interface,
// Since file is provided by a provider that do not support streams.
func (f FileSystem) Open(name string) (fs.File, error) {
	//data, err := f.ReadFile(name)
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

// ReadFile read the entire file at the uri provided.
// Parse all handler and return the first supported handler file response.
//func (f FileSystem) ReadFile(name string) ([]byte, error) {
//	return f.router.Open(f.ctx, name)
//}
