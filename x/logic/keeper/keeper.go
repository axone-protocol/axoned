package keeper

import (
	goctx "context"
	"fmt"
	"io/fs"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/tendermint/tendermint/libs/log"
)

type FSProvider = func(ctx goctx.Context) fs.FS

type (
	Keeper struct {
		cdc        codec.BinaryCodec
		storeKey   storetypes.StoreKey
		memKey     storetypes.StoreKey
		paramstore paramtypes.Subspace

		authKeeper types.AccountKeeper
		bankKeeper types.BankKeeper
		fsProvider FSProvider
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	ps paramtypes.Subspace,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	fsProvider FSProvider,
) *Keeper {
	// set KeyTable if it has not already been set
	if !ps.HasKeyTable() {
		ps = ps.WithKeyTable(types.ParamKeyTable())
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		paramstore: ps,
		authKeeper: authKeeper,
		bankKeeper: bankKeeper,
		fsProvider: fsProvider,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
