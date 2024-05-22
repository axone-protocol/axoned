package keeper

import (
	"github.com/axone-protocol/axoned/v8/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
