package interpreter

import (
	goctx "context"
	"fmt"
	"io/fs"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ichiban/prolog"
)

// Predicates is a map of predicate names to their execution costs.
type Predicates map[string]uint64

// New creates a new prolog.Interpreter with:
// - a list of predefined predicates (with their execution costs).
// - a compiled bootstrap script, that can be used to perform setup tasks.
// - a meter to track gas consumption.
// - a file system to load external files.
//
// The predicates names must be present in the registry, otherwise the function will return an error.
// The bootstrap script can be an empty string if no bootstrap script is needed. If compilation of the bootstrap script
// fails, the function will return an error.
func New(
	ctx goctx.Context,
	predicates Predicates,
	bootstrap string,
	meter sdk.GasMeter,
	fs fs.FS,
) (*prolog.Interpreter, error) {
	var i prolog.Interpreter
	i.FS = fs

	for predicate, cost := range predicates {
		if err := Register(&i, predicate, cost, meter); err != nil {
			return nil, fmt.Errorf("error registering predicate '%s': %w", predicate, err)
		}
	}

	if err := i.Compile(ctx, bootstrap); err != nil {
		return nil, fmt.Errorf("error compiling bootstrap script: %w", err)
	}

	return &i, nil
}
