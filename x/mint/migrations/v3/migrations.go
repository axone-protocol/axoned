package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oldTypes "github.com/okp4/okp4d/x/mint/migrations/v3/types"
	"github.com/okp4/okp4d/x/mint/types"
)

// MigrateStore migrates the x/mint module state from the consensus version 2 to
// version 3.
// This version include new/deleted parameters in store.
func MigrateStore(ctx sdk.Context,
	store sdk.KVStore,
	cdc codec.BinaryCodec,
) error {
	logger := ctx.Logger().
		With("module", "mint").
		With("migration", "v3")

	logger.Debug("starting module migration")

	logger.Debug("migrate mint params")

	var oldParams oldTypes.Params
	d := store.Get(types.ParamsKey)
	err := cdc.Unmarshal(d, &oldParams)
	if err != nil {
		return err
	}

	newParams := types.DefaultParams()
	newParams.MintDenom = oldParams.MintDenom

	bz, err := cdc.Marshal(&newParams)
	if err != nil {
		return err
	}
	store.Set(types.ParamsKey, bz)

	logger.Debug("module migration done")

	return nil
}
