package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/knowledge/types"
)

func (k msgServer) TriggerService(goCtx context.Context, msg *types.MsgTriggerService) (*types.MsgTriggerServiceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// TODO: Handling the message
	_ = ctx

	return &types.MsgTriggerServiceResponse{}, nil
}
