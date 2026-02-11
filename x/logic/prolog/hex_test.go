package prolog

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"

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
				result:      []byte{72, 101, 121, 32, 33, 32, 89, 111, 117, 32, 119, 97, 110, 116, 32, 116, 111, 32, 115, 101, 101, 32, 116, 104, 105, 115, 32, 116, 101, 120, 116, 44, 32, 119, 111, 110, 100, 101, 114, 102, 117, 108, 33}, //nolint:lll
				wantSuccess: true,
			},
			{
				term:        engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				result:      nil,
				wantSuccess: false,
				wantError:   fmt.Errorf("error(type_error(atom,foo(bar)),root)"),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the term #%d: %s", nc, tc.term), func() {
				Convey("when converting hex term to bytes", func() {
					env := engine.NewEnv()
					result, err := TermHexToBytes(tc.term, env)

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
								So(err.Error(), ShouldEqual, tc.wantError.Error())
							})
						})
					}
				})
			})
		}
	})
}
