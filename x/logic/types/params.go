package types

import (
	"fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter store keys.
var (
	KeyInterpreter = []byte("Interpreter")
	KeyLimits      = []byte("Limits")
)

var (
	DefaultPredicatesWhitelist = make([]string, 0)
	DefaultPredicatesBlacklist = make([]string, 0)
	DefaultBootstrap           = ""
	DefaultMaxGas              = math.NewUint(uint64(100000))
	DefaultMaxSize             = math.NewUint(uint64(5000))
	DefaultMaxResultCount      = math.NewUint(uint64(1))
)

// ParamKeyTable the param key table for launch module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params object.
func NewParams(interpreter Interpreter, limits Limits) Params {
	return Params{
		Interpreter: interpreter,
		Limits:      limits,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(NewInterpreter(), DefaultLimits())
}

// ParamSetPairs get the params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return []paramtypes.ParamSetPair{
		paramtypes.NewParamSetPair(KeyInterpreter, &p.Interpreter, validateInterpreter),
		paramtypes.NewParamSetPair(KeyLimits, &p.Limits, validateLimits),
	}
}

// Validate validates the set of params.
func (p Params) Validate() error {
	if err := validateInterpreter(p.Interpreter); err != nil {
		return err
	}
	if err := validateLimits(p.Limits); err != nil {
		return err
	}

	return nil
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

	if i.PredicatesWhitelist == nil {
		i.PredicatesWhitelist = DefaultPredicatesWhitelist
	}

	if i.PredicatesBlacklist == nil {
		i.PredicatesBlacklist = DefaultPredicatesBlacklist
	}

	return i
}

// InterpreterOption is a functional option for configuring the Interpreter.
type InterpreterOption func(*Interpreter)

// WithPredicatesWhitelist sets the whitelist of predicates.
func WithPredicatesWhitelist(whitelist []string) InterpreterOption {
	return func(i *Interpreter) {
		i.PredicatesWhitelist = whitelist
	}
}

// WithPredicatesBlacklist sets the blacklist of predicates.
func WithPredicatesBlacklist(blacklist []string) InterpreterOption {
	return func(i *Interpreter) {
		i.PredicatesBlacklist = blacklist
	}
}

// WithBootstrap sets the bootstrap program.
func WithBootstrap(bootstrap string) InterpreterOption {
	return func(i *Interpreter) {
		i.Bootstrap = bootstrap
	}
}

func validateInterpreter(i interface{}) error {
	_, ok := i.(Interpreter)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO: Validate interpreter params.
	return nil
}

// NewLimits creates a new Limits object.
func NewLimits(maxGas, maxSize, maxResultCount *math.Uint) Limits {
	return Limits{
		MaxGas:         maxGas,
		MaxSize:        maxSize,
		MaxResultCount: maxResultCount,
	}
}

// DefaultLimits return a Limits object with default params.
func DefaultLimits() Limits {
	return NewLimits(&DefaultMaxGas, &DefaultMaxSize, &DefaultMaxResultCount)
}

func validateLimits(i interface{}) error {
	_, ok := i.(Limits)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO: Validate limits params.
	return nil
}
