package util

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	predicateMatches = PredicateMatches
	urlMatches       = func(this string) func(string) bool {
		return func(that string) bool {
			return UrlMatches(ParseUrlMust(this))(ParseUrlMust(that))
		}
	}
)

func TestWhitelistBlacklistMatches(t *testing.T) {

	Convey("Given a test cases", t, func() {
		cases := []struct {
			values     []string
			whitelist  []string
			blacklist  []string
			predicate  func(string) func(string) bool
			wantResult []string
		}{ // predicate filter test cases
			{
				values:     []string{},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{"length/2", "member/2", "call/1", "call/2", "member/2"},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "call/1", "length/2", "member/2"},
				whitelist:  []string{"length/2", "member/2", "call/2", "member/2"},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "member/2", "call/1"},
			},
			{
				values:     []string{},
				whitelist:  []string{},
				blacklist:  []string{"call/1"},
				predicate:  predicateMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{},
				blacklist:  []string{"call/2"},
				predicate:  predicateMatches,
				wantResult: []string{"length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/2", "member/2"},
				whitelist:  []string{"call/2", "length/2", "member/2"},
				blacklist:  []string{"call/1", "member/1", "findall"},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "length/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{"call/1"},
				predicate:  predicateMatches,
				wantResult: []string{"call/2", "member/2"},
			},
			{
				values:     []string{"call/2", "length/1", "member/2", "call/1"},
				whitelist:  []string{"length/2", "member/2", "call", "member/2"},
				blacklist:  []string{"call"},
				predicate:  predicateMatches,
				wantResult: []string{"member/2"},
			},
			// url filter test cases
			{
				values:     []string{},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"https://www.okp4.network"},
				whitelist:  []string{},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"https://www.okp4.network"},
			},
			{
				values:     []string{"https://www.okp4.network", "https://www.okp4.com/foo/bar?baz=qux#frag"},
				whitelist:  []string{"https://www.okp4.com/foo/bar?baz=qux#frag"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"https://www.okp4.com/foo/bar?baz=qux#frag"},
			},
			{
				values:     []string{"https://www.okp4.network", "https://www.okp4.com"},
				whitelist:  []string{"https://www.okp4.com"},
				blacklist:  []string{"https://www.okp4.com"},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"http://example.com/foo/bar"},
				whitelist:  []string{"http://example.com/foo"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"http://example.com/foo"},
				whitelist:  []string{"http://example.com/foo/"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"http://example.com/foo"},
				whitelist:  []string{"http://example.com/foo"},
				blacklist:  []string{"http://example.com/foo?"},
				predicate:  urlMatches,
				wantResult: []string{},
			},
			{
				values:     []string{"mailto:user@example.com"},
				whitelist:  []string{"mailto:user@example.com"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"mailto:user@example.com"},
			},
			{
				values:     []string{"tel:123456789"},
				whitelist:  []string{"tel:"},
				blacklist:  []string{},
				predicate:  urlMatches,
				wantResult: []string{"tel:123456789"},
			},
		}

		for nc, tc := range cases {
			Convey(
				fmt.Sprintf("Given test case #%d with values: %v cheked against whitelist: %v and blacklist: %v",
					nc, tc.values, tc.whitelist, tc.blacklist), func() {
					Convey("When the function WhitelistBlacklistMatches() is called", func() {
						result := lo.Filter(tc.values, Indexed(WhitelistBlacklistMatches(tc.whitelist, tc.blacklist, tc.predicate)))

						Convey(fmt.Sprintf("Then it should return the expected output: %v", tc.wantResult), func() {
							So(result, ShouldResemble, tc.wantResult)
						})
					})
				})
		}
	})
}
