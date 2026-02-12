package predicate

import (
	"context"

	"github.com/axone-protocol/prolog/v3/engine"

	"github.com/axone-protocol/axoned/v13/x/logic/prolog"
)

var atomOpen = engine.NewAtom("open")

// Consult is a predicate which read files as Prolog source code.
//
// # Signature
//
//	consult(+Files) is det
//
// where:
//   - Files represents the source files to be loaded. It can be an atom or a list of atoms representing the source files.
//
// The Files argument are typically URIs that point to the sources file to be loaded through the Virtual File System (VFS).
// Please refer to the open/4 predicate for more information about the VFS.
func Consult(vm *engine.VM, file engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Consult(vm, file, cont, env)
}

// SourceFile is a predicate which unifies the given term with the source file that is currently loaded.
//
// # Signature
//
//	source_file(?File) is det
//
// where:
//   - File represents the loaded source file.
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

// Open is a predicate which opens a stream to a source or sink.
//
// # Signature
//
//	open(+SourceSink, +Mode, -Stream, +Options)
//
// where:
//   - SourceSink is an atom representing the source or sink of the stream, which is typically a URI.
//   - Mode is an atom representing the mode of the stream to be opened. It can be one of "read", "write", or "append".
//   - Stream is the stream to be opened.
//   - Options is a list of options. No options are currently defined, so the list should be empty.
//
// open/4 gives True when SourceSink can be opened in Mode with the given Options.
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
//   - {contract_name}: For informational purposes, indicates the name or type of the smart contract (e.g., "axone-objectarium").
//   - {contract_address}: Specifies the smart contract instance to query.
//   - {contract_query}: The query to be executed on the smart contract. It is a JSON object that specifies the query payload.
//   - base64Decode: (Optional) If true, the response is base64-decoded. Otherwise, the response is returned as is.
func Open(vm *engine.VM, sourceSink, mode, stream, options engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Open(vm, sourceSink, mode, stream, options, k, env)
}

// Open3 is a predicate which opens a stream to a source or sink.
// This predicate is a shorthand for open/4 with an empty list of options.
//
// # Signature
//
//	open(+SourceSink, +Mode, -Stream)
//
// where:
//   - SourceSink is an atom representing the source or sink of the stream, which is typically a URI.
//   - Mode is an atom representing the mode of the stream to be opened. It can be one of "read", "write", or "append".
//   - Stream is the stream to be opened.
//
// open/3 gives True when SourceSink can be opened in Mode.
func Open3(vm *engine.VM, sourceSink, mode, stream engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Call(
		vm,
		atomOpen.Apply(sourceSink, mode, stream, prolog.AtomEmptyList),
		k, env)
}
