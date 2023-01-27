package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/logic/types"
)

// memoizedParams is a minimalistic cache for the params.
// should be used only in go-routine-safe contexts.
var memoizedParams = struct {
	v     types.Params
	isSet bool
}{
	v: types.DefaultParams(),
}

// GetParams get all parameters as types.Params.
// The returned value is memoized so that it can be used without querying the store, and consuming gas.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	if !memoizedParams.isSet {
		k.paramstore.GetParamSet(ctx, &memoizedParams.v)
		memoizedParams.isSet = true
	}

	return memoizedParams.v
}

// SetParams set the cachedParams.
// The params are also memoized so that they can be used without querying the store again.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)

	memoizedParams.v = params
	memoizedParams.isSet = true
}
