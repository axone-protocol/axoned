package keeper

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFilterPredicates(t *testing.T) {
	Convey("Given a test cases", t, func() {
		cases := []struct {
			registry   []string
			whitelist  []string
			blacklist  []string
			wantResult []string
		}{
			{
				registry:   []string{},
				whitelist:  []string{},
				blacklist:  []string{},
				wantResult: []string{},
			},
			{
				registry:   []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{},
				blacklist:  []string{},
				wantResult: []string{},
			},
			{
				registry:   []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{"length/2", "member/2", "call/1", "call/2", "member/2"},
				blacklist:  []string{},
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				registry:   []string{"call/2", "call/1", "length/2", "member/2"},
				whitelist:  []string{"length/2", "member/2", "call/2", "member/2"},
				blacklist:  []string{},
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				registry:   []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{},
				wantResult: []string{"call/2", "member/2", "call/1"},
			},
			{
				registry:   []string{},
				whitelist:  []string{},
				blacklist:  []string{"call/1"},
				wantResult: []string{},
			},
			{
				registry:   []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{},
				blacklist:  []string{"call/2"},
				wantResult: []string{},
			},
			{
				registry:   []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{"call/2", "length/2", "member/2"},
				blacklist:  []string{"call/1", "member/1", "findall"},
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				registry:   []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{"call/1"},
				wantResult: []string{"call/2", "member/2"},
			},
			{
				registry:   []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{"call"},
				wantResult: []string{"member/2"},
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d with registry: %v, whitelist: %v, and blacklist: %v",
					nc, tc.registry, tc.whitelist, tc.blacklist), func() {
					Convey("When the function filterPredicates() is called", func() {
						result := lo.Filter(tc.registry, filterPredicates(tc.whitelist, tc.blacklist))

						Convey(fmt.Sprintf("Then it should return the expected output: %v", tc.wantResult), func() {
							So(result, ShouldResemble, tc.wantResult)
						})
					})
				})
		}
	})
}
