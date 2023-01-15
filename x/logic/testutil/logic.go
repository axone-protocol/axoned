package testutil

import (
	"context"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
)

// NewInterpreterMust returns a new Interpreter with the given context or panics if it fails.
// The Interpreter is configured with minimal settings to support testing.
func NewInterpreterMust(ctx context.Context) (interpreter *prolog.Interpreter) {
	interpreter = &prolog.Interpreter{}
	interpreter.Register3(engine.NewAtom("op"), engine.Op)
	interpreter.Register3(engine.NewAtom("compare"), engine.Compare)
	interpreter.Register2(engine.NewAtom("="), engine.Unify)

	err := interpreter.Compile(ctx, `
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

// CompileMust compiles the given source code and panics if it fails.
// This is a convenience function for testing.
func CompileMust(ctx context.Context, interpreter *prolog.Interpreter, s string, args ...interface{}) {
	err := interpreter.Compile(ctx, s, args...)
	if err != nil {
		panic(err)
	}
}
