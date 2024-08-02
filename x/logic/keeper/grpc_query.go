package keeper

import (
	"github.com/axone-protocol/axoned/v9/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
