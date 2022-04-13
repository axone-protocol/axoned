package types

import (
	"net/url"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTriggerService = "trigger_service"

var _ sdk.Msg = &MsgTriggerService{}

func NewMsgTriggerService(creator string, uri string) *MsgTriggerService {
	return &MsgTriggerService{
		Creator: creator,
		Uri:     uri,
	}
}

func (msg *MsgTriggerService) Route() string {
	return RouterKey
}

func (msg *MsgTriggerService) Type() string {
	return TypeMsgTriggerService
}

func (msg *MsgTriggerService) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgTriggerService) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTriggerService) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if _, err := url.Parse(msg.Uri); err != nil {
		return sdkerrors.Wrapf(ErrInvalidURI, "invalid service uri (%s)", err)
	}
	return nil
}
