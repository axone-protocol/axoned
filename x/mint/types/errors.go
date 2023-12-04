package types

import sdkerrors "cosmossdk.io/errors"

var ErrBondedRatioIsZero = sdkerrors.Register(ModuleName, 1, "bonded ratio is zero")
