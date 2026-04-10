package keeper

import (
	"context"
	"encoding/hex"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"cosmossdk.io/store/prefix"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"

	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

func (k Keeper) Programs(c context.Context, req *types.QueryProgramsRequest) (*types.QueryProgramsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.StoredProgramKeyPrefix)
	programs := make([]types.ProgramMetadata, 0)

	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		if !accumulate {
			return true, nil
		}

		var program types.StoredProgram
		if err := k.cdc.Unmarshal(value, &program); err != nil {
			return false, err
		}

		programs = append(programs, newProgramMetadata(key, program))

		return true, nil
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to paginate programs: %v", err)
	}

	return &types.QueryProgramsResponse{
		Programs:   programs,
		Pagination: pageRes,
	}, nil
}

func (k Keeper) ProgramsByPublisher(
	c context.Context, req *types.QueryProgramsByPublisherRequest,
) (*types.QueryProgramsByPublisherResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.Publisher == "" {
		return nil, status.Error(codes.InvalidArgument, "publisher is required")
	}

	ctx := sdk.UnwrapSDKContext(c)
	publisher, err := sdk.AccAddressFromBech32(req.Publisher)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid publisher %q: %v", req.Publisher, err)
	}

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ProgramPublicationByPublisherPrefix(publisher))
	programs := make([]types.PublishedProgram, 0)

	pageRes, err := query.FilteredPaginate(store, req.Pagination, func(key, value []byte, accumulate bool) (bool, error) {
		if !accumulate {
			return true, nil
		}

		var publication types.ProgramPublication
		if err := k.cdc.Unmarshal(value, &publication); err != nil {
			return false, err
		}

		program, found, err := k.GetStoredProgram(ctx, key)
		if err != nil {
			return false, err
		}
		if !found {
			return false, status.Errorf(
				codes.Internal,
				"program %q referenced by publisher %q was not found",
				hex.EncodeToString(key),
				req.Publisher,
			)
		}

		programs = append(programs, types.PublishedProgram{
			Program:     newProgramMetadata(key, program),
			Publication: publication,
		})

		return true, nil
	})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			return nil, s.Err()
		}

		return nil, status.Errorf(codes.Internal, "failed to paginate publisher programs: %v", err)
	}

	return &types.QueryProgramsByPublisherResponse{
		Programs:   programs,
		Pagination: pageRes,
	}, nil
}
