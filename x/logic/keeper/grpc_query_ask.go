package keeper

import (
	goctx "context"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v11/x/logic/meter"
	"github.com/axone-protocol/axoned/v11/x/logic/types"
)

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

	params := k.GetParams(sdkCtx)
	if err := checkLimits(req, params.Limits); err != nil {
		return nil, err
	}

	return k.execute(
		sdkCtx,
		params,
		req.Program,
		req.Query,
		req.Limit)
}

func checkLimits(request *types.QueryServiceAskRequest, limits types.Limits) error {
	size := uint64(len(request.GetQuery()))

	if limits.MaxSize != 0 && size > limits.MaxSize {
		return errorsmod.Wrapf(types.LimitExceeded, "query: %d > MaxSize: %d", size, limits.MaxSize)
	}

	return nil
}

// withSafeGasMeter returns a new context with a gas meter that has the given limit.
// The gas meter is go-router-safe.
func withSafeGasMeter(sdkCtx sdk.Context) sdk.Context {
	gasMeter := meter.WithSafeMeter(sdkCtx.GasMeter())

	return sdkCtx.WithGasMeter(gasMeter)
}
