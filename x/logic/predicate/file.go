package predicate

import (
	"context"
	"os"
	"reflect"
	"sort"

	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/v7/x/logic/prolog"
)

// SourceFile is a predicate which unifies the given term with the source file that is currently loaded.
//
// # Signature
//
//	source_file(?File) is det
//
// where:
//   - File represents the loaded source file.
func SourceFile(vm *engine.VM, file engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	loaded := getLoadedSources(vm)

	switch file := env.Resolve(file).(type) {
	case engine.Variable:
		promises := make([]func(ctx context.Context) *engine.Promise, 0, len(loaded))
		sortedSource := sortLoadedSources(loaded)
		for i := range sortedSource {
			term := engine.NewAtom(sortedSource[i])
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
		if _, ok := loaded[inputFile]; !ok {
			return engine.Bool(false)
		}
		return cont(env)
	default:
		return engine.Error(engine.TypeError(prolog.AtomTypeAtom, file, env))
	}
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

// Open is a predicate which opens a stream to a source or sink.
//
// # Signature
//
//	open(+SourceSink, +Mode, -Stream, +Options)
//
// where:
//   - SourceSink is an atom representing the source or sink of the stream.
//   - Mode is an atom representing the mode of the stream to be opened. It can be one of "read", "write", or "append".
//   - Stream is the stream to be opened.
//   - Options is a list of options. No options are currently defined, so the list should be empty.
//
// open/4 gives True when SourceSink can be opened in  Mode with the given Options.
//
// # Virtual File System (VFS)
//
// The logical module interprets on-chain Prolog programs, relying on a Virtual Machine that isolates execution from the
// external environment. Consequently, the open/4 predicate doesn't access the physical file system as one might expect.
// Instead, it operates with a Virtual File System (VFS), a conceptual layer that abstracts the file system. This abstraction
// offers a unified view across various storage systems, adhering to the constraints imposed by blockchain technology.
//
// This VFS extends the file concept to resources, which are identified by a Uniform Resource Identifier (URI). A URI
// specifies the access protocol for the resource, its path, and any necessary parameters.
//
// # CosmWasm URI
//
// The cosmwasm URI enables interaction with instantiated CosmWasm smart contract on the blockchain. The URI is used to
// query the smart contract and retrieve the response. The query is executed on the smart contract, and the response is
// returned as a stream. Query parameters are passed as part of the URI to customize the interaction with the smart contract.
//
// Its format is as follows:
//
//	cosmwasm:{contract_name}:{contract_address}?query={contract_query}[&base64Decode={true|false}]
//
// where:
//   - {contract_name}: For informational purposes, indicates the name or type of the smart contract (e.g., "okp4-objectarium").
//   - {contract_address}: Specifies the smart contract instance to query.
//   - {contract_query}: The query to be executed on the smart contract. It is a JSON object that specifies the query payload.
//   - base64Decode: (Optional) If true, the response is base64-decoded. Otherwise, the response is returned as is.
func Open(vm *engine.VM, sourceSink, mode, stream, options engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	var name string
	switch s := env.Resolve(sourceSink).(type) {
	case engine.Variable:
		return engine.Error(engine.InstantiationError(env))
	case engine.Atom:
		name = s.String()
	default:
		return engine.Error(engine.TypeError(prolog.AtomTypeAtom, sourceSink, env))
	}

	if prolog.IsGround(options, env) {
		_, err := prolog.AssertList(options, env)
		switch {
		case err != nil:
			return engine.Error(err)
		case !prolog.IsEmptyList(options, env):
			return engine.Error(engine.DomainError(prolog.ValidEmptyList(), options, env))
		}
	}

	var streamMode ioMode
	switch m := env.Resolve(mode).(type) {
	case engine.Variable:
		return engine.Error(engine.InstantiationError(env))
	case engine.Atom:
		var ok bool
		streamMode, ok = map[engine.Atom]ioMode{
			atomRead:   ioModeRead,
			atomWrite:  ioModeWrite,
			atomAppend: ioModeAppend,
		}[m]
		if !ok {
			return engine.Error(engine.TypeError(prolog.AtomTypeIOMode, mode, env))
		}
	default:
		return engine.Error(engine.TypeError(prolog.AtomTypeIOMode, mode, env))
	}

	if _, ok := env.Resolve(stream).(engine.Variable); !ok {
		// TODO: replace InstantiationError with uninstantiation_error(+Culprit) once it's implemented by ichiban/prolog.
		return engine.Error(engine.InstantiationError(env))
	}

	if streamMode != ioModeRead {
		return engine.Error(engine.PermissionError(prolog.AtomOperationInput, prolog.AtomPermissionTypeStream, sourceSink, env))
	}

	f, err := vm.FS.Open(name)
	if err != nil {
		return engine.Error(engine.ExistenceError(prolog.AtomObjectTypeSourceSink, sourceSink, env))
	}
	s := engine.NewInputTextStream(f)

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
