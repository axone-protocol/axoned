package fs

import (
	"fmt"
	"io/fs"

	"github.com/okp4/okp4d/x/logic/types"
)

type FileSystem struct {
	wasmKeeper types.WasmKeeper
}

// New return a new FileSystem object that will handle all virtual file on the interpreter.
// File can be provided from different sources like CosmWasm cw-storage smart contract.
func New(keeper types.WasmKeeper) FileSystem {
	return FileSystem{
		wasmKeeper: keeper,
	}
}

func (f FileSystem) Open(name string) (fs.File, error) {
	return nil, fmt.Errorf("not implemented")
}
