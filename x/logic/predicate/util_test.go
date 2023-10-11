package predicate

import (
	"fmt"
	"testing"

	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestExtractJsonTerm(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			compound    engine.Compound
			result      map[string]engine.Term
			wantSuccess bool
			wantError   error
		}{
			{
				compound:    engine.NewAtom("foo").Apply(engine.NewAtom("bar")).(engine.Compound),
				wantSuccess: false,
				wantError:   fmt.Errorf("invalid functor foo. Expected json"),
			},
			{
				compound:    engine.NewAtom("json").Apply(engine.NewAtom("bar"), engine.NewAtom("foobar")).(engine.Compound),
				wantSuccess: false,
				wantError:   fmt.Errorf("invalid compound arity : 2 but expected 1"),
			},
			{
				compound:    engine.NewAtom("json").Apply(engine.NewAtom("bar")).(engine.Compound),
				wantSuccess: false,
				wantError:   fmt.Errorf("json compound should contains one list, give engine.Atom"),
			},
			{
				compound: AtomJSON.Apply(engine.List(AtomPair.Apply(engine.NewAtom("foo"), engine.NewAtom("bar")))).(engine.Compound),
				result: map[string]engine.Term{
					"foo": engine.NewAtom("bar"),
				},
				wantSuccess: true,
			},
			{
				compound:    AtomJSON.Apply(engine.List(engine.NewAtom("foo"), engine.NewAtom("bar"))).(engine.Compound),
				wantSuccess: false,
				wantError:   fmt.Errorf("json attributes should be a pair"),
			},
			{
				compound:    AtomJSON.Apply(engine.List(AtomPair.Apply(engine.Integer(10), engine.NewAtom("bar")))).(engine.Compound),
				wantSuccess: false,
				wantError:   fmt.Errorf("first pair arg should be an atom"),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the term compound #%d: %s", nc, tc.compound), func() {
				Convey("when extract json term", func() {
					env := engine.Env{}
					result, err := ExtractJSONTerm(tc.compound, &env)

					if tc.wantSuccess {
						Convey("then no error should be thrown", func() {
							So(err, ShouldBeNil)
							So(result, ShouldNotBeNil)

							Convey("and result should be as expected", func() {
								So(result, ShouldResemble, tc.result)
							})
						})
					} else {
						Convey("then error should occurs", func() {
							So(err, ShouldNotBeNil)

							Convey("and should be as expected", func() {
								So(err, ShouldResemble, tc.wantError)
							})
						})
					}
				})
			})
		}
	})
}

func TestOptionsContains(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			atom        engine.Atom
			options     engine.Term
			result      engine.Compound
			wantSuccess bool
			wantError   error
		}{
			{
				atom:        engine.NewAtom("foo"),
				options:     engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				result:      engine.NewAtom("foo").Apply(engine.NewAtom("bar")).(engine.Compound),
				wantSuccess: true,
			},
			{
				atom:        engine.NewAtom("bar"),
				options:     engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				result:      nil,
				wantSuccess: true,
			},
			{
				atom:        engine.NewAtom("foo"),
				options:     engine.List(engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				result:      engine.NewAtom("foo").Apply(engine.NewAtom("bar")).(engine.Compound),
				wantSuccess: true,
			},
			{
				atom:        engine.NewAtom("bar"),
				options:     engine.List(engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				result:      nil,
				wantSuccess: true,
			},
			{
				atom: engine.NewAtom("foo"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey").Apply(engine.NewAtom("hoo")),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				result:      engine.NewAtom("foo").Apply(engine.NewAtom("bar")).(engine.Compound),
				wantSuccess: true,
			},
			{
				atom: engine.NewAtom("hey"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey").Apply(engine.NewAtom("hoo")),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				result:      engine.NewAtom("hey").Apply(engine.NewAtom("hoo")).(engine.Compound),
				wantSuccess: true,
			},
			{
				atom: engine.NewAtom("foo"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey"),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantSuccess: false,
				wantError:   fmt.Errorf("invalid options term, should be compound, give engine.Atom"),
			},
			{
				atom:        engine.NewAtom("foo"),
				options:     nil,
				wantSuccess: true,
				result:      nil,
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the term option #%d: %s", nc, tc.atom), func() {
				Convey("when check contains", func() {
					env := engine.Env{}
					result, err := OptionsContains(tc.atom, tc.options, &env)

					if tc.wantSuccess {
						Convey("then no error should be thrown", func() {
							So(err, ShouldBeNil)

							Convey("and result should be as expected", func() {
								So(result, ShouldResemble, tc.result)
							})
						})
					} else {
						Convey("then error should occurs", func() {
							So(err, ShouldNotBeNil)

							Convey("and should be as expected", func() {
								So(err, ShouldResemble, tc.wantError)
							})
						})
					}
				})
			})
		}
	})
}
