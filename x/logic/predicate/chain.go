package predicate

import (
	"context"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
)

// ChainID is a predicate which unifies the given term with the current chain ID. The signature is:
//
//	chain_id(?ChainID)
//
// where ChainID represents the current chain ID at the time of the query.
func ChainID(vm *engine.VM, chainID engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, chainID, engine.CharList(sdkContext.ChainID()), cont, env)
	})
}
