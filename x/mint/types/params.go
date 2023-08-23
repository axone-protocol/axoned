package types

import (
	"errors"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewParams(
	mintDenom string, inflationCoef, boundingAdjustment, targetBoundingRatio sdk.Dec, blocksPerYear uint64,
) Params {
	return Params{
		MintDenom:           mintDenom,
		InflationCoef:       inflationCoef,
		BoundingAdjustment:  boundingAdjustment,
		TargetBoundingRatio: targetBoundingRatio,
		BlocksPerYear:       blocksPerYear,
	}
}

// default minting module parameters.
func DefaultParams() Params {
	return Params{
		MintDenom:           sdk.DefaultBondDenom,
		InflationCoef:       sdk.NewDecWithPrec(73, 3),
		BoundingAdjustment:  sdk.NewDecWithPrec(25, 1),
		TargetBoundingRatio: sdk.NewDecWithPrec(66, 2),
		BlocksPerYear:       uint64(60 * 60 * 8766 / 5), // assuming 5-second block times
	}
}

// validate params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateInflationCoef(p.InflationCoef); err != nil {
		return err
	}
	if err := validateBoundingAdjustment(p.BoundingAdjustment); err != nil {
		return err
	}
	if err := validateTargetBoundingRatio(p.TargetBoundingRatio); err != nil {
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
		return fmt.Errorf("inflation coefficient too large: %s", v)
	}

	return nil
}

func validateBoundingAdjustment(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("inflation coefficient cannot be negative: %s", v)
	}

	return nil
}

func validateTargetBoundingRatio(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("target bounding ratio cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("target bounding ratio too large: %s", v)
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
