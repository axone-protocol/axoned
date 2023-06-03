package types

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestToSubstitutions(t *testing.T) {
	Convey("Given test cases", t, func() {
		cases := []struct {
			term       TermResults
			wantResult []Substitution
		}{
			{
				term:       TermResults{},
				wantResult: []Substitution{},
			},
			{
				term: TermResults{
					"X": "foo",
				},
				wantResult: []Substitution{
					{
						Variable: "X",
						Term: Term{
							Name: "foo",
						},
					},
				},
			},
			{
				term: TermResults{
					"X": "foo",
					"Y": "bar",
				},
				wantResult: []Substitution{
					{
						Variable: "X",
						Term: Term{
							Name: "foo",
						},
					},
					{
						Variable: "Y",
						Term: Term{
							Name: "bar",
						},
					},
				},
			},
			{
				term: TermResults{
					"Y": "bar",
					"X": "foo",
				},
				wantResult: []Substitution{
					{
						Variable: "X",
						Term: Term{
							Name: "foo",
						},
					},
					{
						Variable: "Y",
						Term: Term{
							Name: "bar",
						},
					},
				},
			},
		}

		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given test case #%d", nc), func() {
				Convey("When ToSubstitutions() function is called", func() {
					substitutions := tc.term.ToSubstitutions()

					Convey("Then the result should match expectations", func() {
						So(substitutions, ShouldResemble, tc.wantResult)
					})
				})
			})
		}
	})
}
