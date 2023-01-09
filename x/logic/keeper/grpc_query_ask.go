package keeper

import (
	goctx "context"

	sdkerrors "cosmossdk.io/errors"
	"github.com/okp4/okp4d/x/logic/types"
)

func (k Keeper) Ask(ctx goctx.Context, req *types.QueryServiceAskRequest) (*types.QueryServiceAskResponse, error) {
	if req == nil {
		return nil, sdkerrors.Wrap(types.InvalidArgument, "request is nil")
	}

	limits := k.getLimits(ctx)
	if err := checkLimits(req, limits); err != nil {
		return nil, err
	}

	return k.execute(ctx, req.Program, req.Query)
}
