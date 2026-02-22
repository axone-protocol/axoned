package predicate

import (
	"github.com/axone-protocol/prolog/v3/engine"
)

// Call is a predicate that executes a given goal.
//
// # Signature
//
//	call(+Goal)
//
// Where:
//   - Goal is the goal to execute.
func Call(vm *engine.VM, goal engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Call(vm, goal, cont, env)
}

// Catch is a predicate that catches exceptions thrown during the execution of a goal.
//
// # Signature
//
//	catch(+Goal, ?Catcher, +Recover)
//
// Where:
//   - Goal is the goal to execute.
//   - Catcher is the exception pattern to catch.
//   - Recover is the goal to execute when the exception is caught.
//
//nolint:predeclared,revive
func Catch(vm *engine.VM, goal, catcher, recover engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Catch(vm, goal, catcher, recover, cont, env)
}

// Throw is a predicate that throws an exception.
//
// # Signature
//
//	throw(+Exception)
//
// Where:
//   - Exception is the exception term to throw.
func Throw(vm *engine.VM, exception engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Throw(vm, exception, cont, env)
}

// Unify is a predicate that unifies two terms.
//
// # Signature
//
//	=(+Left, +Right)
//
// Where:
//   - Left is the first term to unify.
//   - Right is the second term to unify.
func Unify(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Unify(vm, left, right, cont, env)
}

// UnifyWithOccursCheck is a predicate that unifies two terms with the occurs check enabled.
//
// # Signature
//
//	unify_with_occurs_check(+Left, +Right)
//
// Where:
//   - Left is the first term to unify.
//   - Right is the second term to unify.
//
// The occurs check prevents the creation of infinite structures by ensuring that a variable
// does not occur in the term it is being unified with.
func UnifyWithOccursCheck(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.UnifyWithOccursCheck(vm, left, right, cont, env)
}

// SubsumesTerm is a predicate that checks if one term subsumes another.
//
// # Signature
//
//	subsumes_term(+General, +Specific)
//
// Where:
//   - General is the general term.
//   - Specific is the specific term.
//
// A term General subsumes Specific if there exists a substitution that makes General identical to Specific.
func SubsumesTerm(vm *engine.VM, general, specific engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.SubsumesTerm(vm, general, specific, cont, env)
}

// TypeVar is a predicate that checks if the given term is an uninstantiated variable.
//
// # Signature
//
//	var(@Term)
//
// Where:
//   - Term is the term to check.
func TypeVar(vm *engine.VM, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.TypeVar(vm, term, cont, env)
}

// TypeAtom is a predicate that checks if the given term is an atom.
//
// # Signature
//
//	atom(@Term)
//
// Where:
//   - Term is the term to check.
func TypeAtom(vm *engine.VM, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.TypeAtom(vm, term, cont, env)
}

// TypeInteger is a predicate that checks if the given term is an integer.
//
// # Signature
//
//	integer(@Term)
//
// Where:
//   - Term is the term to check.
func TypeInteger(vm *engine.VM, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.TypeInteger(vm, term, cont, env)
}

// TypeFloat is a predicate that checks if the given term is a floating-point number.
//
// # Signature
//
//	float(@Term)
//
// Where:
//   - Term is the term to check.
func TypeFloat(vm *engine.VM, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.TypeFloat(vm, term, cont, env)
}

// TypeCompound is a predicate that checks if the given term is a compound term.
//
// # Signature
//
//	compound(@Term)
//
// Where:
//   - Term is the term to check.
//
// A compound term is a functor with one or more arguments.
func TypeCompound(vm *engine.VM, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.TypeCompound(vm, term, cont, env)
}

// AcyclicTerm is a predicate that checks if the given term is acyclic (does not contain cycles).
//
// # Signature
//
//	acyclic_term(@Term)
//
// Where:
//   - Term is the term to check.
func AcyclicTerm(vm *engine.VM, term engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.AcyclicTerm(vm, term, cont, env)
}

// Compare is a predicate that compares two terms according to the standard order.
//
// # Signature
//
//	compare(?Order, +Left, +Right)
//
// Where:
//   - Order is unified with '<', '=', or '>' depending on the comparison result.
//   - Left is the first term to compare.
//   - Right is the second term to compare.
func Compare(vm *engine.VM, order, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Compare(vm, order, left, right, cont, env)
}

// Sort is a predicate that sorts a list of terms, removing duplicates.
//
// # Signature
//
//	sort(+List, ?Sorted)
//
// Where:
//   - List is the input list to sort.
//   - Sorted is the sorted list without duplicates.
func Sort(vm *engine.VM, in, out engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Sort(vm, in, out, cont, env)
}

// KeySort is a predicate that sorts a list of Key-Value pairs by Key.
//
// # Signature
//
//	keysort(+List, ?Sorted)
//
// Where:
//   - List is the input list of Key-Value pairs.
//   - Sorted is the sorted list (duplicates are not removed).
func KeySort(vm *engine.VM, in, out engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.KeySort(vm, in, out, cont, env)
}

// Functor is a predicate that relates a compound term with its functor name and arity.
//
// # Signature
//
//	functor(?Term, ?Name, ?Arity)
//
// Where:
//   - Term is the compound term.
//   - Name is the functor name.
//   - Arity is the number of arguments.
func Functor(vm *engine.VM, term, name, arity engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Functor(vm, term, name, arity, cont, env)
}

// Arg is a predicate that accesses the Nth argument of a compound term.
//
// # Signature
//
//	arg(+N, +Term, ?Arg)
//
// Where:
//   - N is the argument index (1-based).
//   - Term is the compound term.
//   - Arg is the Nth argument of Term.
func Arg(vm *engine.VM, index, term, value engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Arg(vm, index, term, value, cont, env)
}

// Univ is a predicate that converts between a term and a list representation.
//
// # Signature
//
//	=..(+Term, ?List)
//
// Where:
//   - Term is a compound term.
//   - List is a list where the first element is the functor and the rest are the arguments.
func Univ(vm *engine.VM, term, list engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Univ(vm, term, list, cont, env)
}

// CopyTerm is a predicate that creates a copy of a term with fresh variables.
//
// # Signature
//
//	copy_term(+Term, ?Copy)
//
// Where:
//   - Term is the term to copy.
//   - Copy is the copy with renamed variables.
func CopyTerm(vm *engine.VM, in, out engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CopyTerm(vm, in, out, cont, env)
}

// TermVariables is a predicate that collects all variables in a term.
//
// # Signature
//
//	term_variables(+Term, ?Variables)
//
// Where:
//   - Term is the term to analyze.
//   - Variables is a list of all variables in Term.
func TermVariables(vm *engine.VM, term, variables engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.TermVariables(vm, term, variables, cont, env)
}

// Is is a predicate that evaluates an arithmetic expression and unifies the result.
//
// # Signature
//
//	is(?Result, +Expression)
//
// Where:
//   - Result is unified with the evaluated result.
//   - Expression is the arithmetic expression to evaluate.
func Is(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Is(vm, left, right, cont, env)
}

// Equal is a predicate that tests arithmetic equality of two expressions.
//
// # Signature
//
//	=:=(+Left, +Right)
//
// Where:
//   - Left is the first arithmetic expression.
//   - Right is the second arithmetic expression.
func Equal(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Equal(vm, left, right, cont, env)
}

// NotEqual is a predicate that tests arithmetic inequality of two expressions.
//
// # Signature
//
//	=\=(+Left, +Right)
//
// Where:
//   - Left is the first arithmetic expression.
//   - Right is the second arithmetic expression.
func NotEqual(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.NotEqual(vm, left, right, cont, env)
}

// LessThan is a predicate that tests if one arithmetic expression is less than another.
//
// # Signature
//
//	<(+Left, +Right)
//
// Where:
//   - Left is the first arithmetic expression.
//   - Right is the second arithmetic expression.
func LessThan(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.LessThan(vm, left, right, cont, env)
}

// LessThanOrEqual is a predicate that tests if one arithmetic expression is less than or equal to another.
//
// # Signature
//
//	=<(+Left, +Right)
//
// Where:
//   - Left is the first arithmetic expression.
//   - Right is the second arithmetic expression.
func LessThanOrEqual(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.LessThanOrEqual(vm, left, right, cont, env)
}

// GreaterThan is a predicate that tests if one arithmetic expression is greater than another.
//
// # Signature
//
//	>(+Left, +Right)
//
// Where:
//   - Left is the first arithmetic expression.
//   - Right is the second arithmetic expression.
func GreaterThan(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.GreaterThan(vm, left, right, cont, env)
}

// GreaterThanOrEqual is a predicate that tests if one arithmetic expression is greater than or equal to another.
//
// # Signature
//
//	>=(+Left, +Right)
//
// Where:
//   - Left is the first arithmetic expression.
//   - Right is the second arithmetic expression.
func GreaterThanOrEqual(vm *engine.VM, left, right engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.GreaterThanOrEqual(vm, left, right, cont, env)
}

// Clause is a predicate that retrieves clauses from the database.
//
// # Signature
//
//	clause(+Head, ?Body)
//
// Where:
//   - Head is the head of the clause to retrieve.
//   - Body is unified with the body of the clause.
func Clause(vm *engine.VM, head, body engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Clause(vm, head, body, cont, env)
}

// CurrentPredicate is a predicate that enumerates currently defined predicates.
//
// # Signature
//
//	current_predicate(?PredicateIndicator)
//
// Where:
//   - PredicateIndicator is a predicate indicator of the form Name/Arity.
func CurrentPredicate(vm *engine.VM, indicator engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CurrentPredicate(vm, indicator, cont, env)
}

// Asserta is a predicate that asserts a clause into the database as the first clause of the predicate.
//
// # Signature
//
//	asserta(+Clause)
//
// Where:
//   - Clause is the clause to assert into the database.
func Asserta(vm *engine.VM, clause engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Asserta(vm, clause, cont, env)
}

// Assertz is a predicate that asserts a clause into the database as the last clause of the predicate.
//
// # Signature
//
//	assertz(+Clause)
//
// Where:
//   - Clause is the clause to assert into the database.
func Assertz(vm *engine.VM, clause engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Assertz(vm, clause, cont, env)
}

// Retract is a predicate that retracts a clause from the database.
//
// # Signature
//
//	retract(+Clause)
//
// Where:
//   - Clause is the clause to retract from the database.
func Retract(vm *engine.VM, clause engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Retract(vm, clause, cont, env)
}

// Abolish is a predicate that abolishes a predicate from the database.
// Removes all clauses of the predicate designated by given predicate indicator Name/Arity.
//
// # Signature
//
//	abolish(+PredicateIndicator)
//
// Where:
//   - PredicateIndicator is the indicator of the predicate to abolish.
func Abolish(vm *engine.VM, indicator engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Abolish(vm, indicator, cont, env)
}

// FindAll is a predicate that collects all solutions to a goal.
//
// # Signature
//
//	findall(?Template, +Goal, ?Bag)
//
// Where:
//   - Template is the term to collect for each solution.
//   - Goal is the goal to find solutions for.
//   - Bag is unified with the list of collected solutions.
func FindAll(vm *engine.VM, template, goal, bag engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.FindAll(vm, template, goal, bag, cont, env)
}

// BagOf is a predicate that collects solutions to a goal grouped by free variables.
//
// # Signature
//
//	bagof(?Template, +Goal, ?Bag)
//
// Where:
//   - Template is the term to collect for each solution.
//   - Goal is the goal to find solutions for.
//   - Bag is unified with the list of collected solutions.
//
// Unlike findall/3, bagof/3 fails if there are no solutions and groups solutions by free variables.
func BagOf(vm *engine.VM, template, goal, bag engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.BagOf(vm, template, goal, bag, cont, env)
}

// SetOf is a predicate that collects unique solutions to a goal in sorted order.
//
// # Signature
//
//	setof(?Template, +Goal, ?Set)
//
// Where:
//   - Template is the term to collect for each solution.
//   - Goal is the goal to find solutions for.
//   - Set is unified with the sorted list of unique solutions.
func SetOf(vm *engine.VM, template, goal, sorted engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.SetOf(vm, template, goal, sorted, cont, env)
}

// CurrentInput is a predicate that unifies the given term with the current input stream.
//
// # Signature
//
//	current_input(?Stream)
//
// Where:
//   - Stream represents the current input stream.
func CurrentInput(vm *engine.VM, input engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CurrentInput(vm, input, cont, env)
}

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

// SetInput is a predicate that sets the current input stream.
//
// # Signature
//
//	set_input(+Stream)
//
// Where:
//   - Stream is the stream to set as current input.
func SetInput(vm *engine.VM, input engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.SetInput(vm, input, cont, env)
}

// SetOutput is a predicate that sets the current output stream.
//
// # Signature
//
//	set_output(+Stream)
//
// Where:
//   - Stream is the stream to set as current output.
func SetOutput(vm *engine.VM, output engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.SetOutput(vm, output, cont, env)
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

// Close is a predicate that closes a stream.
//
// # Signature
//
//	close(+Stream, +Options)
//
// Where:
//   - Stream is the stream to close.
//   - Options is a list of options for closing the stream.
func Close(vm *engine.VM, stream, options engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Close(vm, stream, options, cont, env)
}

// FlushOutput is a predicate that flushes the output of a stream.
//
// # Signature
//
//	flush_output(+Stream)
//
// Where:
//   - Stream is the stream to flush.
func FlushOutput(vm *engine.VM, stream engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.FlushOutput(vm, stream, cont, env)
}

// StreamProperty is a predicate that queries properties of a stream.
//
// # Signature
//
//	stream_property(?Stream, ?Property)
//
// Where:
//   - Stream is the stream to query.
//   - Property is a property of the stream.
func StreamProperty(vm *engine.VM, stream, property engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.StreamProperty(vm, stream, property, cont, env)
}

// SetStreamPosition is a predicate that sets the position of a stream.
//
// # Signature
//
//	set_stream_position(+Stream, +Position)
//
// Where:
//   - Stream is the stream to reposition.
//   - Position is the new position in the stream.
func SetStreamPosition(vm *engine.VM, stream, position engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.SetStreamPosition(vm, stream, position, cont, env)
}

// GetChar is a predicate that reads a character from a stream.
//
// # Signature
//
//	get_char(+Stream, ?Char)
//
// Where:
//   - Stream is the stream to read from.
//   - Char is unified with the character read.
func GetChar(vm *engine.VM, stream, char engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.GetChar(vm, stream, char, cont, env)
}

// PeekChar is a predicate that peeks at the next character from a stream without consuming it.
//
// # Signature
//
//	peek_char(+Stream, ?Char)
//
// Where:
//   - Stream is the stream to peek from.
//   - Char is unified with the next character.
func PeekChar(vm *engine.VM, stream, char engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.PeekChar(vm, stream, char, cont, env)
}

// PutChar is a predicate that writes a character to a stream.
//
// # Signature
//
//	put_char(+Stream, +Char)
//
// Where:
//   - Stream is the stream to write to.
//   - Char is the character to write.
func PutChar(vm *engine.VM, stream, char engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.PutChar(vm, stream, char, cont, env)
}

// GetByte is a predicate that reads a byte from a stream.
//
// # Signature
//
//	get_byte(+Stream, ?Byte)
//
// Where:
//   - Stream is the stream to read from.
//   - Byte is unified with the byte read (0-255).
func GetByte(vm *engine.VM, stream, b engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.GetByte(vm, stream, b, cont, env)
}

// PeekByte is a predicate that peeks at the next byte from a stream without consuming it.
//
// # Signature
//
//	peek_byte(+Stream, ?Byte)
//
// Where:
//   - Stream is the stream to peek from.
//   - Byte is unified with the next byte (0-255).
func PeekByte(vm *engine.VM, stream, b engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.PeekByte(vm, stream, b, cont, env)
}

// PutByte is a predicate that writes a byte to a stream.
//
// # Signature
//
//	put_byte(+Stream, +Byte)
//
// Where:
//   - Stream is the stream to write to.
//   - Byte is the byte to write (0-255).
func PutByte(vm *engine.VM, stream, b engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.PutByte(vm, stream, b, cont, env)
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

// Op is a predicate that declares an operator with a given priority and associativity.
//
// # Signature
//
//	op(+Priority, +Specifier, +Name)
//
// Where:
//   - Priority is an integer between 1 and 1200.
//   - Specifier is the operator type (fx, fy, xf, yf, xfx, xfy, yfx, yfy).
//   - Name is the operator name.
func Op(vm *engine.VM, priority, specifier, operator engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Op(vm, priority, specifier, operator, cont, env)
}

// CurrentOp is a predicate that queries currently defined operators.
//
// # Signature
//
//	current_op(?Priority, ?Specifier, ?Name)
//
// Where:
//   - Priority is the operator priority.
//   - Specifier is the operator type.
//   - Name is the operator name.
func CurrentOp(vm *engine.VM, priority, specifier, operator engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CurrentOp(vm, priority, specifier, operator, cont, env)
}

// CharConversion is a predicate that defines a character conversion rule.
//
// # Signature
//
//	char_conversion(+From, +To)
//
// Where:
//   - From is the source character.
//   - To is the target character.
func CharConversion(vm *engine.VM, from, to engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CharConversion(vm, from, to, cont, env)
}

// CurrentCharConversion is a predicate that queries currently defined character conversions.
//
// # Signature
//
//	current_char_conversion(?From, ?To)
//
// Where:
//   - From is the source character.
//   - To is the target character.
func CurrentCharConversion(vm *engine.VM, from, to engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CurrentCharConversion(vm, from, to, cont, env)
}

// Negate is a predicate that succeeds if the goal fails (negation by failure).
//
// # Signature
//
//	\+(+Goal)
//
// Where:
//   - Goal is the goal to negate.
func Negate(vm *engine.VM, goal engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Negate(vm, goal, cont, env)
}

// Repeat is a predicate that succeeds indefinitely on backtracking.
//
// # Signature
//
//	repeat
//
// This predicate always succeeds and provides an infinite number of choice points.
func Repeat(vm *engine.VM, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Repeat(vm, cont, env)
}

// Call1 is a predicate that calls a goal with one additional argument.
//
// # Signature
//
//	call(+Goal, +Arg1)
//
// Where:
//   - Goal is the goal to call.
//   - Arg1 is an additional argument to append to Goal.
func Call1(vm *engine.VM, goal, arg1 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Call1(vm, goal, arg1, cont, env)
}

// Call2 is a predicate that calls a goal with two additional arguments.
//
// # Signature
//
//	call(+Goal, +Arg1, +Arg2)
//
// Where:
//   - Goal is the goal to call.
//   - Arg1, Arg2 are additional arguments to append to Goal.
func Call2(vm *engine.VM, goal, arg1, arg2 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Call2(vm, goal, arg1, arg2, cont, env)
}

// Call3 is a predicate that calls a goal with three additional arguments.
//
// # Signature
//
//	call(+Goal, +Arg1, +Arg2, +Arg3)
//
// Where:
//   - Goal is the goal to call.
//   - Arg1, Arg2, Arg3 are additional arguments to append to Goal.
func Call3(vm *engine.VM, goal, arg1, arg2, arg3 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Call3(vm, goal, arg1, arg2, arg3, cont, env)
}

// Call4 is a predicate that calls a goal with four additional arguments.
//
// # Signature
//
//	call(+Goal, +Arg1, +Arg2, +Arg3, +Arg4)
//
// Where:
//   - Goal is the goal to call.
//   - Arg1, Arg2, Arg3, Arg4 are additional arguments to append to Goal.
func Call4(vm *engine.VM, goal, arg1, arg2, arg3, arg4 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Call4(vm, goal, arg1, arg2, arg3, arg4, cont, env)
}

// Call5 is a predicate that calls a goal with five additional arguments.
//
// # Signature
//
//	call(+Goal, +Arg1, +Arg2, +Arg3, +Arg4, +Arg5)
//
// Where:
//   - Goal is the goal to call.
//   - Arg1, Arg2, Arg3, Arg4, Arg5 are additional arguments to append to Goal.
func Call5(vm *engine.VM, goal, arg1, arg2, arg3, arg4, arg5 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Call5(vm, goal, arg1, arg2, arg3, arg4, arg5, cont, env)
}

// Call6 is a predicate that calls a goal with six additional arguments.
//
// # Signature
//
//	call(+Goal, +Arg1, +Arg2, +Arg3, +Arg4, +Arg5, +Arg6)
//
// Where:
//   - Goal is the goal to call.
//   - Arg1, Arg2, Arg3, Arg4, Arg5, Arg6 are additional arguments to append to Goal.
func Call6(
	vm *engine.VM, goal, arg1, arg2, arg3, arg4, arg5, arg6 engine.Term, cont engine.Cont, env *engine.Env,
) *engine.Promise {
	return engine.Call6(vm, goal, arg1, arg2, arg3, arg4, arg5, arg6, cont, env)
}

// Call7 is a predicate that calls a goal with seven additional arguments.
//
// # Signature
//
//	call(+Goal, +Arg1, +Arg2, +Arg3, +Arg4, +Arg5, +Arg6, +Arg7)
//
// Where:
//   - Goal is the goal to call.
//   - Arg1, Arg2, Arg3, Arg4, Arg5, Arg6, Arg7 are additional arguments to append to Goal.
func Call7(
	vm *engine.VM, goal, arg1, arg2, arg3, arg4, arg5, arg6, arg7 engine.Term, cont engine.Cont, env *engine.Env,
) *engine.Promise {
	return engine.Call7(vm, goal, arg1, arg2, arg3, arg4, arg5, arg6, arg7, cont, env)
}

// AtomLength is a predicate that determines the length of an atom.
//
// # Signature
//
//	atom_length(+Atom, ?Length)
//
// Where:
//   - Atom is the atom whose length to determine.
//   - Length is unified with the number of characters in Atom.
func AtomLength(vm *engine.VM, atom, length engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.AtomLength(vm, atom, length, cont, env)
}

// AtomConcat is a predicate that concatenates atoms.
//
// # Signature
//
//	atom_concat(?Atom1, ?Atom2, ?Atom3)
//
// Where:
//   - Atom1 is the first atom.
//   - Atom2 is the second atom.
//   - Atom3 is the concatenation of Atom1 and Atom2.
func AtomConcat(vm *engine.VM, left, right, atom engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.AtomConcat(vm, left, right, atom, cont, env)
}

// SubAtom is a predicate that extracts or checks sub-atoms.
//
// # Signature
//
//	sub_atom(+Atom, ?Before, ?Length, ?After, ?SubAtom)
//
// Where:
//   - Atom is the source atom.
//   - Before is the number of characters before the sub-atom.
//   - Length is the length of the sub-atom.
//   - After is the number of characters after the sub-atom.
//   - SubAtom is the extracted sub-atom.
func SubAtom(
	vm *engine.VM, atom, before, length, after, sub engine.Term, cont engine.Cont, env *engine.Env,
) *engine.Promise {
	return engine.SubAtom(vm, atom, before, length, after, sub, cont, env)
}

// AtomChars is a predicate that converts between an atom and a list of characters.
//
// # Signature
//
//	atom_chars(?Atom, ?Chars)
//
// Where:
//   - Atom is the atom.
//   - Chars is a list of single-character atoms.
func AtomChars(vm *engine.VM, atom, chars engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.AtomChars(vm, atom, chars, cont, env)
}

// AtomCodes is a predicate that converts between an atom and a list of character codes.
//
// # Signature
//
//	atom_codes(?Atom, ?Codes)
//
// Where:
//   - Atom is the atom.
//   - Codes is a list of character codes (integers).
func AtomCodes(vm *engine.VM, atom, codes engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.AtomCodes(vm, atom, codes, cont, env)
}

// CharCode is a predicate that converts between a character and its character code.
//
// # Signature
//
//	char_code(?Char, ?Code)
//
// Where:
//   - Char is a single-character atom.
//   - Code is the corresponding character code (integer).
func CharCode(vm *engine.VM, char, code engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CharCode(vm, char, code, cont, env)
}

// NumberChars is a predicate that converts between a number and a list of characters.
//
// # Signature
//
//	number_chars(?Number, ?Chars)
//
// Where:
//   - Number is the number.
//   - Chars is a list of characters representing the number.
func NumberChars(vm *engine.VM, number, chars engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.NumberChars(vm, number, chars, cont, env)
}

// NumberCodes is a predicate that converts between a number and a list of character codes.
//
// # Signature
//
//	number_codes(?Number, ?Codes)
//
// Where:
//   - Number is the number.
//   - Codes is a list of character codes representing the number.
func NumberCodes(vm *engine.VM, number, codes engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.NumberCodes(vm, number, codes, cont, env)
}

// CurrentPrologFlag is a predicate that queries Prolog flags.
//
// # Signature
//
//	current_prolog_flag(?Flag, ?Value)
//
// Where:
//   - Flag is the name of the flag.
//   - Value is the current value of the flag.
func CurrentPrologFlag(vm *engine.VM, flag, value engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CurrentPrologFlag(vm, flag, value, cont, env)
}

// Phrase is a predicate that parses a list using a grammar rule (DCG).
//
// # Signature
//
//	phrase(+GrammarRule, ?Input, ?Rest)
//
// Where:
//   - GrammarRule is the grammar rule to apply.
//   - Input is the input list to parse.
//   - Rest is the remaining unparsed part of the input.
func Phrase(vm *engine.VM, grammarRule, input, rest engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Phrase(vm, grammarRule, input, rest, cont, env)
}

// ExpandTerm is a predicate that expands a term (e.g., DCG rules).
//
// # Signature
//
//	expand_term(+Term, ?Expanded)
//
// Where:
//   - Term is the term to expand.
//   - Expanded is the expanded form of the term.
func ExpandTerm(vm *engine.VM, term, expanded engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.ExpandTerm(vm, term, expanded, cont, env)
}

// Append is a predicate that concatenates two lists.
//
// # Signature
//
//	append(?List1, ?List2, ?List3)
//
// Where:
//   - List1 is the first list.
//   - List2 is the second list.
//   - List3 is the concatenation of List1 and List2.
func Append(vm *engine.VM, left, right, list engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Append(vm, left, right, list, cont, env)
}

// Length is a predicate that determines the length of a list.
//
// # Signature
//
//	length(?List, ?Length)
//
// Where:
//   - List is the list.
//   - Length is the number of elements in the list.
func Length(vm *engine.VM, list, length engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Length(vm, list, length, cont, env)
}

// Between is a predicate that generates or tests integers within a range.
//
// # Signature
//
//	between(+Low, +High, ?Value)
//
// Where:
//   - Low is the lower bound (inclusive).
//   - High is the upper bound (inclusive).
//   - Value is unified with integers between Low and High.
func Between(vm *engine.VM, low, high, value engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Between(vm, low, high, value, cont, env)
}

// Succ is a predicate that relates consecutive integers.
//
// # Signature
//
//	succ(?Int1, ?Int2)
//
// Where:
//   - Int1 is an integer.
//   - Int2 is Int1 + 1.
func Succ(vm *engine.VM, predecessor, successor engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Succ(vm, predecessor, successor, cont, env)
}

// Nth0 is a predicate that accesses the Nth element of a list (0-indexed).
//
// # Signature
//
//	nth0(?N, ?List, ?Elem)
//
// Where:
//   - N is the index (starting from 0).
//   - List is the list.
//   - Elem is the element at position N.
func Nth0(vm *engine.VM, index, list, value engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Nth0(vm, index, list, value, cont, env)
}

// Nth1 is a predicate that accesses the Nth element of a list (1-indexed).
//
// # Signature
//
//	nth1(?N, ?List, ?Elem)
//
// Where:
//   - N is the index (starting from 1).
//   - List is the list.
//   - Elem is the element at position N.
func Nth1(vm *engine.VM, index, list, value engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Nth1(vm, index, list, value, cont, env)
}

// CallNth is a predicate that succeeds for the Nth solution of a goal.
//
// # Signature
//
//	call_nth(+Goal, ?N)
//
// Where:
//   - Goal is the goal to execute.
//   - N is unified with the solution number (starting from 1).
func CallNth(vm *engine.VM, goal, n engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.CallNth(vm, goal, n, cont, env)
}

// Op3 is a predicate that accesses values in a dictionary or option list.
//
// # Signature
//
//	.(+Key, ?Value, +Dict)
//
// Where:
//   - Key is the key to look up.
//   - Value is the associated value.
//   - Dict is the dictionary or option list.
func Op3(vm *engine.VM, key, value, list engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Op3(vm, key, value, list, cont, env)
}

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
