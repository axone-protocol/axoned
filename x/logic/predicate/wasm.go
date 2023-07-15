package predicate

import (
	"fmt"
	"context"
	"encoding/json"
	"strings"


	"github.com/ichiban/prolog/engine"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/logic/util"
	"github.com/okp4/okp4d/x/logic/types"
)

func NewWasmExtension(contractAddress sdk.AccAddress, name string) any {
	parts := strings.Split(name, "/")
	if len(parts) != 2 {
		return nil
	}

	arity := parts[1]
	switch arity {
	case "1":
		return func(vm *engine.VM, arg0 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
			return CosmWasmQuery(name, vm, contractAddress, []engine.Term{arg0}, cont, env)
		}
	case "2":
		return func(vm *engine.VM, arg0, arg1 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
			return CosmWasmQuery(name, vm, contractAddress, []engine.Term{arg0, arg1}, cont, env)
		}
	case "3":
		return func(vm *engine.VM, arg0, arg1, arg2 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
			return CosmWasmQuery(name, vm, contractAddress, []engine.Term{arg0, arg1, arg2}, cont, env)
		}
	default:
		return nil
	}
}

type PrologQueryMsg struct {
	PrologExtensionManifest *PrologExtensionManifest `json:"prolog_extension_manifest,omitempty"`
}

type PrologExtensionManifest struct {
	PredicateName string `json:"predicate_name"`
	Args []string `json:"args"`
}

type PrologQueryResult struct {
	Solve *PrologSolveResult `json:"solve"`
}

type PrologSolveResult struct {
	Solutions [][]string `json:"solutions"`
	// Continuation *SolveContinuation `json:"continuation"`
}

func solvePredicate(ctx sdk.Context, wasm types.WasmKeeper, contractAddr sdk.AccAddress, predicateName string, termArgs []engine.Term) ([][]engine.Term, error) {
	args := make([]string, len(termArgs))
	for i, arg := range termArgs {
		switch arg := arg.(type) {
		case engine.Atom:
			args[i] = arg.String()
		case engine.Variable:
			args[i] = ""
		}
	}

	msg := PrologQueryMsg {
		PrologExtensionManifest: &PrologExtensionManifest {
			PredicateName: predicateName,
			Args: args,
		},
	}
	bz, err := json.Marshal(msg)

	resbz, err := wasm.QuerySmart(ctx, contractAddr, bz)
	if err != nil {
		return nil, err
	}

	var res PrologQueryResult
	err = json.Unmarshal(resbz, &res)
	if err != nil {
		return nil, err
	}

	solutions := make([][]engine.Term, len(res.Solve.Solutions))
	for i, solution := range res.Solve.Solutions {
		solutions[i] = make([]engine.Term, len(solution))
		for j, atom := range solution {
			arg := termArgs[j]
			switch arg := arg.(type) {
			case engine.Atom:
				if arg.String() != atom {
					return nil, fmt.Errorf("unexpected atom: %s", atom)
				}
				solutions[i][j] = engine.NewAtom(atom)
			case engine.Variable:
				solutions[i][j] = engine.NewAtom(atom) // will be unified in CosmwasmQuery
			}
		}
	}

	return solutions, nil
}

func CosmWasmQuery(predicate string, vm *engine.VM, contractAddress sdk.AccAddress, args []engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		sdkContext, err := util.UnwrapSDKContext(ctx)
		if err != nil {
			return engine.Error(err)
		}
		wasmKeeper := sdkContext.Value(types.CosmWasmKeeperContextKey).(types.WasmKeeper)

		solutions, err := solvePredicate(sdkContext, wasmKeeper, contractAddress, predicate, args)
		if err != nil {
			return engine.Error(fmt.Errorf("%s: %w", predicate, err))
		}

		promises := make([]func(ctx context.Context) *engine.Promise, len(solutions))
		for i, solution := range solutions {
			promise := func(ctx context.Context) *engine.Promise {
				return engine.Unify(
					vm,
					Tuple(solution...),
					Tuple(args[2:]...),
					cont,
					env,
				)
			}
			promises[i] = promise
		}

		return engine.Delay(promises...)
	})
}
/*
func CosmWasmQuery3(vm *engine.VM, contractAddress engine.Term, predicateName engine.Term, arg0 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return CosmWasmQuery("cosmwasm_query/3", vm, contractAddress, predicateName, []engine.Term{arg0}, cont, env)
}

func CosmWasmQuery4(vm *engine.VM, contractAddress engine.Term, predicateName engine.Term, arg0, arg1 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return CosmWasmQuery("cosmwasm_query/4", vm, contractAddress, predicateName, []engine.Term{arg0, arg1}, cont, env)
}

func CosmWasmQuery5(vm *engine.VM, contractAddress engine.Term, predicateName engine.Term, arg0, arg1, arg2 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return CosmWasmQuery("cosmwasm_query/5", vm, contractAddress, predicateName, []engine.Term{arg0, arg1, arg2}, cont, env)
}

func CosmWasmQuery6(vm *engine.VM, contractAddress engine.Term, predicateName engine.Term, arg0, arg1, arg2, arg3 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return CosmWasmQuery("cosmwasm_query/6", vm, contractAddress, predicateName, []engine.Term{arg0, arg1, arg2, arg3}, cont, env)
}
*/