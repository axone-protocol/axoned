package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/logic/exported"
	"github.com/okp4/okp4d/x/logic/types"
)

func MigrateStore(ctx sdk.Context, legacySubspace exported.Subspace, cdc codec.BinaryCodec) error {
	logger := ctx.Logger().
		With("module", "logic").
		With("migration", "v2")

	logger.Debug("starting module migration")

	// Add default params keys / values
	logger.Debug("set params default values")
	params := types.DefaultParams()
	legacySubspace.SetParamSet(ctx, &params)

	logger.Debug("module migration done")

	return nil
}
