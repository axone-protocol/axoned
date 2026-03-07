package codec

import (
	"github.com/axone-protocol/prolog/v3/engine"
)

// Codec defines the interface for encode/decode operations.
//
// Each codec implementation handles a specific encoding scheme (e.g., bech32, base64).
// The codec is responsible for validating the request tokens and returning a
// serialized Prolog term response.
type Codec interface {
	// Name returns the codec identifier used in the device path.
	// For example, "bech32" makes the codec available at the path "bech32" relative to the mount point.
	Name() string

	// Decode processes a decode request and returns a Prolog term response.
	// The tokens slice contains the request tokens after the "decode" command.
	//
	// Returns:
	//   - ok(...) term on success
	//   - error(...) term on validation failure
	Decode(tokens [][]byte) engine.Term

	// Encode processes an encode request and returns a Prolog term response.
	// The tokens slice contains the request tokens after the "encode" command.
	//
	// Returns:
	//   - ok(...) term on success
	//   - error(...) term on validation failure
	Encode(tokens [][]byte) engine.Term
}

// registry holds all available codecs indexed by name.
var registry = make(map[string]Codec)

// Register adds a codec to the global registry.
// This should be called during initialization (e.g., in init() functions).
func Register(codec Codec) {
	registry[codec.Name()] = codec
}

// Get retrieves a codec by name from the registry.
// Returns nil if the codec is not registered.
func Get(name string) Codec {
	return registry[name]
}

// All returns all registered codec names.
func All() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
