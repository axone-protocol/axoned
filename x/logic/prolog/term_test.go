//nolint:lll
package prolog

import (
	"fmt"
	"testing"

	"github.com/ichiban/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTermHexToBytes(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			term        engine.Term
			result      []byte
			wantSuccess bool
			wantError   error
		}{
			{ // If no option, by default, given term is in hexadecimal format.
				term:        engine.NewAtom("486579202120596f752077616e7420746f20736565207468697320746578742c20776f6e64657266756c21"),
				result:      []byte{72, 101, 121, 32, 33, 32, 89, 111, 117, 32, 119, 97, 110, 116, 32, 116, 111, 32, 115, 101, 101, 32, 116, 104, 105, 115, 32, 116, 101, 120, 116, 44, 32, 119, 111, 110, 100, 101, 114, 102, 117, 108, 33},
				wantSuccess: true,
			},
			{
				term:        engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(atom,foo(bar)),_2)"),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the term #%d: %s", nc, tc.term), func() {
				Convey("when converting hex term to bytes", func() {
					env := engine.Env{}
					result, err := TermHexToBytes(tc.term, &env)

					if tc.wantSuccess {
						Convey("then no error should be thrown", func() {
							So(err, ShouldBeNil)

							Convey("and result should be as expected", func() {
								So(result, ShouldResemble, tc.result)
							})
						})
					} else {
						Convey("then error should occurs", func() {
							So(err, ShouldNotEqual, nil)

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

func TestTermToBytes(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			term        engine.Term
			encoding    string
			result      []byte
			wantSuccess bool
			wantError   error
		}{
			{
				term:        engine.NewAtom("foo"),
				result:      []byte{102, 111, 111},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.Integer(72), engine.Integer(101), engine.Integer(121), engine.Integer(32), engine.Integer(33), engine.Integer(32), engine.Integer(89), engine.Integer(111), engine.Integer(117), engine.Integer(32), engine.Integer(119), engine.Integer(97), engine.Integer(110), engine.Integer(116), engine.Integer(32), engine.Integer(116), engine.Integer(111), engine.Integer(32), engine.Integer(115), engine.Integer(101), engine.Integer(101), engine.Integer(32), engine.Integer(116), engine.Integer(104), engine.Integer(105), engine.Integer(115), engine.Integer(32), engine.Integer(116), engine.Integer(101), engine.Integer(120), engine.Integer(116), engine.Integer(44), engine.Integer(32), engine.Integer(119), engine.Integer(111), engine.Integer(110), engine.Integer(100), engine.Integer(101), engine.Integer(114), engine.Integer(102), engine.Integer(117), engine.Integer(108), engine.Integer(33)),
				result:      []byte{72, 101, 121, 32, 33, 32, 89, 111, 117, 32, 119, 97, 110, 116, 32, 116, 111, 32, 115, 101, 101, 32, 116, 104, 105, 115, 32, 116, 101, 120, 116, 44, 32, 119, 111, 110, 100, 101, 114, 102, 117, 108, 33},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.NewAtom("f"), engine.NewAtom("o"), engine.NewAtom("o")),
				result:      []byte{102, 111, 111},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.NewAtom("ü")),
				result:      []byte{195, 188},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.NewAtom("ü")),
				encoding:    "utf-8",
				result:      []byte{195, 188},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.NewAtom("ü")),
				encoding:    "octet",
				result:      []byte{252},
				wantSuccess: true,
			},
			{
				term:        engine.NewAtom("ツ"),
				encoding:    "utf8",
				result:      []byte{227, 131, 132},
				wantSuccess: true,
			},
			{
				term:        engine.NewAtom("ツ"),
				encoding:    "text",
				result:      []byte{227, 131, 132},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.NewAtom("ツ")),
				encoding:    "utf8",
				result:      []byte{227, 131, 132},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.Integer(227), engine.Integer(131), engine.Integer(132)),
				encoding:    "utf8",
				result:      []byte{227, 131, 132},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.NewAtom("ツ")),
				encoding:    "shift-jis",
				result:      []byte{131, 99},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.Integer(227), engine.Integer(131), engine.Integer(132)),
				encoding:    "shift-jis",
				result:      []byte{131, 99},
				wantSuccess: true,
			},
			{
				term:        engine.List(engine.NewAtom("ツ")),
				encoding:    "octet",
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(domain_error(valid_byte(12484),[227,131,132]),_3)"),
			},
			{
				term:        engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(character_code,foo(bar)),_4)"),
			},
			{
				term:        engine.List(engine.NewAtom("f"), engine.NewAtom("oo")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(domain_error(valid_character_code(oo),[f,oo]),_5)"),
			},
			{
				term:        engine.NewAtom("foo"),
				encoding:    "foo",
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(domain_error(valid_charset,foo),_6)"),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the term #%d: %s", nc, tc.term), func() {
				Convey("when converting string term to bytes", func() {
					env := engine.Env{}
					result, err := StringTermToBytes(tc.term, engine.NewAtom(tc.encoding), &env)

					if tc.wantSuccess {
						Convey("then no error should be thrown", func() {
							So(err, ShouldEqual, nil)

							Convey("and result should be as expected", func() {
								So(result, ShouldResemble, tc.result)
							})
						})
					} else {
						Convey("then error should occurs", func() {
							So(err, ShouldNotEqual, nil)

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
