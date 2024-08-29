package keeper

import (
	"github.com/axone-protocol/axoned/v10/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
