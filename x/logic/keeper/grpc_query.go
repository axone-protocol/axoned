package keeper

import (
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
