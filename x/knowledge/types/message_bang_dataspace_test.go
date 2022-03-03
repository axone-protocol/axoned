package types_test

import (
	"testing"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/okp4/okp4d/testutil/sample"
	"github.com/okp4/okp4d/x/knowledge/types"
	"github.com/stretchr/testify/require"
)

func TestMsgBangDataspace_ValidateBasic(t *testing.T) {
	tests := []struct {
		name string
		msg  types.MsgBangDataspace
		err  error
	}{
		{
			name: "invalid address",
			msg: types.MsgBangDataspace{
				Creator: "invalid_address",
			},
			err: sdkerrors.ErrInvalidAddress,
		}, {
			name: "valid address",
			msg: types.MsgBangDataspace{
				Creator: sample.AccAddress(),
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
