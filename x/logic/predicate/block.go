package predicate

import (
	"context"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
)

// BlockHeight is higher order function that given a context returns the following predicate:
//
//	block_height(?Height)
//
// where Height represents the current chain height at the time of the query.
// The predicate is non-deterministic, producing a different height each time it is called.
func BlockHeight(ctx context.Context) func(*engine.VM, engine.Term, engine.Cont, *engine.Env) *engine.Promise {
	return func(vm *engine.VM, height engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, height, engine.Integer(sdkContext.BlockHeight()), cont, env)
	}
}

// BlockTime is higher order function that given a context returns the following predicate:
//
//	block_time(?Time)
//
// where Time represents the current chain time at the time of the query.
// The predicate is non-deterministic, producing a different time each time it is called.
func BlockTime(ctx context.Context) func(*engine.VM, engine.Term, engine.Cont, *engine.Env) *engine.Promise {
	return func(vm *engine.VM, time engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, time, engine.Integer(sdkContext.BlockTime().Unix()), cont, env)
	}
}
