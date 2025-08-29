package types

import (
	"fmt"
	"net/url"
)

// Parameter store keys.
var (
	ParamsKey = []byte("Params")
)

// NewParams creates a new Params object.
func NewParams(interpreter Interpreter, limits Limits, gasPolicy GasPolicy) Params {
	return Params{
		Interpreter: interpreter,
		Limits:      limits,
		GasPolicy:   gasPolicy,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		NewInterpreter(),
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
	if err := validateInterpreter(p.Interpreter); err != nil {
		return err
	}
	return validateLimits(p.Limits)
}

// String implements the Stringer interface.
func (p Params) String() string {
	return p.Interpreter.String() + "\n" +
		p.Limits.String()
}

// NewInterpreter creates a new Interpreter with the given options.
func NewInterpreter(opts ...InterpreterOption) Interpreter {
	i := Interpreter{}
	for _, opt := range opts {
		opt(&i)
	}

	if i.PredicatesFilter.Whitelist == nil {
		i.PredicatesFilter.Whitelist = []string{}
	}

	if i.PredicatesFilter.Blacklist == nil {
		i.PredicatesFilter.Blacklist = []string{}
	}

	return i
}

// InterpreterOption is a functional option for configuring the Interpreter.
type InterpreterOption func(*Interpreter)

// WithPredicatesWhitelist sets the whitelist of predicates.
func WithPredicatesWhitelist(whitelist []string) InterpreterOption {
	return func(i *Interpreter) {
		i.PredicatesFilter.Whitelist = whitelist
	}
}

// WithPredicatesBlacklist sets the blacklist of predicates.
func WithPredicatesBlacklist(blacklist []string) InterpreterOption {
	return func(i *Interpreter) {
		i.PredicatesFilter.Blacklist = blacklist
	}
}

// WithVirtualFilesWhitelist sets the whitelist of predicates.
func WithVirtualFilesWhitelist(whitelist []string) InterpreterOption {
	return func(i *Interpreter) {
		i.VirtualFilesFilter.Whitelist = whitelist
	}
}

// WithVirtualFilesBlacklist sets the blacklist of predicates.
func WithVirtualFilesBlacklist(blacklist []string) InterpreterOption {
	return func(i *Interpreter) {
		i.VirtualFilesFilter.Blacklist = blacklist
	}
}

// WithBootstrap sets the bootstrap program.
func WithBootstrap(bootstrap string) InterpreterOption {
	return func(i *Interpreter) {
		i.Bootstrap = bootstrap
	}
}

func validateInterpreter(i any) error {
	interpreter, ok := i.(Interpreter)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, file := range interpreter.VirtualFilesFilter.Whitelist {
		if _, err := url.Parse(file); err != nil {
			return fmt.Errorf("invalid virtual file in whitelist: %s", file)
		}
	}
	for _, file := range interpreter.VirtualFilesFilter.Blacklist {
		if _, err := url.Parse(file); err != nil {
			return fmt.Errorf("invalid virtual file in blacklist: %s", file)
		}
	}

	return nil
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
