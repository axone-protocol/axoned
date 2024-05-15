package keeper

import (
	"github.com/axone-protocol/axoned/v7/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
