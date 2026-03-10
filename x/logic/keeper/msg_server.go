package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServiceServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServiceServer = msgServer{}

// UpdateParams implements the gRPC MsgServer interface. When an UpdateParams
// proposal passes, it updates the module parameters. The update can only be
// performed if the requested authority is the Cosmos SDK governance module
// account.
func (ms msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.authority.String() != req.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner,
			"invalid authority; expected %s, got %s", ms.authority.String(), req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

// StoreProgram implements the gRPC MsgServer interface.
// It stores a canonical program artifact keyed by program_id=SHA256(source),
// then records publisher-specific publication metadata idempotently.
func (ms msgServer) StoreProgram(goCtx context.Context, req *types.MsgStoreProgram) (*types.MsgStoreProgramResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidArgument, "request is nil")
	}

	publisher, err := sdk.AccAddressFromBech32(req.Publisher)
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInvalidArgument, "invalid publisher: %v", err.Error())
	}

	sdkCtx := sdk.UnwrapSDKContext(goCtx)
	params := ms.GetParams(sdkCtx)
	sourceSize := uint64(len(req.GetSource()))

	programIDRaw := sha256.Sum256([]byte(req.GetSource()))
	programID := hex.EncodeToString(programIDRaw[:])

	program, found, err := ms.GetStoredProgram(sdkCtx, programIDRaw[:])
	if err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "error reading stored program: %v", err.Error())
	}

	if !found {
		if params.Limits.MaxSize != 0 && sourceSize > params.Limits.MaxSize {
			return nil, errorsmod.Wrapf(types.ErrLimitExceeded, "source: %d > MaxSize: %d", sourceSize, params.Limits.MaxSize)
		}

		if err := ms.validateProgram(goCtx, params, req.GetSource()); err != nil {
			return nil, err
		}

		if err := ms.SetStoredProgram(sdkCtx, programIDRaw[:], types.StoredProgram{
			Source:     req.GetSource(),
			CreatedAt:  sdkCtx.BlockTime().Unix(),
			SourceSize: sourceSize,
		}); err != nil {
			return nil, errorsmod.Wrapf(types.ErrInternal, "error storing program: %v", err.Error())
		}
	} else if program.Source != req.GetSource() {
		return nil, errorsmod.Wrapf(types.ErrInternal, "program hash collision detected for id %s", programID)
	}

	if err := ms.EnsureProgramPublication(sdkCtx, publisher, programIDRaw[:], sdkCtx.BlockTime().Unix()); err != nil {
		return nil, errorsmod.Wrapf(types.ErrInternal, "error storing program publication: %v", err.Error())
	}

	return &types.MsgStoreProgramResponse{ProgramId: programID}, nil
}
