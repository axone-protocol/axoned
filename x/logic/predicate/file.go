package predicate

import (
	"context"

	"github.com/axone-protocol/prolog/v3/engine"

	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
)

// SourceFile is a predicate which unifies the given term with the source file that is currently loaded.
//
// # Signature
//
//	source_file(?File) is det
//
// where:
//   - File represents the loaded source file.
//
// When File is a variable, solutions are produced in source loading order.
func SourceFile(vm *engine.VM, file engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	loaded := vm.LoadedSources()

	switch file := env.Resolve(file).(type) {
	case engine.Variable:
		promises := make([]func(ctx context.Context) *engine.Promise, 0, len(loaded))
		for i := range loaded {
			term := engine.NewAtom(loaded[i])
			promises = append(
				promises,
				func(_ context.Context) *engine.Promise {
					return engine.Unify(
						vm,
						file,
						term,
						cont,
						env,
					)
				})
		}

		return engine.Delay(promises...)
	case engine.Atom:
		inputFile := file.String()
		for i := range loaded {
			if loaded[i] == inputFile {
				return cont(env)
			}
		}
		return engine.Bool(false)
	default:
		return engine.Error(engine.TypeError(prolog.AtomTypeAtom, file, env))
	}
}
