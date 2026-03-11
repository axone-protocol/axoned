package logic

import (
	"encoding/hex"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/keeper"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	err := k.SetParams(ctx, genState.Params)
	if err != nil {
		panic(errorsmod.Wrapf(err, "error setting params"))
	}

	for _, record := range genState.StoredPrograms {
		programID, err := hex.DecodeString(record.ProgramId)
		if err != nil {
			panic(errorsmod.Wrapf(err, "error decoding stored program id %s", record.ProgramId))
		}

		if err := k.SetStoredProgram(ctx, programID, record.Program); err != nil {
			panic(errorsmod.Wrapf(err, "error setting stored program %s", record.ProgramId))
		}
	}

	for _, record := range genState.ProgramPublications {
		publisher, err := sdk.AccAddressFromBech32(record.Publisher)
		if err != nil {
			panic(errorsmod.Wrapf(err, "error decoding publisher %s", record.Publisher))
		}
		programID, err := hex.DecodeString(record.ProgramId)
		if err != nil {
			panic(errorsmod.Wrapf(err, "error decoding program publication id %s", record.ProgramId))
		}

		if err := k.SetProgramPublication(ctx, publisher, programID, record.Publication); err != nil {
			panic(errorsmod.Wrapf(err, "error setting program publication %s/%s", record.Publisher, record.ProgramId))
		}
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	if err := k.IterateStoredPrograms(ctx, func(programID []byte, program types.StoredProgram) (bool, error) {
		genesis.StoredPrograms = append(genesis.StoredPrograms, types.GenesisStoredProgram{
			ProgramId: hex.EncodeToString(programID),
			Program:   program,
		})
		return false, nil
	}); err != nil {
		panic(errorsmod.Wrapf(err, "error exporting stored programs"))
	}

	if err := k.IterateProgramPublications(
		ctx,
		func(publisher, programID []byte, publication types.ProgramPublication) (bool, error) {
			genesis.ProgramPublications = append(genesis.ProgramPublications, types.GenesisProgramPublication{
				Publisher:   sdk.AccAddress(publisher).String(),
				ProgramId:   hex.EncodeToString(programID),
				Publication: publication,
			})
			return false, nil
		},
	); err != nil {
		panic(errorsmod.Wrapf(err, "error exporting program publications"))
	}

	return genesis
}
