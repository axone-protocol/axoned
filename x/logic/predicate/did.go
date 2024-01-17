package predicate

import (
	"net/url"
	"strings"

	"github.com/ichiban/prolog/engine"
	godid "github.com/nuts-foundation/go-did/did"

	"github.com/okp4/okp4d/x/logic/prolog"
)

// AtomDID is a term which represents a DID as a compound term `did(Method, ID, Path, Query, Fragment)`.
var AtomDID = engine.NewAtom("did")

// DIDPrefix is the prefix for a DID.
const DIDPrefix = "did:"

// DIDComponents is a predicate which breaks down a DID into its components according to the [W3C DID] specification.
//
// The signature is as follows:
//
//	did_components(+DID, -Components) is det
//	did_components(-DID, +Components) is det
//
// where:
//   - DID represents DID URI, given as an Atom, compliant with [W3C DID] specification.
//   - Components is a compound Term in the format did(Method, ID, Path, Query, Fragment), aligned with the [DID syntax],
//     where: Method is The method name, ID is The method-specific identifier, Path is the path component, Query is the
//     query component and Fragment is The fragment component.
//     For any component not present, its value will be null and thus will be left as an uninstantiated variable.
//
// Examples:
//
//	# Decompose a DID into its components.
//	- did_components('did:example:123456?versionId=1', did(Method, ID, Path, Query, Fragment)).
//
//	# Reconstruct a DID from its components.
//	- did_components(DID, did('example', '123456', _, 'versionId=1', _42)).
//
// [W3C DID]: https://w3c.github.io/did-core
// [DID syntax]: https://w3c.github.io/did-core/#did-syntax
//
//nolint:funlen
func DIDComponents(vm *engine.VM, did, components engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	switch t1 := env.Resolve(did).(type) {
	case engine.Variable:
	case engine.Atom:
		parsedDid, err := godid.ParseDIDURL(t1.String())
		if err != nil {
			return engine.Error(prolog.WithError(engine.DomainError(prolog.ValidEncoding("did"), did, env), err, env))
		}

		terms, err := didToTerms(parsedDid, env)
		if err != nil {
			return engine.Error(err)
		}

		return engine.Unify(vm, components, AtomDID.Apply(terms...), cont, env)
	default:
		return engine.Error(engine.TypeError(prolog.AtomTypeAtom, did, env))
	}

	switch t2 := env.Resolve(components).(type) {
	case engine.Variable:
		return engine.Error(engine.InstantiationError(env))
	case engine.Compound:
		if t2.Functor() != AtomDID || t2.Arity() != 5 {
			return engine.Error(engine.DomainError(AtomDID, components, env))
		}

		buf := strings.Builder{}
		buf.WriteString(DIDPrefix)

		processors := []func(engine.Atom){
			func(segment engine.Atom) {
				buf.WriteString(segment.String())
			},
			func(segment engine.Atom) {
				buf.WriteString(":")
				buf.WriteString(url.PathEscape(segment.String()))
			},
			func(segment engine.Atom) {
				for _, s := range strings.FieldsFunc(segment.String(), func(c rune) bool { return c == '/' }) {
					buf.WriteString("/")
					buf.WriteString(url.PathEscape(s))
				}
			},
			func(segment engine.Atom) {
				buf.WriteString("?")
				buf.WriteString(url.PathEscape(segment.String()))
			},
			func(segment engine.Atom) {
				buf.WriteString("#")
				buf.WriteString(url.PathEscape(segment.String()))
			},
		}

		for i := 0; i < t2.Arity(); i++ {
			if err := processSegment(t2, uint8(i), processors[i], env); err != nil {
				return engine.Error(err)
			}
		}

		return engine.Unify(vm, did, engine.NewAtom(buf.String()), cont, env)
	default:
		return engine.Error(engine.TypeError(AtomDID, components, env))
	}
}

// processSegment processes a segment of a DID.
func processSegment(segments engine.Compound, segmentNumber uint8, fn func(segment engine.Atom), env *engine.Env) error {
	term := env.Resolve(segments.Arg(int(segmentNumber)))
	if _, ok := term.(engine.Variable); ok {
		return nil
	}
	segment, err := prolog.AssertAtom(env, segments.Arg(int(segmentNumber)))
	if err != nil {
		return err
	}

	fn(segment)

	return nil
}

// didToTerms converts a DID to a "tuple" of terms (either an Atom or a Variable),
// or returns an error if the conversion fails.
// The returned atoms are url decoded.
func didToTerms(did *godid.DID, env *engine.Env) ([]engine.Term, error) {
	components := []string{did.Method, did.ID, did.Path, did.Query, did.Fragment}
	terms := make([]engine.Term, 0, len(components))

	for _, component := range components {
		r, err := url.PathUnescape(component)
		if err != nil {
			return nil, engine.DomainError(prolog.ValidEncoding("url_encoded"), engine.NewAtom(component), env)
		}
		var r2 engine.Term = engine.NewAtom(r)
		terms = append(terms, r2)
	}

	return terms, nil
}
