package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/mint/exported"
)

// MigrateStore migrates the x/mint module state from the consensus version 2 to
// version 3.
func MigrateStore(ctx sdk.Context,
	storeKey storetypes.StoreKey,
	cdc codec.BinaryCodec,
	legacySubspace exported.Subspace,
) error {
	logger := ctx.Logger().
		With("module", "mint").
		With("migration", "v3")

	logger.Debug("starting module migration")

	logger.Debug("migrate mint params")

	// TODO:

	logger.Debug("module migration done")

	return nil
}
