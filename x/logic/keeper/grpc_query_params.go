package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v13/x/logic/types"
)

func (k Keeper) Params(c context.Context, req *types.QueryServiceParamsRequest) (*types.QueryServiceParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryServiceParamsResponse{Params: k.GetParams(ctx)}, nil
}
