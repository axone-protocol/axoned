package wasm

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	logickeeper "github.com/okp4/okp4d/x/logic/keeper"
	logicwasm "github.com/okp4/okp4d/x/logic/wasm"
)

// customQuery represents the wasm custom query structure, it is intended to allow wasm contracts to execute queries
// against the logic module.
type customQuery struct {
	Ask *logicwasm.AskQuery `json:"ask,omitempty"`
}

// CustomQueryPlugins creates a wasm QueryPlugins containing the custom querier managing wasm contracts queries to the
// logic module.
func CustomQueryPlugins(logicKeeper logickeeper.Keeper) *wasmkeeper.QueryPlugins {
	return &wasmkeeper.QueryPlugins{
		Custom: makeCustomQuerier(
			logicwasm.MakeLogicQuerier(logicKeeper),
		),
	}
}

func makeCustomQuerier(logicQuerier logicwasm.LogicQuerier) wasmkeeper.CustomQuerier {
	return func(ctx sdk.Context, request json.RawMessage) ([]byte, error) {
		var query customQuery
		if err := json.Unmarshal(request, &query); err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		if query.Ask != nil {
			return logicQuerier.Ask(ctx, *query.Ask)
		}

		return nil, sdkerrors.Wrap(wasmtypes.ErrInvalidMsg, "Unknown custom query variant")
	}
}
