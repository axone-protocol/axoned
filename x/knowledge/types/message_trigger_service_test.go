package types

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/okp4/okp4d/testutil/sample"
	"github.com/stretchr/testify/require"
)

func TestMsgTriggerService_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  MsgTriggerService
		err  error
	}{
		{
			name: "invalid address",
			msg: MsgTriggerService{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: MsgTriggerService{
				Creator: sample.AccAddress(),
			},
		},
		{
			name: "invalid uri",
			msg: MsgTriggerService{
				Creator: sample.AccAddress(),
				Uri:     "é:// ",
			},
			err: ErrInvalidURI,
		}, {
			name: "valid uri",
			msg: MsgTriggerService{
				Creator: sample.AccAddress(),
				Uri:     "okp4:service#coucou?target=okp4",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
