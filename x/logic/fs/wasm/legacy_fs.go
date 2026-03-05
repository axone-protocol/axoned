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

	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

const (
	// Scheme is the URI scheme for the legacy CosmWasm filesystem.
	Scheme = "cosmwasm"
)

const (
	queryKey        = "query"
	base64DecodeKey = "base64Decode"
)

type legacyVFS struct {
	ctx        context.Context
	wasmKeeper types.WasmKeeper
}

var (
	_ fs.FS         = (*legacyVFS)(nil)
	_ fs.ReadFileFS = (*legacyVFS)(nil)
)

// NewLegacyFS creates a legacy URI-oriented filesystem.
// This is kept for compatibility in tests while path-based devices are migrated.
func NewLegacyFS(ctx context.Context, wasmKeeper types.WasmKeeper) fs.ReadFileFS {
	return &legacyVFS{ctx: ctx, wasmKeeper: wasmKeeper}
}

func (f *legacyVFS) Open(name string) (fs.File, error) {
	data, err := f.readFile("open", name)
	if err != nil {
		return nil, err
	}

	sdkCtx := sdk.UnwrapSDKContext(f.ctx)
	return NewVirtualFile(name, data, prolog.ResolveHeaderInfo(sdkCtx).Time), nil
}

func (f *legacyVFS) ReadFile(name string) ([]byte, error) {
	return f.readFile("open", name)
}

func (f *legacyVFS) readFile(op string, name string) ([]byte, error) {
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

func (f *legacyVFS) parsePath(op string, path string) (sdk.AccAddress, string, bool, error) {
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
			&fs.PathError{Op: op, Path: path, Err: fmt.Errorf("empty path given, should be '%s:{contractName}:{contractAddr}?query={query}'",
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
