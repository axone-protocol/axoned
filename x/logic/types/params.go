package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// Parameter store keys.
var (
	KeyInterpreter = []byte("Interpreter")
	KeyLimits      = []byte("Limits")
)

var (
	DefaultRegisteredPredicates = make([]string, 0)
	DefaultBootstrap            = ""
	DefaultMaxGas               = uint64(200000)
	DefaultMaxSize              = uint32(0)
	DefaultMaxResultCount       = uint32(50)
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
	return NewParams(DefaultInterpreter(), DefaultLimits())
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

// NewInterpreter creates a new Interpreter object.
func NewInterpreter(registeredPredicates []string, bootstrap string) Interpreter {
	return Interpreter{
		RegisteredPredicates: registeredPredicates,
		Bootstrap:            bootstrap,
	}
}

// DefaultInterpreter return an Interpreter object with default params.
func DefaultInterpreter() Interpreter {
	return NewInterpreter(DefaultRegisteredPredicates, DefaultBootstrap)
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
func NewLimits(maxGas uint64, maxSize, maxResultCount uint32) Limits {
	return Limits{
		MaxGas:         maxGas,
		MaxSize:        maxSize,
		MaxResultCount: maxResultCount,
	}
}

// DefaultLimits return a Limits object with default params.
func DefaultLimits() Limits {
	return NewLimits(DefaultMaxGas, DefaultMaxSize, DefaultMaxResultCount)
}

func validateLimits(i interface{}) error {
	_, ok := i.(Limits)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO: Validate limits params.
	return nil
}
