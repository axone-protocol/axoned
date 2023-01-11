package predicate

import (
	"context"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
)

// ChainID is higher order function that given a context returns the following predicate:
//
//	chain_id(?ChainID)
//
// where ChainID represents the current chain ID at the time of the query.
// The predicate is deterministic, producing the same chain ID each time it is called.
func ChainID(ctx context.Context) func(*engine.VM, engine.Term, engine.Cont, *engine.Env) *engine.Promise {
	return func(vm *engine.VM, chainID engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, chainID, engine.CharList(sdkContext.ChainID()), cont, env)
	}
}
