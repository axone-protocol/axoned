package wasm

import (
	"encoding/json"

	"github.com/okp4/okp4d/x/logic/keeper"
	"github.com/okp4/okp4d/x/logic/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// LogicQuerier ease the bridge between the logic module with the wasm CustomQuerier to allow wasm contracts to query
// the logic module.
type LogicQuerier struct {
	k *keeper.Keeper
}

// MakeLogicQuerier creates a new LogicQuerier based on the logic keeper.
func MakeLogicQuerier(keeper *keeper.Keeper) LogicQuerier {
	return LogicQuerier{
		k: keeper,
	}
}

// Ask is a proxy method with the gRPC request, returning the result in the json format.
func (querier LogicQuerier) Ask(ctx sdk.Context, query AskQuery) ([]byte, error) {
	grpcResp, err := querier.k.Ask(ctx, &types.QueryServiceAskRequest{
		Program: query.Program,
		Query:   query.Query,
	})
	if err != nil {
		return nil, err
	}

	resp := new(AskResponse)
	resp.from(*grpcResp)
	raw, err := json.Marshal(resp)

	querier.k.Logger(ctx).Debug("response to wasm ask", "json", string(raw))

	return raw, err
}
