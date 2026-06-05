package prolog

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAssertIsGround(t *testing.T) {
	X := engine.NewVariable()
	Y := engine.NewVariable()
	foo := engine.NewAtom("foo")
	fortyTwo := engine.Integer(42)

	Convey("Given a test cases", t, func() {
		cases := []struct {
			name     string
			term     engine.Term
			expected error
		}{
			{
				name: "A variable unified",
				term: X,
			},
			{
				name: "an atom",
				term: foo,
			},
			{
				name: "an integer",
				term: fortyTwo,
			},
			{
				name: "a grounded list",
				term: engine.List(foo, X, fortyTwo),
			},
			{
				name: "a grounded compound",
				term: foo.Apply(X, foo.Apply(foo, X, fortyTwo)),
			},
			{
				name:     "a variable",
				term:     Y,
				expected: engine.InstantiationError(engine.NewEnv()),
			},
			{
				name:     "a list containing a variable",
				term:     engine.List(foo, X, Y, fortyTwo),
				expected: engine.InstantiationError(engine.NewEnv()),
			},
			{
				name:     "a compound containing a variable",
				term:     foo.Apply(X, foo.Apply(X, Y, fortyTwo)),
				expected: engine.InstantiationError(engine.NewEnv()),
			},
		}

		Convey("and an environment", func() {
			env, _ := engine.NewEnv().Unify(X, engine.NewAtom("x"))
			for nc, tc := range cases {
				Convey(
					fmt.Sprintf("Given the test case %s (#%d)", tc.name, nc), func() {
						Convey("When the function AreGround() is called", func() {
							result, err := AssertIsGround(tc.term, env)
							Convey("Then it should return the expected output", func() {
								if tc.expected == nil {
									_, err := AssertIsGround(result, env)
									So(err, ShouldBeNil)
								} else {
									So(err, ShouldBeError, tc.expected)
								}
							})
						})
					},
				)
			}
		})
	})
}

func TestAssertKeyValue(t *testing.T) {
	X := engine.NewVariable()

	Convey("Given a test cases", t, func() {
		cases := []struct {
			name      string
			term      engine.Term
			wantKey   engine.Atom
			wantValue engine.Term
			wantError error
		}{
			{
				name:      "a valid key-value pair",
				term:      AtomKeyValue.Apply(StringToAtom("key"), StringToAtom("value")),
				wantKey:   StringToAtom("key"),
				wantValue: StringToAtom("value"),
			},
			{
				name:      "a key-value pair with bounded variable key",
				term:      AtomKeyValue.Apply(X, engine.Integer(42)),
				wantKey:   StringToAtom("x"),
				wantValue: engine.Integer(42),
			},
			{
				name:      "a key-value pair with bounded variable value",
				term:      AtomKeyValue.Apply(StringToAtom("key"), X),
				wantKey:   StringToAtom("key"),
				wantValue: StringToAtom("x"),
			},
			{
				name:      "a key-value pair with non-atom key",
				term:      AtomKeyValue.Apply(engine.Integer(42), StringToAtom("value")),
				wantError: fmt.Errorf("error(type_error(atom,42),root)"),
			},
			{
				name:      "a key-value pair with unbounded variable key",
				term:      AtomKeyValue.Apply(engine.NewVariable(), StringToAtom("value")),
				wantError: fmt.Errorf("error(instantiation_error,root)"),
			},
			{
				name:      "a key-value pair with unbounded variable value",
				term:      AtomKeyValue.Apply(StringToAtom("key"), engine.NewVariable()),
				wantError: fmt.Errorf("error(instantiation_error,root)"),
			},
			{
				name:      "an atom",
				term:      StringToAtom("x"),
				wantError: fmt.Errorf("error(type_error(key_value,x),root)"),
			},
			{
				name:      "an integer",
				term:      engine.Integer(42),
				wantError: fmt.Errorf("error(type_error(key_value,42),root)"),
			},
			{
				name:      "a compound",
				term:      engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				wantError: fmt.Errorf("error(type_error(key_value,foo(bar)),root)"),
			},
			{
				name:      "a key-value pair with arity 1",
				term:      AtomKeyValue.Apply(StringToAtom("key")),
				wantError: fmt.Errorf("error(type_error(key_value,=(key)),root)"),
			},
			{
				name:      "a key-value pair with arity > 2",
				term:      AtomKeyValue.Apply(engine.Integer(1), engine.Integer(2), engine.Integer(3)),
				wantError: fmt.Errorf("error(type_error(key_value,=(1,2,3)),root)"),
			},
		}

		Convey("and an environment", func() {
			env, _ := engine.NewEnv().Unify(X, engine.NewAtom("x"))
			for nc, tc := range cases {
				Convey(
					fmt.Sprintf("Given the test case %s (#%d)", tc.name, nc), func() {
						Convey("When the function AssertKeyValue() is called", func() {
							key, value, err := AssertKeyValue(tc.term, env)
							Convey("Then it should return the expected output", func() {
								if tc.wantError == nil {
									So(key, ShouldEqual, tc.wantKey)
									So(value, ShouldEqual, tc.wantValue)
									So(err, ShouldBeNil)
								} else {
									So(err, ShouldBeError, tc.wantError)
								}
							})
						})
					},
				)
			}
		})
	})
}
