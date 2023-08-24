package types

import (
	"errors"
	"fmt"
	"strings"

	"sigs.k8s.io/yaml"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Parameter store keys.
var (
	KeyMintDenom             = []byte("MintDenom")
	KeyAnnualReductionFactor = []byte("AnnualReductionFactor")
	KeyBlocksPerYear         = []byte("BlocksPerYear")
)

func NewParams(
	mintDenom string, annualReductionFactor sdk.Dec, blocksPerYear uint64,
) Params {
	return Params{
		MintDenom:             mintDenom,
		AnnualReductionFactor: annualReductionFactor,
		BlocksPerYear:         blocksPerYear,
	}
}

// default minting module parameters.
func DefaultParams() Params {
	return Params{
		MintDenom:             sdk.DefaultBondDenom,
		AnnualReductionFactor: sdk.NewDecWithPrec(20, 2),  // Tha annual reduction factor is configured to 20% per year
		BlocksPerYear:         uint64(60 * 60 * 8766 / 5), // assuming 5-second block times
	}
}

// validate params.
func (p Params) Validate() error {
	if err := validateMintDenom(p.MintDenom); err != nil {
		return err
	}
	if err := validateAnnualReductionFactor(p.AnnualReductionFactor); err != nil {
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

func validateAnnualReductionFactor(i interface{}) error {
	v, ok := i.(sdk.Dec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("annual reduction factor cannot be negative: %s", v)
	}
	if v.GT(sdk.OneDec()) {
		return fmt.Errorf("annual reduction factor too large: %s", v)
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
