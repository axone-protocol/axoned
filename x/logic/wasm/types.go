package wasm

import "github.com/okp4/okp4d/x/logic/types"

// AskQuery implements the wasm custom Ask query JSON schema, it basically redefined the Ask gRPC request parameters
// to keep control in case of eventual breaking change in the logic module definition, and to decouple the
// serialization logic.
type AskQuery struct {
	Program string `json:"program"`
	Query   string `json:"query"`
}

// AskResponse implements the Ask query response JSON schema in a wasm custom query purpose, it redefines the existing
// generated type from proto to ensure a dedicated serialization logic.
type AskResponse struct {
	Height  uint64  `json:"height"`
	GasUsed uint64  `json:"gas_used"`
	Answer  *Answer `json:"answer,omitempty"`
}

func (to *AskResponse) from(from types.QueryServiceAskResponse) {
	to.Height = from.Height
	to.GasUsed = from.GasUsed
	to.Answer = nil
	if from.Answer != nil {
		answer := new(Answer)
		answer.from(*from.Answer)
		to.Answer = answer
	}
}

// Answer denotes the Answer element JSON representation in an AskResponse for wasm custom query purpose, it redefines
// the existing generated type from proto to ensure a dedicated serialization logic.
type Answer struct {
	Success   bool     `json:"success"`
	HasMore   bool     `json:"has_more"`
	Variables []string `json:"variables"`
	Results   []Result `json:"results"`
}

func (to *Answer) from(from types.Answer) {
	to.Success = from.Success
	to.HasMore = from.HasMore
	to.Variables = from.Variables
	if to.Variables == nil {
		to.Variables = make([]string, 0)
	}
	to.Results = make([]Result, 0, len(from.Results))
	for _, fromResult := range from.Results {
		result := new(Result)
		result.from(fromResult)
		to.Results = append(to.Results, *result)
	}
}

// Result denotes the Result element JSON representation in an AskResponse for wasm custom query purpose, it redefines
// the existing generated type from proto to ensure a dedicated serialization logic.
type Result struct {
	Substitutions []Substitution `json:"substitutions"`
}

func (to *Result) from(from types.Result) {
	to.Substitutions = make([]Substitution, 0, len(from.Substitutions))
	for _, fromSubstitution := range from.Substitutions {
		substitution := new(Substitution)
		substitution.from(fromSubstitution)
		to.Substitutions = append(to.Substitutions, *substitution)
	}
}

// Substitution denotes the Substitution element JSON representation in an AskResponse for wasm custom query purpose, it redefines
// the existing generated type from proto to ensure a dedicated serialization logic.
type Substitution struct {
	Variable string `json:"variable"`
	Term     Term   `json:"term"`
}

func (to *Substitution) from(from types.Substitution) {
	to.Variable = from.Variable
	term := new(Term)
	term.from(from.Term)
	to.Term = *term
}

// Term denotes the Term element JSON representation in an AskResponse for wasm custom query purpose, it redefines
// the existing generated type from proto to ensure a dedicated serialization logic.
type Term struct {
	Name      string `json:"name"`
	Arguments []Term `json:"arguments"`
}

func (to *Term) from(from types.Term) {
	to.Name = from.Name
	to.Arguments = make([]Term, 0, len(from.Arguments))
	for _, fromTerm := range from.Arguments {
		term := new(Term)
		term.from(fromTerm)
		to.Arguments = append(to.Arguments, *term)
	}
}
