package predicate

import "github.com/axone-protocol/prolog/engine"

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
