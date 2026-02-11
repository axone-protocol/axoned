package keeper

import (
	"fmt"
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStringifyOperand(t *testing.T) {
	Convey("Given various inputs to stringifyOperand", t, func() {
		testCases := []struct {
			description string
			input       engine.Term
			expected    string
			ok          bool
		}{
			{
				description: "an operand implementing fmt.Stringer",
				input:       engine.NewAtom("foo"),
				expected:    "foo",
				ok:          true,
			},
			{
				description: "an operand not implementing fmt.Stringer",
				input:       engine.NewAtom("foo").Apply(engine.NewAtom("bar")),
				expected:    "",
				ok:          false,
			},
			{
				description: "the nil operand",
				input:       nil,
				expected:    "",
				ok:          false,
			},
		}

		for _, tc := range testCases {
			Convey(fmt.Sprintf("When input is %s", tc.description), func() {
				result, ok := stringifyOperand(tc.input)

				Convey("Then the result should match the expected output", func() {
					So(result, ShouldEqual, tc.expected)
					So(ok, ShouldEqual, tc.ok)
				})
			})
		}
	})
}
