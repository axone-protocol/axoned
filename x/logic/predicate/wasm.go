package predicate

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/okp4/okp4d/x/logic/util"
)

func QueryWasm(vm *engine.VM, contractAddr engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(fmt.Errorf("query_wasm/1: %w", err))
		}
		wasmKeeper := sdkContext.Value(types.WasmKeeperContextKey).(types.WasmKeeper)
		addr, err := getBech32(env, contractAddr)
		if err != nil {
			return engine.Error(fmt.Errorf("query_wasm/1: %w", err))
		}

		req := []byte("{\"ask\":{\"query\": \"query_wasm('okp410gnd30r45k9658jm7hzxvp8ehz4ptf33tqjnaepwkunev6kax5ks3mnvmf').\"}}")
		if !json.Valid(req) {
			return engine.Error(fmt.Errorf("query_wasm/1: wasm query must be a valid json"))
		}
		res, err := wasmKeeper.QuerySmart(sdkContext, addr, req)
		if err != nil {
			return engine.Error(fmt.Errorf("query_wasm/1: %w", err))
		}
		fmt.Printf("result %w", string(res))
		return engine.Unify(vm, contractAddr, engine.Integer(sdkContext.BlockHeight()), cont, env)
	})
}
