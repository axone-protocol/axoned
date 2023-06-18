package types

import (
	"google.golang.org/grpc/codes"

	sdkerrors "cosmossdk.io/errors"
)

var (
	InvalidArgument = sdkerrors.RegisterWithGRPCCode(ModuleName, 1, codes.InvalidArgument, "invalid argument")
	// LimitExceeded is returned when a limit is exceeded.
	LimitExceeded = sdkerrors.RegisterWithGRPCCode(ModuleName, 2, codes.InvalidArgument, "limit exceeded")
	// Internal is returned when an internal error occurs.
	Internal = sdkerrors.RegisterWithGRPCCode(ModuleName, 3, codes.Internal, "internal error")
)
