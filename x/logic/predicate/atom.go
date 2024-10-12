package predicate

import (
	"fmt"
	"strings"

	"github.com/axone-protocol/prolog/engine"

	"github.com/axone-protocol/axoned/v10/x/logic/prolog"
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
	if !prolog.IsGround(list, env) {
		return engine.Error(engine.InstantiationError(env))
	}

	it, err := prolog.ListIterator(list, env)
	if err != nil {
		return engine.Error(err)
	}

	switch {
	case !it.Next():
		return engine.Unify(vm, prolog.AtomEmpty, atom, k, env)
	default:
		headAtom := engine.NewVariable()
		return TermToAtom(vm, it.Current(), headAtom, func(env *engine.Env) *engine.Promise {
			tailAtom := engine.NewVariable()
			return AtomicListConcat2(vm, it.Suffix(), tailAtom, func(env *engine.Env) *engine.Promise {
				return engine.AtomConcat(vm, headAtom, tailAtom, atom, k, env)
			}, env)
		}, env)
	}
}
