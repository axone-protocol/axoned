package keeper

import (
	goctx "context"
	"math"

	errorsmod "cosmossdk.io/errors"
	storetypes "cosmossdk.io/store/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/meter"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

func (k Keeper) Ask(ctx goctx.Context, req *types.QueryAskRequest) (response *types.QueryAskResponse, err error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidArgument, "request is nil")
	}

	sdkCtx := withSafeGasMeter(sdk.UnwrapSDKContext(ctx))
	defer func() {
		if r := recover(); r != nil {
			if gasError, ok := r.(storetypes.ErrorOutOfGas); ok {
				response, err = nil, errorsmod.Wrapf(
					types.ErrLimitExceeded, "out of gas: %s <%s> (%d/%d)",
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
	consumeRequestIOGas(sdkCtx.GasMeter(), req, params.GasPolicy.IoCoeff)

	return k.execute(
		sdkCtx,
		params,
		req.Program,
		req.Query,
		req.Limit)
}

func checkLimits(request *types.QueryAskRequest, limits types.Limits) error {
	size := sourceSize(request)

	if limits.MaxSize != 0 && size > limits.MaxSize {
		return errorsmod.Wrapf(types.ErrLimitExceeded, "source: %d > MaxSize: %d", size, limits.MaxSize)
	}

	return nil
}

func consumeRequestIOGas(gasMeter storetypes.GasMeter, request *types.QueryAskRequest, coeff uint64) {
	consumeIOGas(gasMeter, sourceSize(request), coeff)
}

func sourceSize(request *types.QueryAskRequest) uint64 {
	return uint64(len(request.GetProgram()) + len(request.GetQuery()))
}

func consumeIOGas(gasMeter storetypes.GasMeter, units, coeff uint64) {
	if units == 0 {
		return
	}
	if coeff == 0 {
		coeff = 1
	}

	consumed, overflow := meterMultiplyUint64Overflow(units, coeff)
	if overflow {
		gasMeter.ConsumeGas(math.MaxUint64, "IO")
		return
	}

	gasMeter.ConsumeGas(consumed, "IO")
}

func meterMultiplyUint64Overflow(a, b uint64) (uint64, bool) {
	if a == 0 || b == 0 {
		return 0, false
	}

	c := a * b
	if c/a != b || c/b != a {
		return 0, true
	}

	return c, false
}

// withSafeGasMeter returns a new context with a gas meter that has the given limit.
// The gas meter is go-router-safe.
func withSafeGasMeter(sdkCtx sdk.Context) sdk.Context {
	gasMeter := meter.WithSafeMeter(sdkCtx.GasMeter())

	return sdkCtx.WithGasMeter(gasMeter)
}
