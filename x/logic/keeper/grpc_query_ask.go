package keeper

import (
	goctx "context"


	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v8/x/logic/meter"
	"github.com/axone-protocol/axoned/v8/x/logic/types"
	"github.com/axone-protocol/axoned/v8/x/logic/util"
)

var defaultSolutionsLimit = sdkmath.OneUint()

func (k Keeper) Ask(ctx goctx.Context, req *types.QueryServiceAskRequest) (response *types.QueryServiceAskResponse, err error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.InvalidArgument, "request is nil")
	}

	sdkCtx := withSafeGasMeter(sdk.UnwrapSDKContext(ctx))
	defer func() {
		if r := recover(); r != nil {
			if gasError, ok := r.(storetypes.ErrorOutOfGas); ok {
				response, err = nil, errorsmod.Wrapf(
					types.LimitExceeded, "out of gas: %s <%s> (%d/%d)",
					types.ModuleName, gasError.Descriptor, sdkCtx.GasMeter().GasConsumed(), sdkCtx.GasMeter().Limit())

				return
			}

			panic(r)
		}
	}()

	limits := k.limits(ctx)
	if err := checkLimits(req, limits); err != nil {
		return nil, err
	}

	return k.execute(
		sdkCtx,
		req.Program,
		req.Query,
		util.DerefOrDefault(req.Limit, defaultSolutionsLimit))
}

// withSafeGasMeter returns a new context with a gas meter that has the given limit.
// The gas meter is go-router-safe.
func withSafeGasMeter(sdkCtx sdk.Context) sdk.Context {
	gasMeter := meter.WithSafeMeter(sdkCtx.GasMeter())

	return sdkCtx.WithGasMeter(gasMeter)
}
