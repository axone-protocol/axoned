package util

import (
	"context"
	"strings"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"

	sdkmath "cosmossdk.io/math"

	"github.com/okp4/okp4d/x/logic/types"
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
	envs := make([]*engine.Env, 0, limit.Uint64())
	_, callErr := engine.Call(&i.VM, t, func(env *engine.Env) *engine.Promise {
		if count.LT(limit) {
			envs = append(envs, env)
		}
		count = count.Incr()
		return engine.Bool(count.GT(limit))
	}, env).Force(ctx)

	answerErr := lo.IfF(callErr != nil, func() string {
		return callErr.Error()
	}).Else("")
	success := len(envs) > 0
	hasMore := count.GT(limit)
	vars := parsedVarsToVars(p.Vars)
	results, err := envsToResults(envs, p.Vars, i)
	if err != nil {
		return nil, err
	}

	return &types.Answer{
		Success:   success,
		Error:     answerErr,
		HasMore:   hasMore,
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
			var expression prolog.TermString
			err := expression.Scan(&i.VM, v.Variable, rEnv)
			if err != nil {
				return nil, err
			}
			substitution := types.Substitution{
				Variable:   v.Name.String(),
				Expression: string(expression),
			}
			substitutions = append(substitutions, substitution)
		}
		results = append(results, types.Result{Substitutions: substitutions})
	}
	return results, nil
}
