package predicate

import (
	"github.com/axone-protocol/prolog/v2/engine"

	"github.com/axone-protocol/axoned/v11/x/logic/prolog"
)

// URIEncoded is a predicate that unifies the given URI component with the given encoded or decoded string.
//
// The signature is as follows:
//
//	uri_encoded(+Component, +Value, -Encoded) is det
//	uri_encoded(+Component, -Value, +Encoded) is det
//
// Where:
//   - Component represents the component of the URI to be escaped. It can be the atom 'query_path', 'fragment', 'path' or
//     'segment'.
//   - Decoded represents the decoded string to be escaped.
//   - Encoded represents the encoded string.
//
// For more information on URI encoding, refer to [RFC 3986].
//
// # Examples:
//
//	# Escape the given string to be used in the path component.
//	- uri_encoded(path, "foo/bar", Encoded).
//
//	# Unescape the given string to be used in the path component.
//	- uri_encoded(path, Decoded, foo%2Fbar).
//
// [RFC 3986]: https://datatracker.ietf.org/doc/html/rfc3986#section-2.1
func URIEncoded(_ *engine.VM, component, decoded, encoded engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	_, err := prolog.AssertIsGround(component, env)
	if err != nil {
		return engine.Error(err)
	}
	uriComponent, err := prolog.AssertURIComponent(component, env)
	if err != nil {
		return engine.Error(err)
	}
	forwardConverter := func(value []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		in, err := prolog.TextTermToString(value[0], env)
		if err != nil {
			return nil, err
		}
		out := uriComponent.Escape(in)
		return []engine.Term{engine.NewAtom(out)}, nil
	}
	backwardConverter := func(value []engine.Term, _ engine.Term, env *engine.Env) ([]engine.Term, error) {
		in, err := prolog.TextTermToString(value[0], env)
		if err != nil {
			return nil, err
		}
		out, err := uriComponent.Unescape(in)
		if err != nil {
			return nil, prolog.WithError(engine.DomainError(prolog.ValidEncoding("uri"), value[0], env), err, env)
		}
		return []engine.Term{engine.NewAtom(out)}, nil
	}
	return prolog.UnifyFunctionalPredicate(
		[]engine.Term{decoded}, []engine.Term{encoded}, prolog.AtomEmpty, forwardConverter, backwardConverter, cont, env)
}
