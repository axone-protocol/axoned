package logic_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"testing"
	"testing/fstest"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v15/x/logic"
	"github.com/axone-protocol/axoned/v15/x/logic/keeper"
	"github.com/axone-protocol/axoned/v15/x/logic/types"
)

func TestGenesisRoundTripStoredPrograms(t *testing.T) {
	Convey("Given a logic keeper and a genesis state with stored programs", t, func() {
		encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
		key := storetypes.NewKVStoreKey(types.StoreKey)
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

		logicKeeper := keeper.NewKeeper(
			encCfg.Codec,
			encCfg.InterfaceRegistry,
			key,
			key,
			authtypes.NewModuleAddress(govtypes.ModuleName),
			nil,
			nil,
			nil,
			func(_ context.Context) (fs.FS, error) {
				return fstest.MapFS{}, nil
			},
		)

		source := "father(alice, bob)."
		programID := sha256.Sum256([]byte(source))
		programIDHex := hex.EncodeToString(programID[:])
		publisher := authtypes.NewModuleAddress("publisher-a").String()
		genesis := types.GenesisState{
			Params: types.DefaultParams(),
			StoredPrograms: []types.GenesisStoredProgram{{
				ProgramId: programIDHex,
				Program: types.StoredProgram{
					Source:     source,
					CreatedAt:  11,
					SourceSize: uint64(len(source)),
				},
			}},
			ProgramPublications: []types.GenesisProgramPublication{{
				Publisher: publisher,
				ProgramId: programIDHex,
				Publication: types.ProgramPublication{
					PublishedAt: 22,
				},
			}},
		}

		Convey("when initializing and exporting genesis", func() {
			logic.InitGenesis(testCtx.Ctx, *logicKeeper, genesis)
			exported := logic.ExportGenesis(testCtx.Ctx, *logicKeeper)

			Convey("then stored programs and publications should round-trip", func() {
				So(exported.Params, ShouldResemble, genesis.Params)
				So(exported.StoredPrograms, ShouldResemble, genesis.StoredPrograms)
				So(exported.ProgramPublications, ShouldResemble, genesis.ProgramPublications)
			})
		})
	})
}
