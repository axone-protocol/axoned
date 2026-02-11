package prolog

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestJSONNull(t *testing.T) {
	Convey("Given an environment", t, func() {
		env := engine.NewEnv()
		Convey("When calling JSONNull", func() {
			got := JSONNull()
			want := nullTerm
			Convey("Then it should return the JSON null atom", func() {
				So(got, ShouldNotBeNil)
				So(got.Compare(want, env), ShouldEqual, 0)
			})
		})
	})
}

func TestJSONBool(t *testing.T) {
	Convey("Given a boolean value", t, func() {
		cases := []struct {
			input bool
			want  engine.Term
		}{
			{
				input: true,
				want:  trueTerm,
			},
			{
				input: false,
				want:  falseTerm,
			},
		}

		for _, tc := range cases {
			Convey(fmt.Sprintf("When calling JSONBool(%v)", tc.input), func() {
				got := JSONBool(tc.input)

				Convey("Then the result should be as expected", func() {
					So(got, ShouldEqual, tc.want)
				})
			})
		}
	})
}

func TestAssertJSON(t *testing.T) {
	Convey("Given test cases", t, func() {
		env := engine.NewEnv()
		cases := []struct {
			description string
			input       engine.Term
			wantError   error
			wantResult  string
		}{
			{
				description: "valid JSON object",
				input:       AtomJSON.Apply(engine.NewAtom("valid")),
				wantError:   nil,
			},
			{
				description: "non-compound term",
				input:       engine.NewAtom("notACompound"),
				wantError:   fmt.Errorf("error(type_error(json,notACompound),root)"),
			},
			{
				description: "compound term with arity > 1",
				input:       AtomJSON.Apply(engine.NewAtom("foo"), engine.NewAtom("bar")),
				wantError:   fmt.Errorf("error(type_error(json,json(foo,bar)),root)"),
			},
		}

		for _, tc := range cases {
			Convey(fmt.Sprintf("When calling AssertJSON(%s)", tc.input), func() {
				got, err := AssertJSON(tc.input, env)

				Convey("Then the result should match the expected value", func() {
					if tc.wantError != nil {
						So(got, ShouldBeNil)
						So(err, ShouldBeError, tc.wantError)
					} else {
						So(err, ShouldBeNil)
						So(got, ShouldNotBeNil)
						So(got.Compare(tc.input, env), ShouldEqual, 0)
					}
				})
			})
		}
	})
}
