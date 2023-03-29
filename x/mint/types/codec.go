package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

var (
	amino = codec.NewLegacyAmino()
	// ModuleCdc references the global evm module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding.
	ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

	// AminoCdc is a amino codec created to support amino JSON compatible msgs.
	AminoCdc = codec.NewAminoCodec(amino)
)

func init() {
	cryptocodec.RegisterCrypto(amino)
	RegisterLegacyAminoCodec(amino)
	amino.Seal()
}

const (
	// Amino names.
	updateParamsName = "okp4/mint/MsgUpdateParams"
)

// RegisterInterfaces register implementations.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgUpdateParams{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_MsgService_serviceDesc)
}

// RegisterLegacyAminoCodec required for EIP-712.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgUpdateParams{}, updateParamsName, nil)
}
