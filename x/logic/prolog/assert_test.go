package prolog

import (
	"fmt"
	"testing"

	"github.com/ichiban/prolog/engine"
	"github.com/samber/lo"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/okp4/okp4d/x/logic/testutil"
	"github.com/okp4/okp4d/x/logic/util"
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
				values:     []string{"https://www.okp4.network"},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"https://www.okp4.network"},
			},
			{
				values:     []string{"https://www.okp4.network", "https://www.okp4.com/foo/bar?baz=qux#frag"},
				whitelist:  []string{"https://www.okp4.com/foo/bar?baz=qux#frag"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"https://www.okp4.com/foo/bar?baz=qux#frag"},
			},
			{
				values:     []string{"https://www.okp4.network", "https://www.okp4.com"},
				whitelist:  []string{"https://www.okp4.com"},
				blacklist:  []string{"https://www.okp4.com"},
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
				fmt.Sprintf("Given test case #%d with values: %v cheked against whitelist: %v and blacklist: %v",
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
