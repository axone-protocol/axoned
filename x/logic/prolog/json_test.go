package prolog

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
