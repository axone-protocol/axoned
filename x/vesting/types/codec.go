package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/axone-protocol/axoned/v12/x/vesting/exported"
)

// RegisterLegacyAminoCodec registers the vesting interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterInterface((*exported.VestingAccount)(nil), nil)
	cdc.RegisterConcrete(&BaseVestingAccount{}, "cosmos-sdk/BaseVestingAccount", nil)
	cdc.RegisterConcrete(&ContinuousVestingAccount{}, "cosmos-sdk/ContinuousVestingAccount", nil)
	cdc.RegisterConcrete(&DelayedVestingAccount{}, "cosmos-sdk/DelayedVestingAccount", nil)
	cdc.RegisterConcrete(&PeriodicVestingAccount{}, "cosmos-sdk/PeriodicVestingAccount", nil)
	cdc.RegisterConcrete(&PermanentLockedAccount{}, "cosmos-sdk/PermanentLockedAccount", nil)
	cdc.RegisterConcrete(&CliffVestingAccount{}, "cosmos-sdk/CliffVestingAccount", nil)
	legacy.RegisterAminoMsg(cdc, &MsgCreateVestingAccount{}, "cosmos-sdk/MsgCreateVestingAccount")
	legacy.RegisterAminoMsg(cdc, &MsgCreatePermanentLockedAccount{}, "cosmos-sdk/MsgCreatePermLockedAccount")
	legacy.RegisterAminoMsg(cdc, &MsgCreatePeriodicVestingAccount{}, "cosmos-sdk/MsgCreatePeriodVestAccount")
	legacy.RegisterAminoMsg(cdc, &MsgCreateCliffVestingAccount{}, "cosmos-sdk/MsgCreateCliffVestingAccount")
}

// RegisterInterfaces associates protoName with AccountI and VestingAccount
// Interfaces and creates a registry of it's concrete implementations.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterInterface(
		"vesting.v1beta1.VestingAccount",
		(*exported.VestingAccount)(nil),
		&ContinuousVestingAccount{},
		&DelayedVestingAccount{},
		&PeriodicVestingAccount{},
		&PermanentLockedAccount{},
		&CliffVestingAccount{},
	)

	registry.RegisterImplementations(
		(*sdk.AccountI)(nil),
		&BaseVestingAccount{},
		&DelayedVestingAccount{},
		&ContinuousVestingAccount{},
		&PeriodicVestingAccount{},
		&PermanentLockedAccount{},
		&CliffVestingAccount{},
	)

	registry.RegisterImplementations(
		(*authtypes.GenesisAccount)(nil),
		&BaseVestingAccount{},
		&DelayedVestingAccount{},
		&ContinuousVestingAccount{},
		&PeriodicVestingAccount{},
		&PermanentLockedAccount{},
		&CliffVestingAccount{},
	)

	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateVestingAccount{},
		&MsgCreatePermanentLockedAccount{},
		&MsgCreateCliffVestingAccount{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}
