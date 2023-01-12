package predicate

import (
	"context"

	"github.com/ichiban/prolog/engine"
	"github.com/okp4/okp4d/x/logic/util"
)

// BlockHeight is a predicate which unifies the given term with the current block height. The signature is:
//
//	block_height(?Height)
//
// where Height represents the current chain height at the time of the query.
func BlockHeight(vm *engine.VM, height engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, height, engine.Integer(sdkContext.BlockHeight()), cont, env)
	})
}

// BlockTime is a predicate which unifies the given term with the current block time. The signature is:
//
//	block_time(?Time)
//
// where Time represents the current chain time at the time of the query.
func BlockTime(vm *engine.VM, time engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, time, engine.Integer(sdkContext.BlockTime().Unix()), cont, env)
	})
}
