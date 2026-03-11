package types

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// DefaultGenesis returns the default genesis state.
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	storedPrograms := make(map[string]struct{}, len(gs.StoredPrograms))
	for _, record := range gs.StoredPrograms {
		programID, err := decodeProgramIDHex(record.ProgramId)
		if err != nil {
			return fmt.Errorf("invalid stored_programs program_id %q: %w", record.ProgramId, err)
		}
		if _, found := storedPrograms[record.ProgramId]; found {
			return fmt.Errorf("duplicate stored_programs entry for program_id %q", record.ProgramId)
		}

		expectedID := sha256.Sum256([]byte(record.Program.GetSource()))
		if string(programID) != string(expectedID[:]) {
			return fmt.Errorf("stored_programs program_id %q does not match source content hash", record.ProgramId)
		}
		if record.Program.GetSourceSize() != uint64(len(record.Program.GetSource())) {
			return fmt.Errorf("stored_programs source_size mismatch for program_id %q", record.ProgramId)
		}

		storedPrograms[record.ProgramId] = struct{}{}
	}

	publications := make(map[string]struct{}, len(gs.ProgramPublications))
	for _, record := range gs.ProgramPublications {
		if _, err := sdk.AccAddressFromBech32(record.Publisher); err != nil {
			return fmt.Errorf("invalid program_publications publisher %q: %w", record.Publisher, err)
		}
		if _, err := decodeProgramIDHex(record.ProgramId); err != nil {
			return fmt.Errorf("invalid program_publications program_id %q: %w", record.ProgramId, err)
		}
		if _, found := storedPrograms[record.ProgramId]; !found {
			return fmt.Errorf("program_publications references unknown program_id %q", record.ProgramId)
		}

		publicationKey := record.Publisher + "/" + record.ProgramId
		if _, found := publications[publicationKey]; found {
			return fmt.Errorf("duplicate program_publications entry for publisher %q and program_id %q", record.Publisher, record.ProgramId)
		}
		publications[publicationKey] = struct{}{}
	}

	return nil
}

func decodeProgramIDHex(programID string) ([]byte, error) {
	if len(programID) != sha256.Size*2 {
		return nil, fmt.Errorf("want %d hex chars", sha256.Size*2)
	}

	raw, err := hex.DecodeString(programID)
	if err != nil {
		return nil, err
	}
	if hex.EncodeToString(raw) != programID {
		return nil, fmt.Errorf("must be lowercase hexadecimal")
	}

	return raw, nil
}
