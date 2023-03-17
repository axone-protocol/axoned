package predicate

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/ichiban/prolog/engine"
)

// SourceFile is a predicate that unify the given term with the currently loaded source file. The signature is as follows:
//
// source_file(?File).
//
// Where File represents a loaded source file.
//
// Example:
//
//	# Query all the loaded source files, in alphanumeric order.
//	- source_file(File).
//
//	# Query the given source file is loaded.
//	- source_file('foo.pl').
func SourceFile(vm *engine.VM, file engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	loaded := getLoadedSources(vm)

	inputFile, err := getFile(env, file)
	if err != nil {
		return engine.Error(fmt.Errorf("source_file/1: %w", err))
	}

	if inputFile != nil {
		if _, ok := loaded[*inputFile]; ok {
			return engine.Unify(vm, file, engine.NewAtom(*inputFile), cont, env)
		}
		return engine.Delay()
	}

	promises := make([]func(ctx context.Context) *engine.Promise, 0, len(loaded))
	sortedSource := sortLoadedSources(loaded)
	for i := range sortedSource {
		term := engine.NewAtom(sortedSource[i])
		promises = append(
			promises,
			func(ctx context.Context) *engine.Promise {
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
}

func getLoadedSources(vm *engine.VM) map[string]interface{} {
	loadedField := reflect.ValueOf(vm).Elem().FieldByName("loaded").MapKeys()
	loaded := make(map[string]interface{}, len(loadedField))
	for _, value := range loadedField {
		loaded[value.String()] = nil
	}

	return loaded
}

func sortLoadedSources(sources map[string]interface{}) []string {
	result := make([]string, 0, len(sources))
	for filename := range sources {
		result = append(result, filename)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	return result
}

//nolint:nilnil
func getFile(env *engine.Env, term engine.Term) (*string, error) {
	switch file := env.Resolve(term).(type) {
	case engine.Variable:
	case engine.Atom:
		strFile := file.String()
		return &strFile, nil
	default:
		return nil, fmt.Errorf("cannot unify file with %T", term)
	}
	return nil, nil
}
