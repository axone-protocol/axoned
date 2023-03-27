package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

// ParamKeyTable the param key table for launch module.
//
// DEPRECATED: kept for migration purpose,
// will be removed soon.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamSetPairs get the params.ParamSet.
//
// DEPRECATED: kept for migration purpose,
// will be removed soon.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return []paramtypes.ParamSetPair{
		paramtypes.NewParamSetPair(KeyInterpreter, &p.Interpreter, validateInterpreter),
		paramtypes.NewParamSetPair(KeyLimits, &p.Limits, validateLimits),
	}
}
