package fs

import (
	goctx "context"
	"io"
	"io/fs"
	"time"
)

// FileSystem is the custom virtual file system used into the blockchain.
// It will hold a list of handler that can resolve file URI and return the corresponding binary file.
type FileSystem struct {
	ctx    goctx.Context
	parser Parser
}

// New return a new FileSystem object that will handle all virtual file on the interpreter.
// File can be provided from different sources like CosmWasm cw-storage smart contract.
func New(ctx goctx.Context, handlers []URIHandler) FileSystem {
	return FileSystem{
		ctx:    ctx,
		parser: Parser{handlers},
	}
}

// Open will read the entire file from ReadFile interface,
// Since file is provided by a provider that do not support streams.
func (f FileSystem) Open(name string) (fs.File, error) {
	data, err := f.ReadFile(name)
	if err != nil {
		return nil, &fs.PathError{
			Op:   "open",
			Path: name,
			Err:  err,
		}
	}
	return Object(data), nil
}

// ReadFile read the entire file at the uri provided.
// Parse all handler and return the first supported handler file response.
func (f FileSystem) ReadFile(name string) ([]byte, error) {
	return f.parser.Parse(f.ctx, name)
}

type Object []byte

type ObjectInfo struct {
	name string
	size int64
}

func From(object Object) ObjectInfo {
	return ObjectInfo{
		name: "contract",
		size: int64(len(object)),
	}
}

func (o ObjectInfo) Name() string {
	return o.name
}

func (o ObjectInfo) Size() int64 {
	return o.size
}

func (o ObjectInfo) Mode() fs.FileMode {
	return fs.ModeIrregular
}

func (o ObjectInfo) ModTime() time.Time {
	return time.Now()
}

func (o ObjectInfo) IsDir() bool {
	return false
}

func (o ObjectInfo) Sys() any {
	return nil
}

func (o Object) Stat() (fs.FileInfo, error) {
	return From(o), nil
}

func (o Object) Read(bytes []byte) (int, error) {
	copy(bytes, o)
	return 0, io.EOF
}

func (o Object) Close() error {
	return nil
}
