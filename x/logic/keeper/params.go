package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/logic/types"
)

// GetParams get all parameters as types.Params.
func (k Keeper) GetParams(ctx sdk.Context) types.Params {
	// TODO implement me
	return types.Params{}
}

// SetParams set the params.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramstore.SetParamSet(ctx, &params)
}
