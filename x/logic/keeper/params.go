package keeper

import (
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/okp4/okp4d/x/logic/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
    return types.NewParams(
        k.Foo(ctx),
    )
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
    k.paramstore.SetParamSet(ctx, &params)
}

// Foo returns the Foo param
func (k Keeper) Foo(ctx sdk.Context) (res string) {
    k.paramstore.Get(ctx, types.KeyFoo, &res)
    return
}
