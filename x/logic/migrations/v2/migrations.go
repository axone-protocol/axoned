package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/okp4/okp4d/x/logic/types"
)

func MigrateStore(ctx sdk.Context, paramstore paramtypes.Subspace, cdc codec.BinaryCodec) error {
	logger := ctx.Logger().With("migration", "logic")

	logger.Debug("start module migration")

	// Add default params keys / values
	logger.Debug("set params default values")
	params := types.DefaultParams()
	paramstore.SetParamSet(ctx, &params)

	return nil
}
