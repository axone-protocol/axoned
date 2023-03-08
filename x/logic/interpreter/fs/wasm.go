package fs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/logic/types"
)

const queryKey = "query"
const hostName = "wasm"

type WasmFS struct {
	wasmKeeper types.WasmKeeper
}

func NewWasmFS(keeper types.WasmKeeper) WasmFS {
	return WasmFS{wasmKeeper: keeper}
}

func (w WasmFS) CanOpen(ctx context.Context, uri *url.URL) bool {
	return uri.Host == hostName
}

func (w WasmFS) Open(ctx context.Context, uri *url.URL) ([]byte, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	paths := strings.SplitAfter(uri.Path, "/")
	if len(paths) != 2 {
		return nil, fmt.Errorf("incorect path, should contains only contract address : '://wasm/{contractAddr}?query={query}'")
	}

	contractAddr, err := sdk.AccAddressFromBech32(paths[1])
	if err != nil {
		return nil, fmt.Errorf("failed convert path '%s' to contract address: %w", uri.Path, err)
	}

	if !uri.Query().Has(queryKey) {
		return nil, fmt.Errorf("uri should contains query params")
	}
	query := uri.Query().Get(queryKey)

	data, err := w.wasmKeeper.QuerySmart(sdkCtx, contractAddr, []byte(query))
	if err != nil {
		return nil, fmt.Errorf("failed query wasm keeper: %w", err)
	}

	var program string
	err = json.Unmarshal(data, &program)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal json wasm response to string: %w", err)
	}

	decoded, err := base64.StdEncoding.DecodeString(program)
	if err != nil {
		return nil, fmt.Errorf("failed decode wasm base64 respone: %w", err)
	}
	return decoded, nil
}
