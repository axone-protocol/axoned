package prolog

import (
	"bytes"
	"context"

	"github.com/axone-protocol/prolog/v3/engine"

	coreheader "cosmossdk.io/core/header"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

// UnwrapSDKContext retrieves a Context from a context.Context instance
// attached with WrapSDKContext.
func UnwrapSDKContext(ctx context.Context, env *engine.Env) (sdk.Context, error) {
	if sdkCtx, ok := ctx.(sdk.Context); ok {
		return sdkCtx, nil
	}
	if sdkCtx, ok := ctx.Value(sdk.SdkContextKey).(sdk.Context); ok {
		return sdkCtx, nil
	}

	return sdk.Context{}, engine.ResourceError(ResourceContext(), env)
}

// ContextValue returns the value associated with this key in the context.
// If the value is not found, it returns the error: error(resource_error(resource_context(<key>))).
func ContextValue[V any](ctx context.Context, key types.ContextKey, env *engine.Env) (V, error) {
	if value, ok := ctx.Value(key).(V); ok {
		return value, nil
	}

	var zero V
	return zero, engine.ResourceError(ResourceContextValue(string(key)), env)
}

// HeaderInfo retrieves header information from SDK context.
// It reads sdk.Context.HeaderInfo() and returns defensive copies for byte slices.
func HeaderInfo(ctx context.Context, env *engine.Env) (coreheader.Info, error) {
	sdkCtx, err := UnwrapSDKContext(ctx, env)
	if err != nil {
		return coreheader.Info{}, err
	}

	return ResolveHeaderInfo(sdkCtx), nil
}

// ResolveHeaderInfo resolves canonical header information from the given SDK context.
// It exists to centralize access to sdk.Context.HeaderInfo() across the logic module.
func ResolveHeaderInfo(sdkCtx sdk.Context) coreheader.Info {
	info := sdkCtx.HeaderInfo()
	info.Hash = bytes.Clone(info.Hash)
	info.AppHash = bytes.Clone(info.AppHash)

	return info
}
