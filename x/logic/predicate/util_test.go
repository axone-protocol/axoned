//nolint:lll
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

func TestTermToBytes(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			term        engine.Term
			options     engine.Term
			result      []byte
			wantSuccess bool
			wantError   error
		}{
			{ // If no option, by default, given term is in hexadecimal format.
				term:        engine.NewAtom("486579202120596f752077616e7420746f20736565207468697320746578742c20776f6e64657266756c21"),
				options:     nil,
				result:      []byte{72, 101, 121, 32, 33, 32, 89, 111, 117, 32, 119, 97, 110, 116, 32, 116, 111, 32, 115, 101, 101, 32, 116, 104, 105, 115, 32, 116, 101, 120, 116, 44, 32, 119, 111, 110, 100, 101, 114, 102, 117, 108, 33},
				wantSuccess: true,
			},
			{
				term:        engine.NewAtom("486579202120596f752077616e7420746f20736565207468697320746578742c20776f6e64657266756c21"),
				options:     engine.NewAtom("encoding").Apply(engine.NewAtom("hex")),
				result:      []byte{72, 101, 121, 32, 33, 32, 89, 111, 117, 32, 119, 97, 110, 116, 32, 116, 111, 32, 115, 101, 101, 32, 116, 104, 105, 115, 32, 116, 101, 120, 116, 44, 32, 119, 111, 110, 100, 101, 114, 102, 117, 108, 33},
				wantSuccess: true,
			},
			{
				term:        engine.NewAtom("486579202120596f752077616e7420746f20736565207468697320746578742c20776f6e64657266756c21"),
				options:     engine.NewAtom("encoding").Apply(engine.NewAtom("byte")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("term should be a List, given engine.Atom"),
			},
			{
				term:        engine.List(engine.Integer(72), engine.Integer(101), engine.Integer(121), engine.Integer(32), engine.Integer(33), engine.Integer(32), engine.Integer(89), engine.Integer(111), engine.Integer(117), engine.Integer(32), engine.Integer(119), engine.Integer(97), engine.Integer(110), engine.Integer(116), engine.Integer(32), engine.Integer(116), engine.Integer(111), engine.Integer(32), engine.Integer(115), engine.Integer(101), engine.Integer(101), engine.Integer(32), engine.Integer(116), engine.Integer(104), engine.Integer(105), engine.Integer(115), engine.Integer(32), engine.Integer(116), engine.Integer(101), engine.Integer(120), engine.Integer(116), engine.Integer(44), engine.Integer(32), engine.Integer(119), engine.Integer(111), engine.Integer(110), engine.Integer(100), engine.Integer(101), engine.Integer(114), engine.Integer(102), engine.Integer(117), engine.Integer(108), engine.Integer(33)),
				options:     engine.NewAtom("encoding").Apply(engine.NewAtom("byte")),
				result:      []byte{72, 101, 121, 32, 33, 32, 89, 111, 117, 32, 119, 97, 110, 116, 32, 116, 111, 32, 115, 101, 101, 32, 116, 104, 105, 115, 32, 116, 101, 120, 116, 44, 32, 119, 111, 110, 100, 101, 114, 102, 117, 108, 33},
				wantSuccess: true,
			},
			{
				term:        engine.NewAtom("486579202120596f752077616e7420746f20736565207468697320746578742c20776f6e64657266756c21"),
				options:     engine.NewAtom("encoding").Apply(engine.NewAtom("foo")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("invalid encoding option: foo, valid values are 'hex' or 'byte'"),
			},
			{
				term:        engine.NewAtom("486579202120596f752077616e7420746f20736565207468697320746578742c20776f6e64657266756c21"),
				options:     engine.NewAtom("encoding").Apply(engine.NewAtom("foo"), engine.NewAtom("bar")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("invalid arity for compound 'encoding': 2 but expected 1"),
			},
			{
				term:        engine.NewAtom("486579202120596f752077616e7420746f20736565207468697320746578742c20776f6e64657266756c21"),
				options:     engine.NewAtom("encoding").Apply(engine.NewVariable()),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("invalid term '%%!s(engine.Variable=8)' - expected engine.Atom but got engine.Variable"),
			},
			{
				term:        engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				options:     engine.NewAtom("encoding").Apply(engine.NewAtom("byte")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("term should be a List, given *engine.compound"),
			},
			{
				term:        engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				options:     engine.NewAtom("encoding").Apply(engine.NewAtom("hex")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("invalid term type: *engine.compound, should be an atom"),
			},
			{
				term:        engine.NewAtom("486579202120596f752077616e7420746f20736565207468697320746578742c20776f6e64657266756c21"),
				options:     engine.NewAtom("foo"),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("invalid term 'foo' - expected engine.Compound but got engine.Atom"),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the term #%d: %s", nc, tc.term), func() {
				Convey("when check try convert", func() {
					env := engine.Env{}
					result, err := TermToBytes(tc.term, tc.options, &env)

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
