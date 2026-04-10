package keeper

import (
	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
