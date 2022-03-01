package types

const (
	// ModuleName defines the module name.
	ModuleName = "knowledge"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName

	// QuerierRoute defines the module's query routing key.
	QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_knowledge"
)

var (
	// DataspaceKeyPrefix is the prefix under which the dataspace entities are stored.
	DataspaceKeyPrefix = []byte{0x11}
)

// GetDataspaceKey returns the store key referencing a specific Dataspace entity given its id.
func GetDataspaceKey(
	id string,
) []byte {
	return append(DataspaceKeyPrefix, id...)
}
