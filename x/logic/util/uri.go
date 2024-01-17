package util

import (
	"net/url"
)

type URIComponent int

const (
	PathComponent URIComponent = iota
	SegmentComponent
	QueryValueComponent
	FragmentComponent
)

const upperhex = "0123456789ABCDEF"

// Return true if the specified character should be escaped when
// appearing in a URL string depending on the targeted URI component, according
// to [RFC 3986](https://www.rfc-editor.org/rfc/rfc3986).
//
// This is a re-implementation of url.shouldEscape of net/url. Needed since the native implementation doesn't follow
// exactly the [RFC 3986](https://www.rfc-editor.org/rfc/rfc3986) and also because the implementation of component
// escaping is only public for Path component (who in reality is SegmentPath component) and Query component. Otherwise,
// escaping doesn't fit to the SWI-Prolog escaping due to RFC discrepancy between those two implementations.
//
// Another discrepancy is on the query component that escape the space character ' ' to a '+' (plus sign) on the
// golang library and to '%20' escaping on the
// [SWI-Prolog implementation](https://www.swi-prolog.org/pldoc/doc/_SWI_/library/uri.pl?show=src#uri_encoded/3).
//
// Here some reported issues on golang about the RFC non-compliance.
//   - golang.org/issue/5684.
//   - https://github.com/golang/go/issues/27559
func (c URIComponent) shouldEscape(b byte) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= b && b <= 'z' || 'A' <= b && b <= 'Z' || '0' <= b && b <= '9' {
		return false
	}

	switch b {
	case '-', '.', '_', '~': // §2.3 Unreserved characters (mark)
		return false

	case '!', '$', '&', '\'', '(', ')', '*', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch c {
		case PathComponent: // §3.3
			return b == '?' || b == ':'
		case SegmentComponent: // §3.3
			// The RFC allows : @ & = + $
			// meaning to individual path segments.
			return b == '/' || b == '?' || b == ':'
		case QueryValueComponent: // §3.4
			return b == '&' || b == '+' || b == ':' || b == ';' || b == '='
		case FragmentComponent: // §4.1
			return false
		}
	}

	// Everything else must be escaped.
	return true
}

// Escape return the given input string by adding percent encoding depending on the current component where it's
// supposed to be put.
// This is a re-implementation of native url.escape. See shouldEscape() comment's for more details.
func (c URIComponent) Escape(v string) string {
	hexCount := 0
	for i := 0; i < len(v); i++ {
		ch := v[i]
		if c.shouldEscape(ch) {
			hexCount++
		}
	}

	if hexCount == 0 {
		return v
	}

	var buf [64]byte
	var t []byte

	required := len(v) + 2*hexCount
	if required <= len(buf) {
		t = buf[:required]
	} else {
		t = make([]byte, required)
	}

	j := 0
	for i := 0; i < len(v); i++ {
		switch ch := v[i]; {
		case c.shouldEscape(ch):
			t[j] = '%'
			t[j+1] = upperhex[ch>>4]
			t[j+2] = upperhex[ch&15]
			j += 3
		default:
			t[j] = v[i]
			j++
		}
	}
	return string(t)
}

func (c URIComponent) Unescape(v string) (string, error) {
	return url.PathUnescape(v)
}
