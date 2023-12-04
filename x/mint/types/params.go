package types

import (
	"errors"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewParams creates a new Params object.
func NewParams(
	mintDenom string, inflationCoef sdk.Dec, blocksPerYear uint64,
) Params {
	return Params{
		MintDenom:     mintDenom,
		InflationCoef: inflationCoef,
		BlocksPerYear: blocksPerYear,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		MintDenom:     sdk.DefaultBondDenom,
		InflationCoef: sdk.NewDecWithPrec(3, 2),
		BlocksPerYear: uint64(60 * 60 * 8766 / 5), // assuming 5-second block times
	}
}

// Validate is used for validating the params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationCoef(p.InflationCoef); err != nil {
		return err
	}

	return validateBlocksPerYear(p.BlocksPerYear)
}

// String implements the Stringer interface.
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func validateMintDenom(i interface{}) error {
	v, ok := i.(string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if strings.TrimSpace(v) == "" {
		return errors.New("mint denom cannot be blank")
	}

	return sdk.ValidateDenom(v)
}

func validateInflationCoef(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation coefficient cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		// while there's no theoretical limit to the inflation rate, a coefficient of
		// 1 or more would lead to hyper-hyperinflation.
		return fmt.Errorf("inflation coefficient too large: %s", v)
	}

	return nil
}

func validateBlocksPerYear(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("blocks per year must be positive: %d", v)
	}

	return nil
}
