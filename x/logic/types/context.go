package types

// ContextKey is a type for context keys.
type ContextKey string

const (
	// AuthKeeperContextKey is the context key for the auth keeper.
	AuthKeeperContextKey = ContextKey("authKeeper")
	// BankKeeperContextKey is the context key for the bank keeper.
	BankKeeperContextKey = ContextKey("bankKeeper")
	// WasmKeeperContextKey is the context key for the wasm keeper.
	WasmKeeperContextKey = ContextKey("wasmKeeper")
)
