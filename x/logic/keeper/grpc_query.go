package keeper

import (
	"github.com/axone-protocol/axoned/v12/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
