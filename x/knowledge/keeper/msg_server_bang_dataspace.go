package keeper

import (
	"context"
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/knowledge/types"
)

func (k msgServer) BangDataspace(goCtx context.Context, msg *types.MsgBangDataspace) (*types.MsgBangDataspaceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.HasDataspace(ctx, msg.Id) {
		return nil, sdkerrors.Wrap(types.ErrEntityAlreadyExists, fmt.Sprintf("dataspace %s", msg.Id))
	}

	k.SaveDataspace(ctx, msg.Id, msg.Name)

	return &types.MsgBangDataspaceResponse{}, nil
}
