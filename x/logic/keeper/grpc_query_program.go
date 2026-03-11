package keeper

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

func (k Keeper) Program(c context.Context, req *types.QueryProgramRequest) (*types.QueryProgramResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	programID, err := decodeProgramID(req.ProgramId)
	if err != nil {
		return nil, err
	}

	program, found, err := k.GetStoredProgram(ctx, programID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load program %q: %v", req.ProgramId, err)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "program %q not found", req.ProgramId)
	}

	return &types.QueryProgramResponse{
		Program: newProgramMetadata(programID, program),
	}, nil
}

func (k Keeper) ProgramSource(c context.Context, req *types.QueryProgramSourceRequest) (*types.QueryProgramSourceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	programID, err := decodeProgramID(req.ProgramId)
	if err != nil {
		return nil, err
	}

	program, found, err := k.GetStoredProgram(ctx, programID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to load program %q: %v", req.ProgramId, err)
	}
	if !found {
		return nil, status.Errorf(codes.NotFound, "program %q not found", req.ProgramId)
	}

	return &types.QueryProgramSourceResponse{Source: program.Source}, nil
}

func newProgramMetadata(programID []byte, program types.StoredProgram) types.ProgramMetadata {
	return types.ProgramMetadata{
		ProgramId:  hex.EncodeToString(programID),
		CreatedAt:  program.CreatedAt,
		SourceSize: program.SourceSize,
	}
}

func decodeProgramID(programID string) ([]byte, error) {
	if programID == "" {
		return nil, status.Error(codes.InvalidArgument, "program_id is required")
	}

	decoded, err := hex.DecodeString(programID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid program_id %q: %v", programID, err)
	}
	if len(decoded) == 0 {
		return nil, status.Error(codes.InvalidArgument, "program_id is required")
	}
	if len(decoded) != sha256.Size {
		return nil, status.Errorf(codes.InvalidArgument, "invalid program_id %q: expected %d-byte SHA-256 digest", programID, sha256.Size)
	}

	return decoded, nil
}
