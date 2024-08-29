package wasm

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	logickeeper "github.com/axone-protocol/axoned/v10/x/logic/keeper"
	logicwasm "github.com/axone-protocol/axoned/v10/x/logic/wasm"
)

// customQuery represents the wasm custom query structure, it is intended to allow wasm contracts to execute queries
// against the logic module.
type customQuery struct {
	Ask *logicwasm.AskQuery `json:"ask,omitempty"`
}

// CustomQueryPlugins creates a wasm QueryPlugins containing the custom querier managing wasm contracts queries to the
// logic module.
func CustomQueryPlugins(logicKeeper *logickeeper.Keeper) *wasmkeeper.QueryPlugins {
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
			return nil, errorsmod.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
		}

		if query.Ask != nil {
			return logicQuerier.Ask(ctx, *query.Ask)
		}

		return nil, errorsmod.Wrap(wasmtypes.ErrInvalidMsg, "Unknown custom query variant")
	}
}
