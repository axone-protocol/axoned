package keeper

import (
	"github.com/axone-protocol/axoned/v13/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
