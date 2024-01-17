package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ sdk.Msg = &MsgCreateVestingAccount{}

var _ sdk.Msg = &MsgCreatePermanentLockedAccount{}

var _ sdk.Msg = &MsgCreatePeriodicVestingAccount{}

var _ sdk.Msg = &MsgCreateCliffVestingAccount{}

// NewMsgCreateVestingAccount returns a reference to a new MsgCreateVestingAccount.
//
//nolint:interfacer
func NewMsgCreateVestingAccount(fromAddr,
	toAddr sdk.AccAddress,
	amount sdk.Coins,
	endTime int64,
	delayed bool,
) *MsgCreateVestingAccount {
	return &MsgCreateVestingAccount{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      amount,
		EndTime:     endTime,
		Delayed:     delayed,
	}
}

// NewMsgCreatePermanentLockedAccount returns a reference to a new MsgCreatePermanentLockedAccount.
//
//nolint:interfacer
func NewMsgCreatePermanentLockedAccount(fromAddr, toAddr sdk.AccAddress, amount sdk.Coins) *MsgCreatePermanentLockedAccount {
	return &MsgCreatePermanentLockedAccount{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      amount,
	}
}

// NewMsgCreatePeriodicVestingAccount returns a reference to a new MsgCreatePeriodicVestingAccount.
//
//nolint:interfacer
func NewMsgCreatePeriodicVestingAccount(fromAddr,
	toAddr sdk.AccAddress,
	startTime int64,
	periods []Period,
) *MsgCreatePeriodicVestingAccount {
	return &MsgCreatePeriodicVestingAccount{
		FromAddress:    fromAddr.String(),
		ToAddress:      toAddr.String(),
		StartTime:      startTime,
		VestingPeriods: periods,
	}
}

// NewMsgCreateCliffVestingAccount returns a reference to a new MsgCreateCliffVestingAccount.
//
//nolint:interfacer
func NewMsgCreateCliffVestingAccount(fromAddr,
	toAddr sdk.AccAddress,
	amount sdk.Coins,
	endTime, cliffTime int64,
) *MsgCreateCliffVestingAccount {
	return &MsgCreateCliffVestingAccount{
		FromAddress: fromAddr.String(),
		ToAddress:   toAddr.String(),
		Amount:      amount,
		EndTime:     endTime,
		CliffTime:   cliffTime,
	}
}
