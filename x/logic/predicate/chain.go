package predicate

import (
	"context"
	"fmt"

	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/x/logic/util"
)

// ChainID is a predicate which unifies the given term with the current chain ID. The signature is:
//
// The signature is as follows:
//
//	chain_id(?ID)
//
// where:
//   - ID represents the current chain ID at the time of the query.
//
// Examples:
//
//	# Query the current chain ID.
//	- chain_id(ID).
func ChainID(vm *engine.VM, chainID engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(fmt.Errorf("chain_id/1: %w", err))
		}

		return engine.Unify(vm, chainID, engine.NewAtom(sdkContext.ChainID()), cont, env)
	})
}
