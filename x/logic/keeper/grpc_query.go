package keeper

import (
	"github.com/okp4/okp4d/x/logic/types"
)

var _ types.QueryServiceServer = Keeper{}
