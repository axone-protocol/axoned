package wasm

import (
	"github.com/axone-protocol/axoned/v11/x/logic/types"
)

// AskQuery implements the wasm custom Ask query JSON schema, it basically redefined the Ask gRPC request parameters
// to keep control in case of eventual breaking change in the logic module definition, and to decouple the
// serialization logic.
type AskQuery struct {
	Program string `json:"program"`
	Query   string `json:"query"`
	Limit   uint64 `json:"limit"`
}

// AskResponse implements the Ask query response JSON schema in a wasm custom query purpose, it redefines the existing
// generated type from proto to ensure a dedicated serialization logic.
type AskResponse struct {
	Height     uint64  `json:"height"`
	GasUsed    uint64  `json:"gas_used"`
	Answer     *Answer `json:"answer,omitempty"`
	UserOutput string  `json:"user_output,omitempty"`
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
	to.UserOutput = from.UserOutput
}

// Answer denotes the Answer element JSON representation in an AskResponse for wasm custom query purpose, it redefines
// the existing generated type from proto to ensure a dedicated serialization logic.
type Answer struct {
	HasMore   bool     `json:"has_more"`
	Variables []string `json:"variables"`
	Results   []Result `json:"results"`
}

func (to *Answer) from(from types.Answer) {
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
	Error         string         `json:"error,omitempty"`
	Substitutions []Substitution `json:"substitutions"`
}

func (to *Result) from(from types.Result) {
	to.Error = from.Error
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
	Variable   string `json:"variable"`
	Expression string `json:"expression"`
}

func (to *Substitution) from(from types.Substitution) {
	to.Variable = from.Variable
	to.Expression = from.Expression
}
