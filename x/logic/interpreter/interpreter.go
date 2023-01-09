package interpreter

import (
	goctx "context"
	"fmt"

	"github.com/ichiban/prolog"
	"github.com/okp4/okp4d/x/logic/context"
)

// NewInstrumentedInterpreter creates a new prolog.Interpreter with:
// - a list of predefined predicates
// - a compiled bootstrap script, that can be used to perform setup tasks.
// - a function that can be used to increment the gas meter.
//
// The predicates names must be present in the registry, otherwise the function will return an error.
// The bootstrap script can be an empty string if no bootstrap script is needed. If compilation of the bootstrap script
// fails, the function will return an error.
func NewInstrumentedInterpreter(
	ctx goctx.Context,
	predicates []string,
	bootstrap string,
	inc context.IncrementCountByFunc,
) (*prolog.Interpreter, error) {
	var i prolog.Interpreter

	for _, o := range predicates {
		if err := Register(ctx, &i, o, inc); err != nil {
			return nil, fmt.Errorf("error registering predicate '%s': %w", o, err)
		}
	}

	if err := i.Compile(ctx, bootstrap); err != nil {
		return nil, fmt.Errorf("error compiling bootstrap script: %w", err)
	}

	return &i, nil
}
