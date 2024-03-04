package keeper

import (
	"github.com/okp4/okp4d/v7/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
