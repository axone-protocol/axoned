package keeper

import (
	goctx "context"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/logic/meter"
	"github.com/okp4/okp4d/x/logic/types"
)

func (k Keeper) Ask(ctx goctx.Context, req *types.QueryServiceAskRequest) (response *types.QueryServiceAskResponse, err error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if req == nil {
		return nil, sdkerrors.Wrap(types.InvalidArgument, "request is nil")
	}

	limits := k.limits(ctx)
	if err := checkLimits(req, limits); err != nil {
		return nil, err
	}

	sdkCtx = withGasMeter(sdkCtx, limits)
	defer func() {
		if r := recover(); r != nil {
			if gasError, ok := r.(sdk.ErrorOutOfGas); ok {
				response, err = nil, sdkerrors.Wrapf(
					types.LimitExceeded, "out of gas: %s <%s> (%d/%d)",
					types.ModuleName, gasError.Descriptor, sdkCtx.GasMeter().GasConsumed(), sdkCtx.GasMeter().Limit())

				return
			}

			panic(r)
		}
		if sdkCtx.GasMeter().IsOutOfGas() {
			response, err = nil, sdkerrors.Wrapf(
				types.LimitExceeded, "out of gas: %s (%d/%d)",
				types.ModuleName, sdkCtx.GasMeter().GasConsumed(), sdkCtx.GasMeter().Limit())

			return
		}
	}()
	sdkCtx.GasMeter().ConsumeGas(sdkCtx.GasMeter().GasConsumed(), types.ModuleName)

	//nolint:contextcheck
	return k.execute(
		sdkCtx,
		req.Program,
		req.Query)
}

// withGasMeter returns a new context with a gas meter that has the given limit.
// The gas meter is go-router-safe.
func withGasMeter(sdkCtx sdk.Context, limits types.Limits) sdk.Context {
	gasMeter := meter.WithSafeMeter(sdk.NewGasMeter(limits.MaxGas.Uint64()))

	return sdkCtx.WithGasMeter(gasMeter)
}
