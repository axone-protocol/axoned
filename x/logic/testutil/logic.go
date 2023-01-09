package testutil

import (
	"context"

	"github.com/ichiban/prolog/engine"
)

// NewVMMust returns a new VM with the given context or panics if it fails.
// The VM is configured with minimal settings to support testing.
func NewVMMust(ctx context.Context) (vm *engine.VM) {
	vm = &engine.VM{}
	vm.Register3(engine.NewAtom("op"), engine.Op)
	vm.Register3(engine.NewAtom("compare"), engine.Compare)

	err := vm.Compile(ctx, `
						:-(op(1200, xfx, ':-')).
						:-(op(1000, xfy, ',')).
						:-(op(700, xfx, '==')).
						X == Y :- compare(=, X, Y).`)
	if err != nil {
		panic(err)
	}

	return
}

// CompileMust compiles the given source code and panics if it fails.
// This is a convenience function for testing.
func CompileMust(ctx context.Context, vm *engine.VM, s string, args ...interface{}) {
	err := vm.Compile(ctx, s, args...)
	if err != nil {
		panic(err)
	}
}
