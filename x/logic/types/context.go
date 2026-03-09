package types

// ContextKey is a type for context keys.
type ContextKey string

const (
	// InterfaceRegistryContextKey is the context key for the interface registry.
	InterfaceRegistryContextKey = ContextKey("interfaceRegistry")
	// AuthKeeperContextKey is the context key for the auth keeper.
	AuthKeeperContextKey = ContextKey("authKeeper")
	// AuthQueryServiceContextKey is the context key for the auth query service.
	AuthQueryServiceContextKey = ContextKey("authQueryService")
	// BankKeeperContextKey is the context key for the bank keeper.
	BankKeeperContextKey = ContextKey("bankKeeper")
	// IOCoeffContextKey is the context key for the I/O gas coefficient used by devices.
	IOCoeffContextKey = ContextKey("ioCoeff")
)
