package keeper

import (
	"github.com/axone/axoned/v7/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
