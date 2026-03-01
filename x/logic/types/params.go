package types

import (
	"fmt"
)

// Parameter store keys.
var (
	ParamsKey = []byte("Params")
)

// NewParams creates a new Params object.
func NewParams(limits Limits, gasPolicy GasPolicy) Params {
	return Params{
		Limits:    limits,
		GasPolicy: gasPolicy,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		NewLimits(
			WithMaxSize(5000),
			WithMaxResultCount(3),
			WithMaxVariables(100000)),
		NewGasPolicy(
			WithWeightingFactor(1),
			WithDefaultPredicateCost(1)),
	)
}

// Validate validates the set of params.
func (p Params) Validate() error {
	return validateLimits(p.Limits)
}

// String implements the Stringer interface.
func (p Params) String() string {
	return p.Limits.String() + "\n" + p.GasPolicy.String()
}

// LimitsOption is a functional option for configuring the Limits.
type LimitsOption func(*Limits)

// WithMaxSize sets the max size limits accepted for a prolog program.
func WithMaxSize(maxSize uint64) LimitsOption {
	return func(i *Limits) {
		i.MaxSize = maxSize
	}
}

// WithMaxResultCount sets the maximum number of results that can be requested for a query.
func WithMaxResultCount(maxResultCount uint64) LimitsOption {
	return func(i *Limits) {
		i.MaxResultCount = maxResultCount
	}
}

// WithMaxUserOutputSize specifies the maximum number of bytes to keep in the user output.
func WithMaxUserOutputSize(maxUserOutputSize uint64) LimitsOption {
	return func(i *Limits) {
		i.MaxUserOutputSize = maxUserOutputSize
	}
}

// WithMaxVariables sets the maximum number of variables that can be created by the interpreter.
func WithMaxVariables(maxVariables uint64) LimitsOption {
	return func(i *Limits) {
		i.MaxVariables = maxVariables
	}
}

// NewLimits creates a new Limits object.
func NewLimits(opts ...LimitsOption) Limits {
	l := Limits{}
	for _, opt := range opts {
		opt(&l)
	}

	return l
}

func validateLimits(i any) error {
	_, ok := i.(Limits)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO: Validate limits params.
	return nil
}

// GasPolicyOption is a functional option for configuring the GasPolicy.
type GasPolicyOption func(*GasPolicy)

// WithWeightingFactor sets the weighting factor.
func WithWeightingFactor(weightingFactor uint64) GasPolicyOption {
	return func(i *GasPolicy) {
		i.WeightingFactor = weightingFactor
	}
}

// WithDefaultPredicateCost sets the default cost of a predicate.
func WithDefaultPredicateCost(defaultPredicateCost uint64) GasPolicyOption {
	return func(i *GasPolicy) {
		i.DefaultPredicateCost = defaultPredicateCost
	}
}

// WithPredicateCosts sets the cost of a predicate.
func WithPredicateCosts(predicateCosts []PredicateCost) GasPolicyOption {
	return func(i *GasPolicy) {
		i.PredicateCosts = predicateCosts
	}
}

// NewGasPolicy creates a new GasPolicy object.
func NewGasPolicy(opts ...GasPolicyOption) GasPolicy {
	g := GasPolicy{}
	for _, opt := range opts {
		opt(&g)
	}

	return g
}
