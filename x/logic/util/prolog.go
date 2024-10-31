package util

import (
	"context"
	"errors"
	"strings"

	"github.com/axone-protocol/prolog"
	"github.com/axone-protocol/prolog/engine"
	"github.com/samber/lo"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v10/x/logic/types"
)

const (
	defaultEnvCap = uint64(50)
)

// QueryInterpreter interprets a query and returns the solutions up to the given limit.
//
//nolint:nestif,funlen
func QueryInterpreter(
	ctx context.Context, i *prolog.Interpreter, query string, solutionsLimit uint64,
) (*types.Answer, error) {
	p := engine.NewParser(&i.VM, strings.NewReader(query))
	t, err := p.Term()
	if err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error executing query: %v", err.Error())
	}

	var env *engine.Env
	count := uint64(0)
	envs := make([]*engine.Env, 0, min(solutionsLimit, defaultEnvCap))
	_, callErr := engine.Call(&i.VM, t, func(env *engine.Env) *engine.Promise {
		if count < solutionsLimit {
			envs = append(envs, env)
		}
		count++
		return engine.Bool(count > solutionsLimit)
	}, env).Force(ctx)

	vars := parsedVarsToVars(p.Vars)
	results, err := envsToResults(envs, p.Vars, i)
	if err != nil {
		return nil, errorsmod.Wrapf(types.InvalidArgument, "error executing query: %v", err.Error())
	}

	if callErr != nil {
		if uint64(len(results)) < solutionsLimit {
			// error is not part of the look-ahead and should be included in the solutions
			if errors.Is(callErr, types.LimitExceeded) {
				return nil, callErr
			}

			var panicErr engine.PanicError
			if errors.As(callErr, &panicErr) && errors.Is(panicErr.OriginErr, engine.ErrMaxVariables) {
				return nil, errorsmod.Wrapf(types.LimitExceeded, panicErr.OriginErr.Error()) //nolint:govet
			}

			if err = func() error {
				defer func() {
					_ = recover()
				}()
				sdkCtx := sdk.UnwrapSDKContext(ctx)
				if sdkCtx.GasMeter().IsOutOfGas() {
					return errorsmod.Wrapf(
						types.LimitExceeded, "out of gas: %s <%s> (%d/%d)",
						types.ModuleName, callErr.Error(), sdkCtx.GasMeter().GasConsumed(), sdkCtx.GasMeter().Limit())
				}
				return nil
			}(); err != nil {
				return nil, err
			}

			results = append(results, types.Result{Error: callErr.Error()})
		} else {
			// error is part of the look-ahead, so let's consider that there's one more solution
			count++
		}
	}

	return &types.Answer{
		HasMore:   count > solutionsLimit,
		Variables: vars,
		Results:   results,
	}, nil
}

func parsedVarsToVars(vars []engine.ParsedVariable) []string {
	return lo.Map(vars, func(v engine.ParsedVariable, _ int) string {
		return v.Name.String()
	})
}

func envsToResults(envs []*engine.Env, vars []engine.ParsedVariable, i *prolog.Interpreter) ([]types.Result, error) {
	results := make([]types.Result, 0, len(envs))
	for _, rEnv := range envs {
		substitutions := make([]types.Substitution, 0, len(vars))
		for _, v := range vars {
			if !isBound(v, rEnv) {
				// skip parsed variables that are not bound (singletons variables or other)
				continue
			}

			substitution, err := scanExpression(i, v, rEnv)
			if err != nil {
				return results, err
			}
			substitutions = append(substitutions, substitution)
		}
		results = append(results, types.Result{Substitutions: substitutions})
	}
	return results, nil
}

func scanExpression(i *prolog.Interpreter, v engine.ParsedVariable, rEnv *engine.Env) (types.Substitution, error) {
	var expression prolog.TermString
	err := expression.Scan(&i.VM, v.Variable, rEnv)
	if err != nil {
		return types.Substitution{}, err
	}
	substitution := types.Substitution{
		Variable:   v.Name.String(),
		Expression: string(expression),
	}

	return substitution, nil
}

// isBound returns true if the given parsed variable is bound in the given environment.
func isBound(v engine.ParsedVariable, env *engine.Env) bool {
	_, ok := env.Resolve(v.Variable).(engine.Variable)

	return !ok
}
