package prolog

import (
	"context"

	"github.com/ichiban/prolog/engine"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
