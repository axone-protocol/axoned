package types

import (
	sdkerrors "cosmossdk.io/errors"
	"google.golang.org/grpc/codes"
)

var (
	InvalidArgument = sdkerrors.RegisterWithGRPCCode(ModuleName, 1, codes.InvalidArgument, "limit exceeded")
	// LimitExceeded is returned when a limit is exceeded.
	LimitExceeded = sdkerrors.RegisterWithGRPCCode(ModuleName, 2, codes.InvalidArgument, "limit exceeded")
	// Internal is returned when an internal error occurs.
	Internal = sdkerrors.RegisterWithGRPCCode(ModuleName, 3, codes.Internal, "internal error")
)
