package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v8/x/mint/types"
)

// InitGenesis new mint genesis.
func (k Keeper) InitGenesis(ctx sdk.Context, ak types.AccountKeeper, data *types.GenesisState) {
	if err := k.Minter.Set(ctx, data.Minter); err != nil {
		panic(err)
	}

	if err := k.Params.Set(ctx, data.Params); err != nil {
		panic(err)
	}

	ak.GetModuleAccount(ctx, types.ModuleName)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func (k Keeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	minter, err := k.Minter.Get(ctx)
	if err != nil {
		panic(err)
	}

	params, err := k.Params.Get(ctx)
	if err != nil {
		panic(err)
	}

	return types.NewGenesisState(minter, params)
}
