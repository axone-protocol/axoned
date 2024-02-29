package interpreter

import (
	goctx "context"
	"fmt"
	"io"
	"io/fs"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"
)

// Option is a function that configures an Interpreter.
type Option func(*prolog.Interpreter) error

// WithPredicates configures the interpreter to register the specified predicates.
// See WithPredicate for more details.
func WithPredicates(ctx goctx.Context, predicates []string, hook Hook) Option {
	return func(i *prolog.Interpreter) error {
		for _, predicate := range predicates {
			if err := WithPredicate(ctx, predicate, hook)(i); err != nil {
				return err
			}
		}
		return nil
	}
}

// WithPredicate configures the interpreter to register the specified predicate with the specified hook.
// The hook is a function that is called before the predicate is executed and can be used to check some conditions,
// like the gas consumption or the permission to execute the predicate.
//
// The predicates names must be present in the registry, otherwise the function will return an error.
func WithPredicate(_ goctx.Context, predicate string, hook Hook) Option {
	return func(i *prolog.Interpreter) error {
		if err := Register(i, predicate, hook); err != nil {
			return fmt.Errorf("error registering predicate '%s': %w", predicate, err)
		}
		return nil
	}
}

// WithBootstrap configures the interpreter to compile the specified bootstrap script to serve as setup context.
// If compilation of the bootstrap script fails, the function will return an error.
func WithBootstrap(ctx goctx.Context, bootstrap string) Option {
	return func(i *prolog.Interpreter) error {
		if err := i.Compile(ctx, bootstrap); err != nil {
			return fmt.Errorf("error compiling bootstrap script: %w", err)
		}
		return nil
	}
}

// WithUserOutputWriter configures the interpreter to use the specified writer for user output.
func WithUserOutputWriter(w io.Writer) Option {
	return func(i *prolog.Interpreter) error {
		i.SetUserOutput(engine.NewOutputTextStream(w))

		return nil
	}
}

// WithFS configures the interpreter to use the specified file system.
func WithFS(fs fs.FS) Option {
	return func(i *prolog.Interpreter) error {
		i.FS = fs
		return nil
	}
}

// New creates a new prolog.Interpreter with the specified options.
func New(
	opts ...Option,
) (*prolog.Interpreter, error) {
	var i prolog.Interpreter

	for _, opt := range opts {
		if err := opt(&i); err != nil {
			return nil, err
		}
	}

	return &i, nil
}
