package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/mint/exported"
	"github.com/okp4/okp4d/x/mint/types"
)

// MigrateStore migrates the x/mint module state from the consensus version 1 to
// version 2.
// Specifically, it takes the parameters that are currently stored
// and managed by the x/params modules and stores them directly into the x/mint
// module state.
func MigrateStore(ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	legacySubspace exported.Subspace,
) error {
	logger := ctx.Logger().
		With("module", "mint").
		With("migration", "v2")

	logger.Debug("starting module migration")

	logger.Debug("migrate mint params")

	store := ctx.KVStore(storeKey)

	var params types.Params
	legacySubspace.WithKeyTable(types.ParamKeyTable()).
		GetParamSet(ctx, &params)

	if err := params.Validate(); err != nil {
		return err
	}

	bz, err := cdc.Marshal(&params)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	logger.Debug("module migration done")

	return nil
}
