package predicate

import (
	"context"

	"github.com/ichiban/prolog/engine"
)

// ChainID is higher order function that given a context returns the following predicate:
//
//	chain_id(?ChainID)
//
// where ChainID represents the current chain ID at the time of the query.
// The predicate is deterministic, producing the same chain ID each time it is called.
func ChainID(ctx context.Context) engine.Predicate1 {
	return func(vm *engine.VM, chainID engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		sdkContext, err := UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, chainID, engine.CharList(sdkContext.ChainID()), cont, env)
	}
}
