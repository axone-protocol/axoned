package keeper

import (
	"context"

	"cosmossdk.io/errors"

	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/okp4/okp4d/v7/x/mint/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

// UpdateParams implements the gRPC MsgServer interface. When an UpdateParams
// proposal passes, it updates the module parameters. The update can only be
// performed if the requested authority is the Cosmos SDK governance module
// account.
func (ms msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.authority != msg.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, msg.Authority)
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := ms.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}
