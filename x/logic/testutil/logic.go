package testutil

import (
	"context"
	"fmt"
	"regexp"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/x/logic/interpreter/bootstrap"
)

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
// The Interpreter is configured with the full boostrap but with a minimal set of predicates.
func NewComprehensiveInterpreterMust(ctx context.Context) (i *prolog.Interpreter) {
	i = &prolog.Interpreter{}
	i.Register3(engine.NewAtom("op"), engine.Op)
	i.Register3(engine.NewAtom("compare"), engine.Compare)
	i.Register2(engine.NewAtom("="), engine.Unify)
	i.Register1(engine.NewAtom("consult"), engine.Consult)
	i.Register3(engine.NewAtom("bagof"), engine.BagOf)
	i.Register1(engine.NewAtom("current_output"), engine.CurrentOutput)
	i.Register2(engine.NewAtom("put_char"), engine.PutChar)
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
	return prolog.TermString(re.ReplaceAllStringFunc(string(s), func(m string) string {
		index++
		return fmt.Sprintf("_%d", index)
	}))
}
