package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

// Parameter store keys.
var (
	KeyInterpreter = []byte("Interpreter")
	KeyLimits      = []byte("Limits")
)

// String implements the Stringer interface.
func (p *Params) String() string {
	return p.Interpreter.String() + "\n" +
		p.Limits.String()
}

// ParamSetPairs get the params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	noopValidator := func(i interface{}) error { return nil }

	return []paramtypes.ParamSetPair{
		paramtypes.NewParamSetPair(KeyInterpreter, &p.Interpreter, noopValidator),
		paramtypes.NewParamSetPair(KeyLimits, &p.Limits, noopValidator),
	}
}
