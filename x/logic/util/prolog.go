package util

import (
	"context"
	"strings"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"

	sdkmath "cosmossdk.io/math"

	"github.com/okp4/okp4d/v7/x/logic/types"
)

const (
	defaultEnvCap = uint64(50)
)

// QueryInterpreter interprets a query and returns the solutions up to the given limit.
func QueryInterpreter(
	ctx context.Context, i *prolog.Interpreter, query string, limit sdkmath.Uint,
) (*types.Answer, error) {
	p := engine.NewParser(&i.VM, strings.NewReader(query))
	t, err := p.Term()
	if err != nil {
		return nil, err
	}

	var env *engine.Env
	count := sdkmath.ZeroUint()
	envs := make([]*engine.Env, 0, sdkmath.MinUint(limit, sdkmath.NewUint(defaultEnvCap)).Uint64())
	_, callErr := engine.Call(&i.VM, t, func(env *engine.Env) *engine.Promise {
		if count.LT(limit) {
			envs = append(envs, env)
		}
		count = count.Incr()
		return engine.Bool(count.GT(limit))
	}, env).Force(ctx)

	vars := parsedVarsToVars(p.Vars)

	results, err := envsToResults(envs, p.Vars, i)
	if err != nil {
		return nil, err
	}

	if callErr != nil {
		if sdkmath.NewUint(uint64(len(results))).LT(limit) {
			// error is not part of the look-ahead and should be included in the solutions
			results = append(results, types.Result{Error: callErr.Error()})
		} else {
			// error is part of the look-ahead, so let's consider that there's one more solution
			count = count.Incr()
		}
	}

	return &types.Answer{
		HasMore:   count.GT(limit),
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
