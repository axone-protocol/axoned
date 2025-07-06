package keeper

import (
	"fmt"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v12/x/logic/fs"
	"github.com/axone-protocol/axoned/v12/x/logic/types"
)

type (
	Keeper struct {
		cdc               codec.BinaryCodec
		interfaceRegistry cdctypes.InterfaceRegistry
		storeKey          storetypes.StoreKey
		memKey            storetypes.StoreKey
		// the address capable of executing a MsgUpdateParams message. Typically, this should be the x/gov module account.
		authority sdk.AccAddress

		authKeeper       types.AccountKeeper
		authQueryService types.AuthQueryService
		bankKeeper       types.BankKeeper
		fsProvider       fs.Provider
	}
)

func NewKeeper(cdc codec.BinaryCodec, interfaceRegistry cdctypes.InterfaceRegistry, storeKey, memKey storetypes.StoreKey,
	authority sdk.AccAddress, authKeeper types.AccountKeeper, authQueryService types.AuthQueryService, bankKeeper types.BankKeeper,
	fsProvider fs.Provider,
) *Keeper {
	// ensure gov module account is set and is not nil
	if err := sdk.VerifyAddressFormat(authority); err != nil {
		panic(err)
	}

	return &Keeper{
		cdc:               cdc,
		interfaceRegistry: interfaceRegistry,
		storeKey:          storeKey,
		memKey:            memKey,
		authority:         authority,
		authKeeper:        authKeeper,
		authQueryService:  authQueryService,
		bankKeeper:        bankKeeper,
		fsProvider:        fsProvider,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetAuthority returns the x/logic module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority.String()
}
