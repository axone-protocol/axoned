package util

import (
	"net/url"
)

// URLMatches is a function that returns a function that matches the given url against the given other item.
//
// The function matches the components of the given url against the components of the given other url. If the component
// of the given other url is empty, it is considered to match the component of the given url.
// For example:
//   - URLMatches("http://example.com/foo")("http://example.com/foo") -> true
//   - URLMatches("http://example.com/foo")("http://example.com/foo?bar=baz") -> false
//   - URLMatches("tel:123456789")("tel:") -> true
//
// The function is curried, and is a binary relation that is reflexive, associative (but not commutative).
func URLMatches(this *url.URL) func(*url.URL) bool {
	return func(that *url.URL) bool {
		return (that.Scheme == "" || that.Scheme == this.Scheme) &&
			(that.Opaque == "" || that.Opaque == this.Opaque) &&
			(that.User == nil || that.User.String() == "" || that.User.String() == this.User.String()) &&
			(that.Host == "" || that.Host == this.Host) &&
			(that.Path == "" || that.Path == this.Path) &&
			(that.RawQuery == "" || that.RawQuery == this.RawQuery) &&
			(that.Fragment == "" || that.Fragment == this.Fragment)
	}
}

// ParseURLMust parses the given url and panics if it fails.
// You have been warned.
func ParseURLMust(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}
