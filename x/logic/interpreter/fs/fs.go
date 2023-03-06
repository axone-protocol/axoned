package fs

import (
	goctx "context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/logic/types"
)

type FileSystem struct {
	ctx        goctx.Context
	wasmKeeper types.WasmKeeper
}

// New return a new FileSystem object that will handle all virtual file on the interpreter.
// File can be provided from different sources like CosmWasm cw-storage smart contract.
func New(ctx goctx.Context, keeper types.WasmKeeper) FileSystem {
	return FileSystem{
		ctx:        ctx,
		wasmKeeper: keeper,
	}
}

func (f FileSystem) Open(name string) (fs.File, error) {
	data, err := f.ReadFile(name)
	return Object(data), err
}

// ReadFile reads the named file and returns its contents.
// A successful call returns a nil error, not io.EOF.
// (Because ReadFi   le reads the whole file, the expected EOF
// from the final Read is not treated as an error to be reported.)
//
// The caller is permitted to modify the returned byte slice.
// This method should return a copy of the underlying data.
func (f FileSystem) ReadFile(name string) ([]byte, error) {
	sdkCtx := sdk.UnwrapSDKContext(f.ctx)

	req := []byte(fmt.Sprintf("{\"object_data\":{\"id\": \"%s\"}}", name))
	contractAddr, err := sdk.AccAddressFromBech32("okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht")
	if err != nil {
		return nil, err
	}

	data, err := f.wasmKeeper.QuerySmart(sdkCtx, contractAddr, req)
	if err != nil {
		return nil, err
	}
	var program string
	err = json.Unmarshal(data, &program)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(program)
	return decoded, err
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
	//TODO implement me
	panic("implement me")
}

func (o Object) Close() error {
	return nil
}
