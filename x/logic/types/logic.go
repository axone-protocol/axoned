package types

import (
	"sort"

	"github.com/ichiban/prolog"
)

// TermResults is a map from variable strings to prolog term values.
type TermResults map[string]prolog.TermString

// ToSubstitutions converts a TermResults value to a slice of Substitution values.
func (t TermResults) ToSubstitutions() []Substitution {
	substitutions := make([]Substitution, 0, len(t))
	for v, ts := range t {
		substitution := Substitution{
			Variable: v,
			Term: Term{
				Name: string(ts),
			},
		}
		substitutions = append(substitutions, substitution)
	}

	return substitutions
}

// ToVariables extract from a TermResults value the variable names.
// The variable names are sorted in ascending order.
func (t TermResults) ToVariables() []string {
	variables := make([]string, 0, len(t))
	for v := range t {
		variables = append(variables, v)
	}
	sort.Strings(variables)

	return variables
}
