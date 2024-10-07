package prolog

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetOption(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			option     engine.Atom
			options    engine.Term
			wantResult engine.Term
			wantError  error
		}{
			{
				option:     engine.NewAtom("foo"),
				options:    nil,
				wantResult: nil,
				wantError:  nil,
			},
			{
				option:     engine.NewAtom("foo"),
				options:    engine.List(),
				wantResult: nil,
				wantError:  nil,
			},
			{
				option:     engine.NewAtom("foo"),
				options:    engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				wantResult: engine.NewAtom("bar"),
				wantError:  nil,
			},
			{
				option:     engine.NewAtom("bar"),
				options:    engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				wantResult: nil,
				wantError:  nil,
			},
			{
				option:     engine.NewAtom("foo"),
				options:    engine.List(engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: engine.NewAtom("bar"),
				wantError:  nil,
			},
			{
				option:     engine.NewAtom("bar"),
				options:    engine.List(engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: nil,
				wantError:  nil,
			},
			{
				option: engine.NewAtom("foo"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey").Apply(engine.NewAtom("hoo")),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: engine.NewAtom("bar"),
				wantError:  nil,
			},
			{
				option: engine.NewAtom("foo"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar1")),
					engine.NewAtom("hey").Apply(engine.NewAtom("hoo")),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar1"))),
				wantResult: engine.NewAtom("bar1"),
				wantError:  nil,
			},
			{
				option: engine.NewAtom("hey"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey").Apply(engine.NewAtom("hoo")),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: engine.NewAtom("hoo"),
				wantError:  nil,
			},
			{
				option: engine.NewAtom("hey"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey").Apply(engine.NewAtom("jo").Apply(engine.NewAtom("bi"))),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
				wantError:  nil,
			},
			{
				option: engine.NewAtom("hey"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey").Apply(engine.NewAtom("jo").Apply(engine.NewAtom("bi"))),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
				wantError:  nil,
			},
			{
				option: engine.NewAtom("hey"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey").Apply(engine.List(engine.NewAtom("bi"), engine.NewAtom("bar"))),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: engine.List(engine.NewAtom("bi"), engine.NewAtom("bar")),
				wantError:  nil,
			},
			{
				option: engine.NewAtom("hey"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.List(engine.NewAtom("hey").Apply(engine.NewAtom("joe"))),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: nil,
				wantError:  nil,
			},
			{
				option:     engine.NewAtom("foo"),
				options:    engine.NewAtom("foo"),
				wantResult: nil,
				wantError:  fmt.Errorf("error(type_error(option,foo),root)"),
			},
			{
				option: engine.NewAtom("foo"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey"),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar"))),
				wantResult: nil,
				wantError:  fmt.Errorf("error(type_error(option,hey),root)"),
			},
			{
				option: engine.NewAtom("foo"),
				options: engine.List(
					engine.NewAtom("jo").Apply(engine.NewAtom("bi")),
					engine.NewAtom("hey").Apply(engine.NewAtom("hoo")),
					engine.NewAtom("foo").Apply(engine.NewAtom("bar1"), engine.NewAtom("bar2"))),
				wantResult: nil,
				wantError:  fmt.Errorf("error(type_error(option,foo(bar1,bar2)),root)"),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the term option #%d: %s", nc, tc.option), func() {
				Convey("when getting option", func() {
					env := engine.NewEnv()
					result, err := GetOption(tc.option, tc.options, env)

					if tc.wantError == nil {
						Convey("then no error should be thrown", func() {
							So(err, ShouldBeNil)

							Convey("and result should be as expected", func() {
								So(result, ShouldEqual, tc.wantResult)
							})
						})
					} else {
						Convey("then atom returned should be the empty one", func() {
							So(result, ShouldEqual, tc.wantResult)
						})
						Convey("then error should occurs", func() {
							So(err, ShouldNotBeNil)

							Convey("and should be as expected", func() {
								So(err.Error(), ShouldEqual, tc.wantError.Error())
							})
						})
					}
				})
			})
		}
	})
}
