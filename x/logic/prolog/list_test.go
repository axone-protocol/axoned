package prolog

import (
	"errors"
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestForEach(t *testing.T) {
	Convey("Given test cases", t, func(c C) {
		env := engine.NewEnv()
		cases := []struct {
			description string
			list        engine.Term
			f           func(v engine.Term, hasNext bool) error
			wantError   error
		}{
			{
				description: "Empty list",
				list:        engine.NewAtom("[]"),

				f: func(_ engine.Term, _ bool) error {
					t.Errorf("Function should not be called for empty list")
					return nil
				},
			},
			{
				description: "Non-list term",
				list:        engine.NewAtom("not_a_list"),

				f: func(_ engine.Term, _ bool) error {
					t.Errorf("Function should not be called for non-list term")
					return nil
				},
				wantError: fmt.Errorf("error(type_error(list,not_a_list),root)"),
			},
			{
				description: "List with elements",
				list:        engine.List(engine.NewAtom("a"), engine.NewAtom("b"), engine.NewAtom("c")),

				f: func() func(v engine.Term, hasNext bool) error {
					i := 0
					values := []string{"a", "b", "c"}
					return func(v engine.Term, hasNext bool) error {
						defer func() { i++ }()

						c.So(i, ShouldBeLessThan, len(values))
						got, err := AssertAtom(v, env)
						c.So(err, ShouldBeNil)
						c.So(got.String(), ShouldEqual, values[i])
						c.So(hasNext, ShouldEqual, i < len(values)-1)

						return nil
					}
				}(),
				wantError: nil,
			},
			{
				description: "Function returns error",
				list:        engine.List(engine.NewAtom("a"), engine.NewAtom("b")),

				f: func() func(v engine.Term, hasNext bool) error {
					i := 0
					return func(_ engine.Term, hasNext bool) error {
						defer func() { i++ }()
						c.So(i, ShouldEqual, 0)
						c.So(hasNext, ShouldBeTrue)

						return errors.New("test error")
					}
				}(),
				wantError: fmt.Errorf("test error"),
			},
		}

		for tn, tc := range cases {
			Convey(fmt.Sprintf("When calling ForEach (case %d)", tn), func() {
				err := ForEach(tc.list, env, tc.f)
				if tc.wantError != nil {
					So(err, ShouldBeError, tc.wantError)
				} else {
					So(err, ShouldBeNil)
				}
			})
		}
	})
}
