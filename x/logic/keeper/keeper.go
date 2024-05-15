package keeper

import (
	goctx "context"
	"fmt"
	"io/fs"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v7/x/logic/types"
)

type FSProvider = func(ctx goctx.Context) fs.FS

type (
	Keeper struct {
		cdc      codec.BinaryCodec
		storeKey storetypes.StoreKey
		memKey   storetypes.StoreKey
		// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
		authority sdk.AccAddress

		authKeeper types.AccountKeeper
		bankKeeper types.BankKeeper
		fsProvider FSProvider
	}
)

func NewKeeper(
	cdc codec.BinaryCodec,
	storeKey,
	memKey storetypes.StoreKey,
	authority sdk.AccAddress,
	authKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	fsProvider FSProvider,
) *Keeper {
	// ensure gov module account is set and is not nil
	if err := sdk.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:        cdc,
		storeKey:   storeKey,
		memKey:     memKey,
		authority:  authority,
		authKeeper: authKeeper,
		bankKeeper: bankKeeper,
		fsProvider: fsProvider,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the x/logic module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority.String()
}
