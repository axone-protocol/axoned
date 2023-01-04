package interpreter

import (
	"github.com/ichiban/prolog/engine"
)

type Callback[T any] func() T

// Instrument1 is a higher order function that given a 1arg-predicate and a callback returns a new predicate that calls the
// callback before calling the predicate.
func Instrument1[T any](callback Callback[T], p engine.Predicate1) engine.Predicate1 {
	return func(vm *engine.VM, t1 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		callback()

		return p(vm, t1, cont, env)
	}
}

// Instrument2 is a higher order function that given a 2args-predicate and a callback returns a new predicate that calls the
// callback before calling the predicate.
func Instrument2[T any](callback Callback[T], p engine.Predicate2) engine.Predicate2 {
	return func(vm *engine.VM, t1 engine.Term, t2 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		callback()

		return p(vm, t1, t2, cont, env)
	}
}

// Instrument3 is a higher order function that given a 3args-predicate and a callback returns a new predicate that calls the
// callback before calling the predicate.
func Instrument3[T any](callback Callback[T], p engine.Predicate3) engine.Predicate3 {
	return func(vm *engine.VM, t1 engine.Term, t2 engine.Term, t3 engine.Term, cont engine.Cont,
		env *engine.Env,
	) *engine.Promise {
		callback()

		return p(vm, t1, t2, t3, cont, env)
	}
}

// Instrument4 is a higher order function that given a 4args-predicate and a callback returns a new predicate that calls the
// callback before calling the predicate.
//
//nolint:lll
func Instrument4[T any](callback Callback[T], p engine.Predicate4) engine.Predicate4 {
	return func(vm *engine.VM, t1 engine.Term, t2 engine.Term, t3 engine.Term, t4 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		callback()

		return p(vm, t1, t2, t3, t4, cont, env)
	}
}

// Instrument5 is a higher order function that given a 5args-predicate and a callback returns a new predicate that calls the
// callback before calling the predicate.
//
//nolint:lll
func Instrument5[T any](callback Callback[T], p engine.Predicate5) engine.Predicate5 {
	return func(vm *engine.VM, t1 engine.Term, t2 engine.Term, t3 engine.Term, t4 engine.Term, t5 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		callback()

		return p(vm, t1, t2, t3, t4, t5, cont, env)
	}
}

// Instrument6 is a higher order function that given a 6args-predicate and a callback returns a new predicate that calls the
// callback before calling the predicate.
//
//nolint:lll
func Instrument6[T any](callback Callback[T], p engine.Predicate6) engine.Predicate6 {
	return func(vm *engine.VM, t1 engine.Term, t2 engine.Term, t3 engine.Term, t4 engine.Term, t5 engine.Term, t6 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		callback()

		return p(vm, t1, t2, t3, t4, t5, t6, cont, env)
	}
}

// Instrument7 is a higher order function that given a 7args-predicate and a callback returns a new predicate that calls the
// callback before calling the predicate.
//
//nolint:lll
func Instrument7[T any](callback Callback[T], p engine.Predicate7) engine.Predicate7 {
	return func(vm *engine.VM, t1 engine.Term, t2 engine.Term, t3 engine.Term, t4 engine.Term, t5 engine.Term, t6 engine.Term, t7 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		callback()

		return p(vm, t1, t2, t3, t4, t5, t6, t7, cont, env)
	}
}

// Instrument8 is a higher order function that given a 8args-predicate and a callback returns a new predicate that calls the
// callback before calling the predicate.
//
//nolint:lll
func Instrument8[T any](callback Callback[T], p engine.Predicate8) engine.Predicate8 {
	return func(vm *engine.VM, t1 engine.Term, t2 engine.Term, t3 engine.Term, t4 engine.Term, t5 engine.Term, t6 engine.Term, t7 engine.Term, t8 engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
		callback()

		return p(vm, t1, t2, t3, t4, t5, t6, t7, t8, cont, env)
	}
}
