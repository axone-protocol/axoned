package predicate

import "github.com/axone-protocol/prolog/engine"

// CurrentOutput is a predicate that unifies the given term with the current output stream.
//
// # Signature
//
//	current_output(-Stream) is det
//
// where:
//   - Stream represents the current output stream.
//
// This predicate connects to the default output stream available for user interactions, allowing the user to perform
// write operations.
//
// The outcome of the stream's content throughout the execution of a query is provided as a string within the
// user_output field in the query's response. However, it's important to note that the maximum length of the output
// is constrained by the max_query_output_size setting, meaning only the final max_query_output_size bytes (not characters)
// of the output are included in the response.
func CurrentOutput(vm *engine.VM, output engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CurrentOutput(vm, output, cont, env)
}
