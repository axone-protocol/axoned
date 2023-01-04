package predicate

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UnwrapSDKContext retrieves a Context from a context.Context instance
// attached with WrapSDKContext. It panics if a Context was not properly
// attached.
func UnwrapSDKContext(ctx context.Context) (sdk.Context, error) {
	if sdkCtx, ok := ctx.(sdk.Context); ok {
		return sdkCtx, nil
	}
	if sdkCtx, ok := ctx.Value(sdk.SdkContextKey).(sdk.Context); ok {
		return sdkCtx, nil
	}

	return sdk.Context{}, fmt.Errorf("no sdk.Context found in context")
}
