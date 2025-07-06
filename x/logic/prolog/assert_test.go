package prolog

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v2/engine"
	"github.com/samber/lo"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/axone-protocol/axoned/v12/x/logic/testutil"
	"github.com/axone-protocol/axoned/v12/x/logic/util"
)

var (
	predicateMatches = PredicateMatches
	urlMatches       = func(this string) func(string) bool {
		return func(that string) bool {
			return util.URLMatches(util.ParseURLMust(this))(util.ParseURLMust(that))
		}
	}
)

func TestWhitelistBlacklistMatches(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			values     []string
			whitelist  []string
			blacklist  []string
			predicate  func(string) func(string) bool
			wantResult []string
		}{ // predicate filter test cases
			{
				values:     []string{},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{"length/2", "member/2", "call/1", "call/2", "member/2"},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "call/1", "length/2", "member/2"},
				whitelist:  []string{"length/2", "member/2", "call/2", "member/2"},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "member/2", "call/1"},
			},
			{
				values:     []string{},
				whitelist:  []string{},
				blacklist:  []string{"call/1"},
				predicate:  predicateMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{},
				blacklist:  []string{"call/2"},
				predicate:  predicateMatches,
				wantResult: []string{"length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{"call/2", "length/2", "member/2"},
				blacklist:  []string{"call/1", "member/1", "findall"},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{"call/1"},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{"call"},
				predicate:  predicateMatches,
				wantResult: []string{"member/2"},
			},
			// url filter test cases
			{
				values:     []string{},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"https://www.axone.network"},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"https://www.axone.network"},
			},
			{
				values:     []string{"https://www.axone.network", "https://www.axone.xyz/foo/bar?baz=qux#frag"},
				whitelist:  []string{"https://www.axone.xyz/foo/bar?baz=qux#frag"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"https://www.axone.xyz/foo/bar?baz=qux#frag"},
			},
			{
				values:     []string{"https://www.axone.network", "https://www.axone.xyz"},
				whitelist:  []string{"https://www.axone.xyz"},
				blacklist:  []string{"https://www.axone.xyz"},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"http://example.com/foo/bar"},
				whitelist:  []string{"http://example.com/foo"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"http://example.com/foo"},
				whitelist:  []string{"http://example.com/foo/"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"http://example.com/foo"},
				whitelist:  []string{"http://example.com/foo"},
				blacklist:  []string{"http://example.com/foo?"},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"mailto:user@example.com"},
				whitelist:  []string{"mailto:user@example.com"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"mailto:user@example.com"},
			},
			{
				values:     []string{"tel:123456789"},
				whitelist:  []string{"tel:"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"tel:123456789"},
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d with values: %v checked against whitelist: %v and blacklist: %v",
					nc, tc.values, tc.whitelist, tc.blacklist), func() {
					Convey("When the function WhitelistBlacklistMatches() is called", func() {
						result := lo.Filter(tc.values, util.Indexed(util.WhitelistBlacklistMatches(tc.whitelist, tc.blacklist, tc.predicate)))

						Convey(fmt.Sprintf("Then it should return the expected output: %v", tc.wantResult), func() {
							So(result, ShouldResemble, tc.wantResult)
						})
					})
				})
		}
	})
}

func TestAreGround(t *testing.T) {
	X := engine.NewVariable()
	Y := engine.NewVariable()
	foo := engine.NewAtom("foo")
	fortyTwo := engine.Integer(42)

	Convey("Given a test cases", t, func() {
		cases := []struct {
			name     string
			terms    []engine.Term
			expected bool
		}{
			{
				name:     "all terms are ground",
				terms:    []engine.Term{X, foo, foo.Apply(X), fortyTwo, engine.List(X, fortyTwo)},
				expected: true,
			},
			{
				name:     "one term is a variable",
				terms:    []engine.Term{X, foo, Y, foo.Apply(X)},
				expected: false,
			},
			{
				name:     "one term is a list containing a variable",
				terms:    []engine.Term{X, foo, engine.List(X, Y, foo), fortyTwo},
				expected: false,
			},
			{
				name:     "one term is a compound containing a variable",
				terms:    []engine.Term{X, foo, foo.Apply(X, foo.Apply(X, Y, fortyTwo)), fortyTwo},
				expected: false,
			},
			{
				name:     "no terms",
				terms:    []engine.Term{},
				expected: true,
			},
			{
				name:     "no terms (2)",
				terms:    []engine.Term{AtomEmptyList},
				expected: true,
			},
		}

		Convey("and an environment", func() {
			env, _ := engine.NewEnv().Unify(X, engine.NewAtom("x"))
			for nc, tc := range cases {
				Convey(
					fmt.Sprintf("Given the test case %s (#%d)", tc.name, nc), func() {
						Convey("When the function AreGround() is called", func() {
							result := AreGround(tc.terms, env)

							Convey("Then it should return the expected output", func() {
								So(result, ShouldEqual, tc.expected)
							})
						})
					})
			}
		})
	})
}

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
									So(result, testutil.ShouldBeGrounded)
								} else {
									So(err, ShouldBeError, tc.expected)
								}
							})
						})
					})
			}
		})
	})
}

func TestAssertPair(t *testing.T) {
	X := engine.NewVariable()

	Convey("Given a test cases", t, func() {
		cases := []struct {
			name       string
			term       engine.Term
			wantFirst  engine.Term
			wantSecond engine.Term
			wantError  error
		}{
			{
				name:       "a valid pair",
				term:       AtomPair.Apply(StringToAtom("foo"), StringToAtom("bar")),
				wantFirst:  StringToAtom("foo"),
				wantSecond: StringToAtom("bar"),
			},
			{
				name:       "a pair with bounded variable",
				term:       AtomPair.Apply(X, engine.Integer(42)),
				wantFirst:  StringToAtom("x"),
				wantSecond: engine.Integer(42),
			},
			{
				name:      "a pair with unbounded variable",
				term:      AtomPair.Apply(engine.NewVariable(), StringToAtom("bar")),
				wantError: fmt.Errorf("error(instantiation_error,root)"),
			},
			{
				name:      "an atom",
				term:      StringToAtom("x"),
				wantError: fmt.Errorf("error(type_error(pair,x),root)"),
			},
			{
				name:      "an integer",
				term:      engine.Integer(42),
				wantError: fmt.Errorf("error(type_error(pair,42),root)"),
			},
			{
				name:      "a compound",
				term:      engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				wantError: fmt.Errorf("error(type_error(pair,foo(bar)),root)"),
			},
			{
				name:      "a pair with arity 1",
				term:      AtomPair.Apply(StringToAtom("foo")),
				wantError: fmt.Errorf("error(type_error(pair,-(foo)),root)"),
			},
			{
				name:      "a pair with arity > 1",
				term:      AtomPair.Apply(engine.Integer(1), engine.Integer(2), engine.Integer(3)),
				wantError: fmt.Errorf("error(type_error(pair,-(1,2,3)),root)"),
			},
		}

		Convey("and an environment", func() {
			env, _ := engine.NewEnv().Unify(X, engine.NewAtom("x"))
			for nc, tc := range cases {
				Convey(
					fmt.Sprintf("Given the test case %s (#%d)", tc.name, nc), func() {
						Convey("When the function AssertPair() is called", func() {
							first, second, err := AssertPair(tc.term, env)
							Convey("Then it should return the expected output", func() {
								if tc.wantError == nil {
									So(first, ShouldEqual, tc.wantFirst)
									So(second, ShouldEqual, tc.wantSecond)
									So(err, ShouldBeNil)
								} else {
									So(err, ShouldBeError, tc.wantError)
								}
							})
						})
					})
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
					})
			}
		})
	})
}
