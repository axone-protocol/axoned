package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBangDataspace = "bang_dataspace"

var _ sdk.Msg = &MsgBangDataspace{}

func NewMsgBangDataspace(creator string, name string, description string) *MsgBangDataspace {
	return &MsgBangDataspace{
		Creator:     creator,
		Name:        name,
		Description: description,
	}
}

func (msg *MsgBangDataspace) Route() string {
	return RouterKey
}

func (msg *MsgBangDataspace) Type() string {
	return TypeMsgBangDataspace
}

func (msg *MsgBangDataspace) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBangDataspace) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBangDataspace) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
