package predicate

import (
	"context"
	"fmt"
	"net/url"

	"github.com/ichiban/prolog/engine"

	"github.com/okp4/okp4d/x/logic/prolog"
)

type Component string

const (
	QueryComponent    Component = "query"
	FragmentComponent Component = "fragment"
	PathComponent     Component = "path"
	SegmentComponent  Component = "segment"
)

const upperhex = "0123456789ABCDEF"

func NewComponent(v string) (Component, error) {
	switch v {
	case string(QueryComponent):
		return QueryComponent, nil
	case string(FragmentComponent):
		return FragmentComponent, nil
	case string(PathComponent):
		return PathComponent, nil
	case string(SegmentComponent):
		return SegmentComponent, nil
	default:
		return "", fmt.Errorf("invalid component name %s, expected `query`, `fragment`, `path` or `segment`", v)
	}
}

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
func shouldEscape(c byte, comp Component) bool {
	// §2.3 Unreserved characters (alphanum)
	if 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || '0' <= c && c <= '9' {
		return false
	}

	switch c {
	case '-', '.', '_', '~': // §2.3 Unreserved characters (mark)
		return false

	case '!', '$', '&', '\'', '(', ')', '*', '+', ',', '/', ':', ';', '=', '?', '@': // §2.2 Reserved characters (reserved)
		// Different sections of the URL allow a few of
		// the reserved characters to appear unescaped.
		switch comp {
		case PathComponent: // §3.3
			return c == '?' || c == ':'

		case SegmentComponent: // §3.3
			// The RFC allows : @ & = + $
			// meaning to individual path segments.
			return c == '/' || c == '?' || c == ':'

		case QueryComponent: // §3.4
			return c == '&' || c == '+' || c == ':' || c == ';' || c == '='
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
func (comp Component) Escape(v string) string {
	hexCount := 0
	for i := 0; i < len(v); i++ {
		ch := v[i]
		if shouldEscape(ch, comp) {
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
		case shouldEscape(ch, comp):
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

func (comp Component) Unescape(v string) (string, error) {
	return url.PathUnescape(v)
}

// URIEncoded is a predicate that unifies the given URI component with the given encoded or decoded string.
//
// The signature is as follows:
//
//	uri_encoded(+Component, +Decoded, -Encoded)
//
// Where:
//   - Component represents the component of the URI to be escaped. It can be the atom query, fragment, path or
//     segment.
//   - Decoded represents the decoded string to be escaped.
//   - Encoded represents the encoded string.
//
// For more information on URI encoding, refer to [RFC 3986].
//
// Examples:
//
//	# Escape the given string to be used in the path component.
//	- uri_encoded(path, "foo/bar", Encoded).
//
//	# Unescape the given string to be used in the path component.
//	- uri_encoded(path, Decoded, foo%2Fbar).
//
// [RFC 3986]: https://datatracker.ietf.org/doc/html/rfc3986#section-2.1
func URIEncoded(vm *engine.VM, component, decoded, encoded engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	return engine.Delay(func(ctx context.Context) *engine.Promise {
		var comp Component
		switch c := env.Resolve(component).(type) {
		case engine.Atom:
			cc, err := NewComponent(c.String())
			if err != nil {
				return engine.Error(fmt.Errorf("uri_encoded/3: %w", err))
			}
			comp = cc
		default:
			return engine.Error(fmt.Errorf("uri_encoded/3: invalid component type: %T, should be Atom", component))
		}

		var dec string
		switch d := env.Resolve(decoded).(type) {
		case engine.Variable:
		case engine.Atom:
			dec = comp.Escape(d.String())
		default:
			return engine.Error(fmt.Errorf("uri_encoded/3: invalid decoded type: %T, should be Variable or Atom", d))
		}

		switch e := env.Resolve(encoded).(type) {
		case engine.Variable:
			return engine.Unify(vm, encoded, prolog.StringToTerm(dec), cont, env)
		case engine.Atom:
			enc, err := comp.Unescape(e.String())
			if err != nil {
				return engine.Error(fmt.Errorf("uri_encoded/3: %w", err))
			}
			return engine.Unify(vm, decoded, prolog.StringToTerm(enc), cont, env)
		default:
			return engine.Error(fmt.Errorf("uri_encoded/3: invalid encoded type: %T, should be Variable or Atom", e))
		}
	})
}
