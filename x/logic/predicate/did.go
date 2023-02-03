package predicate

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/ichiban/prolog/engine"
	godid "github.com/nuts-foundation/go-did/did"
	"github.com/okp4/okp4d/x/logic/util"
)

// AtomDID is a term which represents a DID as a compound term `did(Method, ID, Path, Query, Fragment)`.
var AtomDID = engine.NewAtom("did")

// DIDComponents is a predicate which breaks down a DID into its components according to the [W3C DID] specification.
//
// did_components(+DID, -Components) is det
// did_components(-DID, +Components) is det
//
// where:
//   - `DID` represents the DID URI as a `text`, compliant with the [W3C DID] specification.
//   - `Components` is a term `did(Method, ID, Path, Query, Fragment)` following the [DID syntax] which represents
//     respectively the method name, the method-specific ID, the path, the query, and the fragment of the DID, in decoded
//     form. Components that are not found (i.e. `null`) are left uninstantiated (variable).
//
// Example:
//
// # Decompose a DID into its components.
// - did_components('did:example:123456?versionId=1', did(Method, ID, Path, Query, Fragment)).
//
// # Reconstruct a DID from its components.
// - did_components(DID, did('example', '123456', null, 'versionId=1', _42)).
//
// [W3C DID]: https://w3c.github.io/did-core
// [DID syntax]: https://w3c.github.io/did-core/#did-syntax
func DIDComponents(vm *engine.VM, did, components engine.Term, cont engine.Cont, env *engine.Env) *engine.Promise {
	switch t1 := env.Resolve(did).(type) {
	case engine.Variable:
	case engine.Atom:
		parsedDid, err := godid.ParseDIDURL(t1.String())
		if err != nil {
			return engine.Error(fmt.Errorf("did_components/2: %w", err))
		}

		terms, err := didToTerms(parsedDid)
		if err != nil {
			return engine.Error(fmt.Errorf("did_components/2: %w", err))
		}

		return engine.Unify(vm, components, AtomDID.Apply(terms...), cont, env)
	default:
		return engine.Error(fmt.Errorf("did_components/2: cannot unify did with %T", t1))
	}

	switch t2 := env.Resolve(components).(type) {
	case engine.Variable:
		return engine.Error(fmt.Errorf("did_components/2: at least one argument must be instantiated"))
	case engine.Compound:
		if t2.Functor() != AtomDID {
			return engine.Error(fmt.Errorf("did_components/2: invalid functor %s. Expected %s", t2.Functor().String(), AtomDID.String()))
		}
		if t2.Arity() != 5 {
			return engine.Error(fmt.Errorf("did_components/2: invalid arity %d. Expected 5", t2.Arity()))
		}

		buf := strings.Builder{}
		buf.WriteString("did:")
		if segment, ok := util.Resolve(env, t2.Arg(0)); ok {
			buf.WriteString(url.PathEscape(segment.String()))
		}
		if segment, ok := util.Resolve(env, t2.Arg(1)); ok {
			buf.WriteString(":")
			buf.WriteString(url.PathEscape(segment.String()))
		}
		if segment, ok := util.Resolve(env, t2.Arg(2)); ok {
			for _, s := range strings.FieldsFunc(segment.String(), func(c rune) bool { return c == '/' }) {
				buf.WriteString("/")
				buf.WriteString(url.PathEscape(s))
			}
		}
		if segment, ok := util.Resolve(env, t2.Arg(3)); ok {
			buf.WriteString("?")
			buf.WriteString(url.PathEscape(segment.String()))
		}
		if segment, ok := util.Resolve(env, t2.Arg(4)); ok {
			buf.WriteString("#")
			buf.WriteString(url.PathEscape(segment.String()))
		}
		return engine.Unify(vm, did, engine.NewAtom(buf.String()), cont, env)
	default:
		return engine.Error(fmt.Errorf("did_components/2: cannot unify did with %T", t2))
	}
}

// didToTerms converts a DID to a "tuple" of terms (either an Atom or a Variable),
// or returns an error if the conversion fails.
// The returned atoms are url decoded.
func didToTerms(did *godid.DID) ([]engine.Term, error) {
	components := []string{did.Method, did.ID, did.Path, did.Query, did.Fragment}
	terms := make([]engine.Term, 0, len(components))

	for _, component := range components {
		r, err := url.PathUnescape(component)
		if err != nil {
			return nil, err
		}
		terms = append(terms, util.StringToTerm(r))
	}

	return terms, nil
}
