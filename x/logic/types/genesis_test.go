package types_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

func TestGenesisState_Validate(t *testing.T) {
	Convey("Given genesis state validation cases", t, func() {
		source := "father(alice, bob)."
		programID := sha256.Sum256([]byte(source))
		programIDHex := hex.EncodeToString(programID[:])
		publisher := authtypes.NewModuleAddress("publisher-a").String()

		for _, tc := range []struct {
			desc     string
			genState *types.GenesisState
			valid    bool
		}{
			{
				desc:     "default is valid",
				genState: types.DefaultGenesis(),
				valid:    true,
			},
			{
				desc:     "valid genesis state",
				genState: &types.GenesisState{},
				valid:    true,
			},
			{
				desc: "valid stored program publication graph",
				genState: &types.GenesisState{
					StoredPrograms: []types.GenesisStoredProgram{{
						ProgramId: programIDHex,
						Program: types.StoredProgram{
							Source:     source,
							CreatedAt:  1,
							SourceSize: uint64(len(source)),
						},
					}},
					ProgramPublications: []types.GenesisProgramPublication{{
						Publisher: publisher,
						ProgramId: programIDHex,
						Publication: types.ProgramPublication{
							PublishedAt: 2,
						},
					}},
				},
				valid: true,
			},
			{
				desc: "stored program id must match source hash",
				genState: &types.GenesisState{
					StoredPrograms: []types.GenesisStoredProgram{{
						ProgramId: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
						Program: types.StoredProgram{
							Source:     source,
							CreatedAt:  1,
							SourceSize: uint64(len(source)),
						},
					}},
				},
				valid: false,
			},
			{
				desc: "publication must reference an existing stored program",
				genState: &types.GenesisState{
					ProgramPublications: []types.GenesisProgramPublication{{
						Publisher: publisher,
						ProgramId: programIDHex,
						Publication: types.ProgramPublication{
							PublishedAt: 2,
						},
					}},
				},
				valid: false,
			},
		} {
			Convey(tc.desc, func() {
				err := tc.genState.Validate()
				if tc.valid {
					So(err, ShouldBeNil)
				} else {
					So(err, ShouldNotBeNil)
				}
			})
		}
	})
}
