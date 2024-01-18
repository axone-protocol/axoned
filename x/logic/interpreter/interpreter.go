package interpreter

import (
	goctx "context"
	"fmt"
	"io"
	"io/fs"

	"github.com/ichiban/prolog"
	"github.com/ichiban/prolog/engine"

	storetypes "cosmossdk.io/store/types"
)

// Predicates is a map of predicate names to their execution costs.
type Predicates map[string]uint64

// Option is a function that configures an Interpreter.
type Option func(*prolog.Interpreter) error

// WithPredicates configures the interpreter to register the specified predicates.
// The predicates names must be present in the registry, otherwise the function will return an error.
func WithPredicates(_ goctx.Context, predicates Predicates, meter storetypes.GasMeter) Option {
	return func(i *prolog.Interpreter) error {
		for predicate, cost := range predicates {
			if err := Register(i, predicate, cost, meter); err != nil {
				return fmt.Errorf("error registering predicate '%s': %w", predicate, err)
			}
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
