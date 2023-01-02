package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	// TODO implement me
	return Params{}
}

// ParamSetPairs get the params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	// TODO implement me
	return []paramtypes.ParamSetPair{}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	// TODO implement me
	return nil
}

// String implements the Stringer interface.
func (p Params) String() string {
	// TODO implement me
	return ""
}
