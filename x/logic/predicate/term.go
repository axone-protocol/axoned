package predicate

import "github.com/axone-protocol/prolog/engine"

// WriteTerm3 is a predicate that writes a term to a stream or alias.
//
// The signature is as follows:
//
//	write_term(+Stream, +Term, +Options)
//
// where:
//   - Stream represents the stream or alias to write the term to.
//   - Term represents the term to write.
//   - Options represents the options to control the writing process.
//
// Valid options are:
//
//   - quoted(Bool): If true, atoms and strings that need quotes will be quoted. The default is false.
//   - ignore_ops(Bool): If true, the generic term representation (<functor>(<args> ... )) will be used for all terms.
//     Otherwise (default), operators will be used where appropriate.
//   - numbervars(Bool): If true, variables will be numbered. The default is false.
//   - variable_names(+List): Assign names to variables in Term. List is a list of Name = Var terms, where Name is an atom
//     and Var is a variable.
//   - max_depth(+Int): The maximum depth to which the term is written. The default is infinite.
func WriteTerm3(vm *engine.VM, streamOrAlias, t, options engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	return engine.WriteTerm(vm, streamOrAlias, t, options, k, env)
}

// ReadTerm3 is a predicate that reads a term from a stream or alias.
//
// The signature is as follows:
//
//	read_term(+Stream, -Term, +Options)
//
// where:
//   - Stream represents the stream or alias to read the term from.
//   - Term represents the term to read.
//   - Options represents the options to control the reading process.
//
// Valid options are:
//
//   - singletons(Vars): Vars is unified with a list of variables that occur only once in the term.
//   - variables(Vars): Vars is unified with a list of variables that occur in the term.
//   - variable_names(Vars): Vars is unified with a list of Name = Var terms, where Name is an atom and Var is a variable.
func ReadTerm3(vm *engine.VM, streamOrAlias, t, options engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	return engine.ReadTerm(vm, streamOrAlias, t, options, k, env)
}
