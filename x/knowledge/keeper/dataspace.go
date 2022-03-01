package keeper

import (
    "github.com/cosmos/cosmos-sdk/store/prefix"
    sdk "github.com/cosmos/cosmos-sdk/types"
    "github.com/okp4/okp4d/x/knowledge/types"
)

func (k Keeper) SaveDataspace(
    ctx sdk.Context,
    id string,
    name string) (*types.MsgBangDataspaceResponse, error) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataspaceKeyPrefix)
    key := types.GetDataspaceKey(id)

    store.Set(key, []byte(name))

    return &types.MsgBangDataspaceResponse{}, nil
}

func (k Keeper) HasDataspace(
    ctx sdk.Context,
    id string,
) bool {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DataspaceKeyPrefix)
    key := types.GetDataspaceKey(id)

    return store.Has(key)
}
