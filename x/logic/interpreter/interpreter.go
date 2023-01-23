package interpreter

import (
	goctx "context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog"
)

// New creates a new prolog.Interpreter with:
// - a list of predefined predicates
// - a compiled bootstrap script, that can be used to perform setup tasks.
//
// The predicates names must be present in the registry, otherwise the function will return an error.
// The bootstrap script can be an empty string if no bootstrap script is needed. If compilation of the bootstrap script
// fails, the function will return an error.
func New(
	ctx goctx.Context,
	predicates []string,
	bootstrap string,
	meter sdk.GasMeter,
) (*prolog.Interpreter, error) {
	var i prolog.Interpreter

	for _, o := range predicates {
		if err := Register(&i, o, meter); err != nil {
			return nil, fmt.Errorf("error registering predicate '%s': %w", o, err)
		}
	}

	if err := i.Compile(ctx, bootstrap); err != nil {
		return nil, fmt.Errorf("error compiling bootstrap script: %w", err)
	}

	return &i, nil
}
