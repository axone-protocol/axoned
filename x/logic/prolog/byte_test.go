package prolog

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBytesToAtom(t *testing.T) {
	Convey("Given the BytesToAtom function", t, func() {
		Convey("It should correctly convert byte slices to atoms", func() {
			cases := []struct {
				bytes []byte
				want  engine.Atom
			}{
				{
					bytes: []byte(""),
					want:  engine.NewAtom(""),
				},
				{
					bytes: []byte("foo bar"),
					want:  engine.NewAtom("foo bar"),
				},
				{
					bytes: []byte("こんにちは"),
					want:  engine.NewAtom("こんにちは"),
				},
				{
					bytes: []byte{0xF0, 0x9F, 0x98, 0x80},
					want:  engine.NewAtom("😀"),
				},
			}

			for _, tc := range cases {
				Convey(fmt.Sprintf("When converting '%s", tc.bytes), func() {
					got := BytesToAtom(tc.bytes)

					Convey("Then the result should match the expected value", func() {
						So(got, ShouldEqual, tc.want)
					})
				})
			}
		})
	})
}
