package predicate

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"sort"

	"github.com/ichiban/prolog/engine"
)

// SourceFile is a predicate that unify the given term with the currently loaded source file.
//
// The signature is as follows:
//
//	source_file(?File).
//
// Where:
//   - File represents a loaded source file.
//
// Examples:
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

// ioMode describes what operations you can perform on the stream.
type ioMode int

const (
	// ioModeRead means you can read from the stream.
	ioModeRead = ioMode(os.O_RDONLY)
	// ioModeWrite means you can write to the stream.
	ioModeWrite = ioMode(os.O_CREATE | os.O_WRONLY)
	// ioModeAppend means you can append to the stream.
	ioModeAppend = ioMode(os.O_APPEND) | ioModeWrite
)

var (
	atomRead   = engine.NewAtom("read")
	atomWrite  = engine.NewAtom("write")
	atomAppend = engine.NewAtom("append")
)

func (m ioMode) Term() engine.Term {
	return [...]engine.Term{
		ioModeRead:   atomRead,
		ioModeWrite:  atomWrite,
		ioModeAppend: atomAppend,
	}[m]
}

// Open is a predicate that unify a stream with a source sink on a virtual file system.
//
// The signature is as follows:
//
//	open(+SourceSink, +Mode, ?Stream, +Options)
//
// Where:
//   - SourceSink is an atom representing the source or sink of the stream. The atom typically represents a resource
//     that can be opened, such as a URI. The URI scheme determines the type of resource that is opened.
//   - Mode is an atom representing the mode of the stream (read, write, append).
//   - Stream is the stream to be opened.
//   - Options is a list of options.
//
// Examples:
//
//	# Open a stream from a cosmwasm query.
//	# The Stream should be read as a string with a read_string/3 predicate, and then closed with the close/1 predicate.
//	- open('cosmwasm:okp4-objectarium:okp41lppz4x9dtmccek2m6cezjlwwzup6pdqrkvxjpk95806c3dewgrfq602kgx?query=%7B%22object_data%22%3A%7B%22id%22%3A%222625337e6025495a87cb32eb7f5a042f31e4385fd7e34c90d661bfc94dd539e3%22%7D%7D', 'read', Stream)
func Open(vm *engine.VM, sourceSink, mode, stream, options engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	var name string
	switch s := env.Resolve(sourceSink).(type) {
	case engine.Variable:
		return engine.Error(fmt.Errorf("open/4: source cannot be a variable"))
	case engine.Atom:
		name = s.String()
	default:
		return engine.Error(fmt.Errorf("open/4: invalid domain for source, should be an atom, got %T", s))
	}

	var streamMode ioMode
	switch m := env.Resolve(mode).(type) {
	case engine.Variable:
		return engine.Error(fmt.Errorf("open/4: streamMode cannot be a variable"))
	case engine.Atom:
		var ok bool
		streamMode, ok = map[engine.Atom]ioMode{
			atomRead:   ioModeRead,
			atomWrite:  ioModeWrite,
			atomAppend: ioModeAppend,
		}[m]
		if !ok {
			return engine.Error(fmt.Errorf("open/4: invalid open mode (read | write | append)"))
		}
	default:
		return engine.Error(fmt.Errorf("open/4: invalid domain for open mode, should be an atom, got %T", m))
	}

	if _, ok := env.Resolve(stream).(engine.Variable); !ok {
		return engine.Error(fmt.Errorf("open/4: stream can only be a variable, got %T", env.Resolve(stream)))
	}

	if streamMode != ioModeRead {
		return engine.Error(fmt.Errorf("open/4: only read mode is allowed here"))
	}

	f, err := vm.FS.Open(name)
	if err != nil {
		return engine.Error(fmt.Errorf("open/4: couldn't open stream: %w", err))
	}
	s := engine.NewInputTextStream(f)

	iter := engine.ListIterator{List: options, Env: env}
	for iter.Next() {
		return engine.Error(fmt.Errorf("open/4: options is not allowed here"))
	}

	return engine.Unify(vm, stream, s, k, env)
}

func getLoadedSources(vm *engine.VM) map[string]struct{} {
	loadedField := reflect.ValueOf(vm).Elem().FieldByName("loaded").MapKeys()
	loaded := make(map[string]struct{}, len(loadedField))
	for _, value := range loadedField {
		loaded[value.String()] = struct{}{}
	}

	return loaded
}

func sortLoadedSources(sources map[string]struct{}) []string {
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
