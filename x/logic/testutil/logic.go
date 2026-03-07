package testutil

import (
	"context"
	"fmt"
	"regexp"

	"github.com/axone-protocol/prolog/v3"
	"github.com/axone-protocol/prolog/v3/engine"

	logicembeddedfs "github.com/axone-protocol/axoned/v14/x/logic/fs/embedded"
	logicvfs "github.com/axone-protocol/axoned/v14/x/logic/fs/vfs"
	"github.com/axone-protocol/axoned/v14/x/logic/interpreter/bootstrap"
	logiclib "github.com/axone-protocol/axoned/v14/x/logic/lib"
)

type TermResults map[string]prolog.TermString

// NewLightInterpreterMust returns a new Interpreter with the given context or panics if it fails.
// The Interpreter is configured with minimal settings to support testing.
func NewLightInterpreterMust(ctx context.Context) (i *prolog.Interpreter) {
	i = &prolog.Interpreter{}
	i.Register3(engine.NewAtom("op"), engine.Op)
	i.Register3(engine.NewAtom("compare"), engine.Compare)
	i.Register2(engine.NewAtom("="), engine.Unify)
	i.Register1(engine.NewAtom("consult"), engine.Consult)

	err := i.Compile(ctx, `
						:-(op(1200, xfx, ':-')).
						:-(op(1000, xfy, ',')).
						:-(op(700, xfx, [==, \==, @<, @=<, @>, @>=])).
						:-(op(700, xfx, '=')).
						:-(op(500, yfx, [+, -, /\, \/])).

						member(X, [X|_]).
						member(X, [_|Xs]) :- member(X, Xs).
						X == Y :- compare(=, X, Y).`)
	if err != nil {
		panic(err)
	}

	return
}

// NewComprehensiveInterpreterMust returns a new Interpreter with the given context or panics if it fails.
// The Interpreter is configured with the full bootstrap but with a minimal set of predicates.
func NewComprehensiveInterpreterMust(ctx context.Context) *prolog.Interpreter {
	i := &prolog.Interpreter{}
	i.Register3(engine.NewAtom("op"), engine.Op)
	i.Register3(engine.NewAtom("compare"), engine.Compare)
	i.Register2(engine.NewAtom("="), engine.Unify)
	i.Register1(engine.NewAtom("consult"), engine.Consult)
	i.Register3(engine.NewAtom("bagof"), engine.BagOf)
	i.Register1(engine.NewAtom("current_output"), engine.CurrentOutput)
	i.Register1(engine.NewAtom("current_input"), engine.CurrentInput)
	i.Register2(engine.NewAtom("put_char"), engine.PutChar)
	i.Register2(engine.NewAtom("get_char"), engine.GetChar)
	i.Register3(engine.NewAtom("write_term"), engine.WriteTerm)
	i.Register1(engine.NewAtom("call"), engine.Call)
	i.Register1(engine.NewAtom("\\+"), engine.Negate)
	i.Register1(engine.NewAtom("var"), engine.TypeVar)
	i.Register1(engine.NewAtom("atom"), engine.TypeAtom)
	i.Register1(engine.NewAtom("integer"), engine.TypeInteger)
	i.Register1(engine.NewAtom("float"), engine.TypeFloat)
	i.Register1(engine.NewAtom("compound"), engine.TypeCompound)
	i.Register2(engine.NewAtom("term_variables"), engine.TermVariables)
	i.Register3(engine.NewAtom("catch"), engine.Catch)
	i.Register1(engine.NewAtom("throw"), engine.Throw)
	i.Register2(engine.NewAtom(">="), engine.GreaterThanOrEqual)
	i.Register2(engine.NewAtom("=<"), engine.LessThanOrEqual)
	i.Register2(engine.NewAtom(">"), engine.GreaterThan)
	i.Register2(engine.NewAtom("<"), engine.LessThan)
	i.Register2(engine.NewAtom("is"), engine.Is)
	i.Register2(engine.NewAtom("atom_chars"), engine.AtomChars)
	i.Register2(engine.NewAtom("char_code"), engine.CharCode)

	err := i.Compile(ctx, bootstrap.Bootstrap())
	if err != nil {
		panic(err)
	}

	pathFS := logicvfs.New()
	if err := pathFS.Mount("/v1/lib", logicembeddedfs.NewFS(logiclib.Files)); err != nil {
		panic(fmt.Errorf("failed to mount /v1/lib: %w", err))
	}
	i.FS = pathFS

	return i
}

// CompileMust compiles the given source code and panics if it fails.
// This is a convenience function for testing.
func CompileMust(ctx context.Context, i *prolog.Interpreter, s string, args ...any) {
	err := i.Compile(ctx, s, args...)
	if err != nil {
		panic(err)
	}
}

// ReindexUnknownVariables reindexes the variables in the given term so that the variables are numbered sequentially.
// This is required for test predictability when the term is a result of a query and the variables are unknown.
//
// For example, the following term:
//
//	foo(_1, _2, _3, _1)
//
// is re-indexed as:
//
//	foo(_1, _2, _3, _4)
func ReindexUnknownVariables(s prolog.TermString) prolog.TermString {
	re := regexp.MustCompile("_([0-9]+)")
	var index int
	return prolog.TermString(re.ReplaceAllStringFunc(string(s), func(_ string) string {
		index++
		return fmt.Sprintf("_%d", index)
	}))
}

// ShouldBeGrounded is a goconvey assertion that asserts that the given term does not hold any
// uninstantiated variables.
func ShouldBeGrounded(actual any, expected ...any) string {
	if len(expected) != 0 {
		return fmt.Sprintf("This assertion requires exactly %d comparison values (you provided %d).", 0, len(expected))
	}

	var containsVariable func(engine.Term) bool
	containsVariable = func(term engine.Term) bool {
		switch t := term.(type) {
		case engine.Variable:
			return true
		case engine.Compound:
			for i := 0; i < t.Arity(); i++ {
				if containsVariable(t.Arg(i)) {
					return true
				}
			}
		}
		return false
	}
	if t, ok := actual.(engine.Term); ok {
		if containsVariable(t) {
			return "Expected term to NOT hold a free variable (but it was)."
		}

		return ""
	}

	return fmt.Sprintf("The argument to this assertion must be a term (you provided %v).", actual)
}
