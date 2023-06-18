package keeper

import (
	"github.com/okp4/okp4d/x/mint/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis new mint genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, ak types.AccountKeeper, data *types.GenesisState) {
	k.SetMinter(ctx, data.Minter)
	err := k.SetParams(ctx, data.Params)
	if err != nil {
		panic(errorsmod.Wrapf(err, "error setting params"))
	}
	ak.GetModuleAccount(ctx, types.ModuleName)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)
	return types.NewGenesisState(minter, params)
}
