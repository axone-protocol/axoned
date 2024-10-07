package testutil

import (
	"context"
	"fmt"
	"regexp"

	"github.com/axone-protocol/prolog"
	"github.com/axone-protocol/prolog/engine"

	"github.com/axone-protocol/axoned/v10/x/logic/interpreter/bootstrap"
)

type TermResults map[string]prolog.TermString

// NewLightInterpreterMust returns a new Interpreter with the given context or panics if it fails.
// The Interpreter is configured with minimal settings to support testing.
func NewLightInterpreterMust(ctx context.Context) (i *prolog.Interpreter) {
	i = &prolog.Interpreter{}
	i.Register3(engine.NewAtom("op"), engine.Op)
	i.Register3(engine.NewAtom("compare"), engine.Compare)
	i.Register2(engine.NewAtom("="), engine.Unify)
	i.Register1(engine.NewAtom("consult"), engine.Consult)

	err := i.Compile(ctx, `
						:-(op(1200, xfx, ':-')).
						:-(op(1000, xfy, ',')).
						:-(op(700, xfx, [==, \==, @<, @=<, @>, @>=])).
						:-(op(700, xfx, '=')).
						:-(op(500, yfx, [+, -, /\, \/])).

						member(X, [X|_]).
						member(X, [_|Xs]) :- member(X, Xs).
						X == Y :- compare(=, X, Y).`)
	if err != nil {
		panic(err)
	}

	return
}

// NewComprehensiveInterpreterMust returns a new Interpreter with the given context or panics if it fails.
// The Interpreter is configured with the full bootstrap but with a minimal set of predicates.
func NewComprehensiveInterpreterMust(ctx context.Context) (i *prolog.Interpreter) {
	i = &prolog.Interpreter{}
	i.Register3(engine.NewAtom("op"), engine.Op)
	i.Register3(engine.NewAtom("compare"), engine.Compare)
	i.Register2(engine.NewAtom("="), engine.Unify)
	i.Register1(engine.NewAtom("consult"), engine.Consult)
	i.Register3(engine.NewAtom("bagof"), engine.BagOf)
	i.Register1(engine.NewAtom("current_output"), engine.CurrentOutput)
	i.Register1(engine.NewAtom("current_input"), engine.CurrentInput)
	i.Register2(engine.NewAtom("put_char"), engine.PutChar)
	i.Register2(engine.NewAtom("get_char"), engine.GetChar)
	i.Register3(engine.NewAtom("write_term"), engine.WriteTerm)

	err := i.Compile(ctx, bootstrap.Bootstrap())
	if err != nil {
		panic(err)
	}

	return
}

// CompileMust compiles the given source code and panics if it fails.
// This is a convenience function for testing.
func CompileMust(ctx context.Context, i *prolog.Interpreter, s string, args ...interface{}) {
	err := i.Compile(ctx, s, args...)
	if err != nil {
		panic(err)
	}
}

// ReindexUnknownVariables reindexes the variables in the given term so that the variables are numbered sequentially.
// This is required for test predictability when the term is a result of a query and the variables are unknown.
//
// For example, the following term:
//
//	foo(_1, _2, _3, _1)
//
// is re-indexed as:
//
//	foo(_1, _2, _3, _4)
func ReindexUnknownVariables(s prolog.TermString) prolog.TermString {
	re := regexp.MustCompile("_([0-9]+)")
	var index int
	return prolog.TermString(re.ReplaceAllStringFunc(string(s), func(_ string) string {
		index++
		return fmt.Sprintf("_%d", index)
	}))
}

// ShouldBeGrounded is a goconvey assertion that asserts that the given term does not hold any
// uninstantiated variables.
func ShouldBeGrounded(actual any, expected ...any) string {
	if len(expected) != 0 {
		return fmt.Sprintf("This assertion requires exactly %d comparison values (you provided %d).", 0, len(expected))
	}

	var containsVariable func(engine.Term) bool
	containsVariable = func(term engine.Term) bool {
		switch t := term.(type) {
		case engine.Variable:
			return true
		case engine.Compound:
			for i := 0; i < t.Arity(); i++ {
				if containsVariable(t.Arg(i)) {
					return true
				}
			}
		}
		return false
	}
	if t, ok := actual.(engine.Term); ok {
		if containsVariable(t) {
			return "Expected term to NOT hold a free variable (but it was)."
		}

		return ""
	}

	return fmt.Sprintf("The argument to this assertion must be a term (you provided %v).", actual)
}
