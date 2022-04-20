package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/knowledge module sentinel errors.
var (
	ErrEntityAlreadyExists = sdkerrors.Register(ModuleName, 1100, "entity already exists")
	ErrInvalidURI          = sdkerrors.Register(ModuleName, 1101, "invalid uri")
)
