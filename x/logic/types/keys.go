package types

import (
	"encoding/binary"
	"errors"
	"math"
)

const (
	// ModuleName defines the module name.
	ModuleName = "logic"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key.
	RouterKey = ModuleName

	// MemStoreKey defines the in-memory store key.
	MemStoreKey = "mem_logic"
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	StoredProgramKeyPrefix      = []byte{0x01}
	ProgramPublicationKeyPrefix = []byte{0x02}
)

// StoredProgramKey returns the store key for a stored program identified by the raw bytes of program_id.
func StoredProgramKey(programID []byte) []byte {
	key := make([]byte, len(StoredProgramKeyPrefix)+len(programID))
	copy(key, StoredProgramKeyPrefix)
	copy(key[len(StoredProgramKeyPrefix):], programID)

	return key
}

// ProgramPublicationKey returns the store key for a publication by publisher for the given program_id.
func ProgramPublicationKey(publisher, programID []byte) []byte {
	key := make([]byte, len(ProgramPublicationKeyPrefix)+binary.MaxVarintLen64+len(publisher)+len(programID))
	offset := copy(key, ProgramPublicationKeyPrefix)
	offset += binary.PutUvarint(key[offset:], uint64(len(publisher)))
	offset += copy(key[offset:], publisher)
	offset += copy(key[offset:], programID)

	return key[:offset]
}

// ParseStoredProgramKey extracts the raw program_id bytes from a stored program key.
func ParseStoredProgramKey(key []byte) ([]byte, error) {
	if len(key) < len(StoredProgramKeyPrefix) || !hasPrefix(key, StoredProgramKeyPrefix) {
		return nil, errors.New("invalid stored program key")
	}

	return key[len(StoredProgramKeyPrefix):], nil
}

// ParseProgramPublicationKey extracts publisher and program_id bytes from a publication key.
func ParseProgramPublicationKey(key []byte) ([]byte, []byte, error) {
	if len(key) < len(ProgramPublicationKeyPrefix) || !hasPrefix(key, ProgramPublicationKeyPrefix) {
		return nil, nil, errors.New("invalid program publication key")
	}

	rest := key[len(ProgramPublicationKeyPrefix):]
	publisherLen, n := binary.Uvarint(rest)
	if n <= 0 {
		return nil, nil, errors.New("invalid program publication key length prefix")
	}

	rest = rest[n:]
	if uint64(len(rest)) < publisherLen {
		return nil, nil, errors.New("invalid program publication key publisher length")
	}
	if publisherLen > uint64(math.MaxInt) {
		return nil, nil, errors.New("invalid program publication key publisher length overflow")
	}

	publisherSize := int(publisherLen)
	publisher := rest[:publisherSize]
	programID := rest[publisherSize:]
	if len(programID) == 0 {
		return nil, nil, errors.New("invalid program publication key program_id")
	}

	return publisher, programID, nil
}

func hasPrefix(key, prefix []byte) bool {
	if len(key) < len(prefix) {
		return false
	}

	for i := range prefix {
		if key[i] != prefix[i] {
			return false
		}
	}

	return true
}
