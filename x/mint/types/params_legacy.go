package types

import paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

// ParamKeyTable for minting module.
//
// Deprecated: kept for migration purpose,
// will be removed soon.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamSetPairs Implements params.ParamSet.
//
// Deprecated: kept for migration purpose,
// will be removed soon.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMintDenom, &p.MintDenom, validateMintDenom),
		paramtypes.NewParamSetPair(KeyAnnualReductionFactor, &p.AnnualReductionFactor, validateAnnualReductionFactor),
		paramtypes.NewParamSetPair(KeyBlocksPerYear, &p.BlocksPerYear, validateBlocksPerYear),
	}
}
