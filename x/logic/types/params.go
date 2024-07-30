package types

import (
	"fmt"
	"github.com/axone-protocol/axoned/v8/x/logic/interpreter"
	"net/url"

	"cosmossdk.io/math"
)

// Parameter store keys.
var (
	ParamsKey = []byte("Params")
)

var (
	DefaultPredicatesWhitelist = make([]string, 0)
	DefaultPredicatesBlacklist = make([]string, 0)
	DefaultMaxSize             = math.NewUint(uint64(5000))
	DefaultMaxResultCount      = math.NewUint(uint64(1))
	DefaultMaxVariables        = math.NewUint(uint64(100000))
)

// NewParams creates a new Params object.
func NewParams(interpreter Interpreter, limits Limits) Params {
	return Params{
		Interpreter: interpreter,
		Limits:      limits,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(NewInterpreter(), NewLimits())
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
		i.PredicatesFilter.Whitelist = DefaultPredicatesWhitelist
	}

	if i.PredicatesFilter.Blacklist == nil {
		i.PredicatesFilter.Blacklist = DefaultPredicatesBlacklist
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

func validateInterpreter(i interface{}) error {
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

	for _, predicate := range interpreter.PredicatesFilter.Whitelist {
		if err := validatePredicate(predicate); err != nil {
			return fmt.Errorf("invalide predicates filter whitelist : %v", err)
		}
	}

	for _, predicate := range interpreter.PredicatesFilter.Blacklist {
		if err := validatePredicate(predicate); err != nil {
			return fmt.Errorf("invalide predicates filter blacklist : %v", err)
		}
	}
	return nil
}

func validatePredicate(predicate string) error {
	for _, p := range interpreter.RegistryNames {
		if predicate == p {
			return nil
		}
	}
	return fmt.Errorf("unknown predicate: %s", predicate)
}

// LimitsOption is a functional option for configuring the Limits.
type LimitsOption func(*Limits)

// WithMaxSize sets the max size limits accepted for a prolog program.
func WithMaxSize(maxSize math.Uint) LimitsOption {
	return func(i *Limits) {
		i.MaxSize = &maxSize
	}
}

// WithMaxResultCount sets the maximum number of results that can be requested for a query.
func WithMaxResultCount(maxResultCount math.Uint) LimitsOption {
	return func(i *Limits) {
		i.MaxResultCount = &maxResultCount
	}
}

// WithMaxUserOutputSize specifies the maximum number of bytes to keep in the user output.
func WithMaxUserOutputSize(maxUserOutputSize math.Uint) LimitsOption {
	return func(i *Limits) {
		i.MaxUserOutputSize = &maxUserOutputSize
	}
}

// WithMaxVariables sets the maximum number of variables that can be created by the interpreter.
func WithMaxVariables(maxVariables math.Uint) LimitsOption {
	return func(i *Limits) {
		i.MaxVariables = &maxVariables
	}
}

// NewLimits creates a new Limits object.
func NewLimits(opts ...LimitsOption) Limits {
	l := Limits{}
	for _, opt := range opts {
		opt(&l)
	}

	if l.MaxSize == nil {
		l.MaxSize = &DefaultMaxSize
	}

	if l.MaxResultCount == nil {
		l.MaxResultCount = &DefaultMaxResultCount
	}

	if l.MaxVariables == nil {
		l.MaxVariables = &DefaultMaxVariables
	}

	return l
}

func validateLimits(i interface{}) error {
	_, ok := i.(Limits)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	// TODO: Validate limits params.
	return nil
}
