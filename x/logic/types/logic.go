package types

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
)

// TermResults is a map from variable strings to prolog term values.
type TermResults map[string]interface{}

// ToSubstitutions converts a TermResults value to a slice of Substitution values.
func (t TermResults) ToSubstitutions() []Substitution {
	substitutions := make([]Substitution, 0, len(t))
	for v, ts := range t {
		var term string
		if reflect.TypeOf(ts).Kind() == reflect.Slice {
			var buf bytes.Buffer
			for _, t := range ts.([]interface{}) {
				buf.WriteString(fmt.Sprintf("%v", t))
			}
			term = buf.String()
		} else {
			term = fmt.Sprintf("%v", ts)
		}

		substitution := Substitution{
			Variable: v,
			Term: Term{
				Name: term,
			},
		}
		substitutions = append(substitutions, substitution)
	}

	return substitutions
}

// ToVariables extract from a TermResults value the variable names.
func (t TermResults) ToVariables() []string {
	variables := make([]string, 0, len(t))
	for v := range t {
		variables = append(variables, v)
	}
	sort.Strings(variables)

	return variables
}
