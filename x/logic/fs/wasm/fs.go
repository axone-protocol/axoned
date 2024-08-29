package wasm

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/url"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v10/x/logic/types"
)

const (
	// scheme is the URI scheme for the WASM filesystem.
	Scheme = "cosmwasm"
)

const (
	queryKey        = "query"
	base64DecodeKey = "base64Decode"
)

type vfs struct {
	ctx        context.Context
	wasmKeeper types.WasmKeeper
}

var (
	_ fs.FS         = (*vfs)(nil)
	_ fs.ReadFileFS = (*vfs)(nil)
)

// NewFS creates a new filesystem that can read data from a WASM contract.
// The URI should be in the format `cosmwasm:{contractName}:{contractAddr}?query={query}`.
func NewFS(ctx context.Context, wasmKeeper types.WasmKeeper) fs.ReadFileFS {
	return &vfs{ctx: ctx, wasmKeeper: wasmKeeper}
}

func (f *vfs) Open(name string) (fs.File, error) {
	data, err := f.readFile("open", name)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return NewVirtualFile(name, data, sdkCtx.BlockTime()), nil
}

func (f *vfs) ReadFile(name string) ([]byte, error) {
	return f.readFile("readfile", name)
}

func (f *vfs) readFile(op string, name string) ([]byte, error) {
	contractAddr, query, base64Decode, err := f.parsePath(op, name)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	data, err := f.wasmKeeper.QuerySmart(sdkCtx, contractAddr, []byte(query))
	if err != nil {
		return nil, &fs.PathError{
			Op:   op,
			Path: name,
			Err:  fmt.Errorf("failed to query WASM contract %s: %w", contractAddr, err),
		}
	}

	if base64Decode {
		var program string
		err = json.Unmarshal(data, &program)
		if err != nil {
			return nil, &fs.PathError{
				Op:   op,
				Path: name,
				Err:  fmt.Errorf("failed to unmarshal JSON WASM response to string: %w", err),
			}
		}

		data, err = base64.StdEncoding.DecodeString(program)
		if err != nil {
			return nil, &fs.PathError{
				Op:   op,
				Path: name,
				Err:  fmt.Errorf("failed to decode WASM base64 response: %w", err),
			}
		}
	}

	return data, nil
}

// parsePath parses the provided path and returns its component.
func (f *vfs) parsePath(op string, path string) (sdk.AccAddress, string, bool, error) {
	uri, err := url.Parse(path)
	if err != nil {
		return nil, "", false,
			&fs.PathError{Op: op, Path: path, Err: fs.ErrInvalid}
	}

	if uri.Scheme != Scheme {
		return nil, "", false,
			&fs.PathError{Op: op, Path: path, Err: fmt.Errorf("invalid scheme, expected '%s', got '%s'", Scheme, uri.Scheme)}
	}

	paths := strings.SplitAfter(uri.Opaque, ":")
	pathsLen := len(paths)
	if pathsLen < 1 || paths[pathsLen-1] == "" {
		return nil, "", false,
			&fs.PathError{Op: op, Path: path, Err: fmt.Errorf("emtpy path given, should be '%s:{contractName}:{contractAddr}?query={query}'",
				Scheme)}
	}

	lastPart := paths[len(paths)-1]
	contractAddr, err := sdk.AccAddressFromBech32(lastPart)
	if err != nil {
		return nil, "", false,
			&fs.PathError{
				Op:   op,
				Path: path,
				Err:  fmt.Errorf("failed to convert path '%s' to contract address: %w", lastPart, err),
			}
	}

	query := uri.Query().Get(queryKey)
	if query == "" {
		return nil, "", false,
			&fs.PathError{
				Op:   op,
				Path: path,
				Err:  fmt.Errorf("uri should contains `query` params"),
			}
	}

	base64Decode := true
	if uri.Query().Has(base64DecodeKey) {
		if base64Decode, err = strconv.ParseBool(uri.Query().Get(base64DecodeKey)); err != nil {
			return nil, "", false,
				&fs.PathError{
					Op:   op,
					Path: path,
					Err:  fmt.Errorf("failed to convert 'base64Decode' query value to boolean: %w", err),
				}
		}
	}

	return contractAddr, query, base64Decode, nil
}
