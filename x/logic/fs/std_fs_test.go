package fs

import (
	"crypto/sha256"
	"encoding/hex"
	"io/fs"
	"slices"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	coreheader "cosmossdk.io/core/header"
	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	logictypes "github.com/axone-protocol/axoned/v14/x/logic/types"
)

type stdFSTestProgramKeeper struct {
	programs     map[string]logictypes.StoredProgram
	publications map[string]logictypes.ProgramPublication
}

func (k stdFSTestProgramKeeper) GetStoredProgram(
	_ sdk.Context, programID []byte,
) (logictypes.StoredProgram, bool, error) {
	program, found := k.programs[string(programID)]
	return program, found, nil
}

func (k stdFSTestProgramKeeper) GetProgramPublication(
	_ sdk.Context, publisher, programID []byte,
) (logictypes.ProgramPublication, bool, error) {
	publication, found := k.publications[string(publisher)+":"+string(programID)]
	return publication, found, nil
}

func TestStandardMountsLayout(t *testing.T) {
	Convey("StandardMounts should expose the canonical layout", t, func() {
		key := storetypes.NewKVStoreKey("test")
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

		paths := make([]string, 0, len(StandardMounts(testCtx.Ctx, nil, nil)))
		for _, mount := range StandardMounts(testCtx.Ctx, nil, nil) {
			paths = append(paths, mount.Path)
		}

		So(paths, ShouldHaveLength, 7)
		So(slices.Contains(paths, libPath), ShouldBeTrue)
		So(slices.Contains(paths, runHeaderPath), ShouldBeTrue)
		So(slices.Contains(paths, runCometPath), ShouldBeTrue)
		So(slices.Contains(paths, varLibBankPath), ShouldBeTrue)
		So(slices.Contains(paths, varLibLogicPath), ShouldBeTrue)
		So(slices.Contains(paths, devCodecPath), ShouldBeTrue)
		So(slices.Contains(paths, devWasmPath), ShouldBeTrue)
	})
}

func TestNewVFSCanonicalPaths(t *testing.T) {
	Convey("NewVFS should expose canonical paths", t, func() {
		key := storetypes.NewKVStoreKey("test")
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
		testCtx.Ctx = testCtx.Ctx.WithHeaderInfo(coreheader.Info{
			Height:  42,
			Time:    time.Date(2026, time.March, 11, 9, 0, 0, 0, time.UTC),
			ChainID: "axone-testchain-1",
		})

		source := "published_fact(alice)."
		programID := sha256.Sum256([]byte(source))
		publisherText := authtypes.NewModuleAddress("publisher-a").String()
		publisher, err := sdk.AccAddressFromBech32(publisherText)
		So(err, ShouldBeNil)

		programKeeper := stdFSTestProgramKeeper{
			programs: map[string]logictypes.StoredProgram{
				string(programID[:]): {
					Source:     source,
					CreatedAt:  10,
					SourceSize: uint64(len(source)),
				},
			},
			publications: map[string]logictypes.ProgramPublication{
				string(publisher) + ":" + string(programID[:]): {
					PublishedAt: 11,
				},
			},
		}

		vfs, err := NewVFS(testCtx.Ctx, nil, programKeeper)
		So(err, ShouldBeNil)

		header, err := fs.ReadFile(vfs, runHeaderPath+"/height")
		So(err, ShouldBeNil)

		programPath := publisherText + "/programs/" + fmtProgramID(programID) + ".pl"
		program, err := fs.ReadFile(vfs, varLibLogicPath+"/"+programPath)
		So(err, ShouldBeNil)

		So(string(header), ShouldEqual, "42.\n")
		So(string(program), ShouldEqual, source)
	})
}

func fmtProgramID(programID [sha256.Size]byte) string {
	return hex.EncodeToString(programID[:])
}
