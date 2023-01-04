package predicate

import (
	"context"

	"github.com/ichiban/prolog/engine"
)

// BlockHeight is higher order function that given a context returns the following predicate:
//
//	block_height(?Height)
//
// where Height represents the current chain height at the time of the query.
// The predicate is non-deterministic, producing a different height each time it is called.
func BlockHeight(ctx context.Context) engine.Predicate1 {
	return func(vm *engine.VM, height engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		sdkContext, err := UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, height, engine.Integer(sdkContext.BlockHeight()), cont, env)
	}
}

