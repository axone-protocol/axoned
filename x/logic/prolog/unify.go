package prolog

import "github.com/ichiban/prolog/engine"

// ConvertFunc is a function mapping a domain which is a list of terms with a codomain which is a set of terms.
// Domains and co-domains can have different cardinalities.
// options is a list of options that can be used to parameterize the conversion.
// All the terms provided are fully instantiated (i.e. no variables).
type ConvertFunc func(value []engine.Term, options engine.Term, env *engine.Env) ([]engine.Term, error)

// UnifyFunctional is a generic unification which unifies a set of input terms with a set of output terms, using the
// given conversion functions maintaining the function's relationship.
//
// The aim of this function is to simplify the implementation of a wide range of predicates which are essentially
// functional, like hash functions, encoding functions, etc.
//
// The semantic of the unification is as follows:
//  1. first all the variables are resolved
//  2. if there's variables in the input and the output,
//     the conversion is not possible and a not sufficiently instantiated error is returned.
//  3. if there's no variables in the input,
//     then the conversion is attempted from the input to the output and the result is unified with the output.
//  4. if there's no variables in the output,
//     then the conversion is attempted from the output to the input and the result is unified with the input.
//
// The following table summarizes the behavior, where:
// - fi = fully instantiated (i.e. no variables)
// - !fi = not fully instantiated (i.e. at least one variable)
//
// | input | output | result                               |
// |-------|--------|--------------------------------------|
// | !fi   | !fi    | error: not sufficiently instantiated |
// |  fi   | !fi    | unify(forward(input), output)        |
// |  fi   |  fi    | unify(forward(input), output)        |
// | !fi   |  fi    | unify(input,backward(output))        |
//
// Conversion functions may produce an error in scenarios where the conversion is unsuccessful or infeasible due to
// the inherent characteristics of the function's relationship, such as the absence of a one-to-one correspondence
// (e.g. hash functions).
func UnifyFunctional(
	in,
	out []engine.Term,
	options engine.Term,
	forwardConverter ConvertFunc,
	backwardConverter ConvertFunc,
	env *engine.Env,
) (bool, *engine.Env, error) {
	isInFI, isOutFi := AreGround(in, env), AreGround(out, env)
	if !isInFI && !isOutFi {
		return false, env, engine.InstantiationError(env)
	}

	var err error
	from, to := in, out

	switch {
	case forwardConverter == nil && backwardConverter == nil:
		// no-op
	case isInFI && forwardConverter != nil:
		from, err = forwardConverter(in, options, env)
		if err != nil {
			return false, env, err
		}
	case isOutFi && backwardConverter != nil:
		to, err = backwardConverter(out, options, env)
		if err != nil {
			return false, env, err
		}
	default:
		return false, env, engine.InstantiationError(env)
	}

	env, result := env.Unify(
		Tuple(from...),
		Tuple(to...),
	)

	return result, env, nil
}

// UnifyFunctionalPredicate is the predicate version of UnifyFunctional returning a promise.
func UnifyFunctionalPredicate(
	in,
	out []engine.Term,
	options engine.Term,
	forwardConverter ConvertFunc,
	backwardConverter ConvertFunc,
	cont engine.Cont,
	env *engine.Env,
) *engine.Promise {
	ok, env, err := UnifyFunctional(in, out, options, forwardConverter, backwardConverter, env)
	if err != nil {
		return engine.Error(err)
	}
	if !ok {
		return engine.Bool(false)
	}
	return cont(env)
}
