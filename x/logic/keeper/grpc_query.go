package keeper

import (
	"github.com/axone-protocol/axoned/v11/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
