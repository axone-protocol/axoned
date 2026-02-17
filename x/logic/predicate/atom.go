package predicate

import (
	"fmt"
	"strings"

	"github.com/axone-protocol/prolog/v3/engine"

	"github.com/axone-protocol/axoned/v14/x/logic/prolog"
)

// TermToAtom is a predicate that describes Atom as a term that unifies with Term.
//
// # Signature
//
//	term_to_atom(?Term, ?Atom)
//
// where:
//   - Term is a term that unifies with Atom.
//   - Atom is an atom.
//
// When Atom is instantiated, Atom is parsed and the result unified with Term. If Atom has no valid syntax,
// a syntax_error exception is raised. Otherwise, Term is “written” on Atom using write_term/2 with the option quoted(true).
//
// # Example
//
//	# Convert the atom to a term.
//	- term_to_atom(foo, foo).
func TermToAtom(vm *engine.VM, term, atom engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	switch {
	case prolog.IsGround(term, env):
		var strBuilder strings.Builder
		os := engine.NewOutputTextStream(&strBuilder)
		return engine.WriteTerm(vm, os, term, engine.List(engine.NewAtom("quoted").Apply(prolog.AtomTrue)),
			func(env *engine.Env) *engine.Promise {
				return engine.Unify(vm, prolog.StringToAtom(strBuilder.String()), atom, k, env)
			}, env)
	case prolog.IsGround(atom, env):
		atom, err := prolog.AssertAtom(atom, env)
		if err != nil {
			return engine.Error(err)
		}
		is := engine.NewInputTextStream(strings.NewReader(fmt.Sprintf("%s.", atom)))
		parsedAtom := engine.NewVariable()
		return engine.ReadTerm(vm, is, parsedAtom, engine.List(),
			func(env *engine.Env) *engine.Promise {
				return engine.Unify(vm, term, parsedAtom, k, env)
			}, env)
	}

	return engine.Error(engine.InstantiationError(env))
}

// AtomicListConcat2 is a predicate that unifies an Atom with the concatenated elements of a List.
//
// # Signature
//
//	atomic_list_concat(+List, ?Atom)
//
// where:
//   - List is a list of strings, atoms, integers, floating point numbers or non-integer rationals
//   - Atom is an Atom representing the concatenation of the elements of List
func AtomicListConcat2(vm *engine.VM, list, atom engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	return AtomicListConcat3(vm, list, prolog.AtomEmpty, atom, k, env)
}

// AtomicListConcat3 is a predicate that unifies an Atom with the concatenated elements of a List
// using a given separator.
//
// The atomic_list_concat/3 predicate creates an atom just like atomic_list_concat/2, but inserts Separator
// between each pair of inputs.
//
// # Signature
//
//	atomic_list_concat(+List, +Separator, ?Atom)
//
// where:
//   - List is a list of strings, atoms, integers, floating point numbers or non-integer rationals
//   - Separator is an atom (possibly empty)
//   - Atom is an Atom representing the concatenation of the elements of List
func AtomicListConcat3(vm *engine.VM, list, sep, atom engine.Term, k engine.Cont, env *engine.Env) *engine.Promise {
	if !prolog.IsGround(list, env) {
		return engine.Error(engine.InstantiationError(env))
	}

	it, err := prolog.ListIterator(list, env)
	if err != nil {
		return engine.Error(err)
	}

	if !it.Next() {
		return engine.Unify(vm, prolog.AtomEmpty, atom, k, env)
	}

	if !prolog.IsGround(sep, env) {
		return engine.Error(engine.InstantiationError(env))
	}

	head := engine.NewVariable()
	return TermToAtom(vm, it.Current(), head, func(env *engine.Env) *engine.Promise {
		tail := engine.NewVariable()
		return AtomicListConcat3(vm, it.Suffix(), sep, tail, func(env *engine.Env) *engine.Promise {
			temp := engine.NewVariable()
			if tail.Compare(prolog.AtomEmpty, env) != 0 {
				return engine.AtomConcat(vm, head, sep, temp, func(env *engine.Env) *engine.Promise {
					return engine.AtomConcat(vm, temp, tail, atom, k, env)
				}, env)
			}

			return engine.Unify(vm, atom, head, k, env)
		}, env)
	}, env)
}
