package share

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io/fs"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"

	logictypes "github.com/axone-protocol/axoned/v14/x/logic/types"
)

type testKeeper struct {
	programs       map[string]logictypes.StoredProgram
	publications   map[string]logictypes.ProgramPublication
	programErr     error
	publicationErr error
}

func (k *testKeeper) GetStoredProgram(_ sdk.Context, programID []byte) (logictypes.StoredProgram, bool, error) {
	if k.programErr != nil {
		return logictypes.StoredProgram{}, false, k.programErr
	}

	p, found := k.programs[string(programID)]
	return p, found, nil
}

func (k *testKeeper) GetProgramPublication(
	_ sdk.Context, publisher, programID []byte,
) (logictypes.ProgramPublication, bool, error) {
	if k.publicationErr != nil {
		return logictypes.ProgramPublication{}, false, k.publicationErr
	}

	p, found := k.publications[string(publisher)+":"+string(programID)]
	return p, found, nil
}

func TestUserFSReadFile(t *testing.T) {
	Convey("Given a user library filesystem", t, func() {
		key := storetypes.NewKVStoreKey("test")
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

		source := "father(alice, bob)."
		id := sha256.Sum256([]byte(source))
		now := time.Date(2026, time.March, 9, 15, 0, 0, 0, time.UTC)
		publisher := sdk.AccAddress([]byte("publisher-address-01"))

		fsys := NewFS(testCtx.Ctx, &testKeeper{
			programs: map[string]logictypes.StoredProgram{
				string(id[:]): {
					Source:     source,
					CreatedAt:  now.Unix(),
					SourceSize: uint64(len(source)),
				},
			},
			publications: map[string]logictypes.ProgramPublication{
				string(publisher) + ":" + string(id[:]): {
					PublishedAt: now.Unix(),
				},
			},
		})

		Convey("when reading an existing program file", func() {
			content, err := fsys.ReadFile(publishedPath(publisher.String(), id))

			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, source)
		})

		Convey("when reading with a leading slash", func() {
			content, err := fsys.ReadFile("/" + publishedPath(publisher.String(), id))

			So(err, ShouldBeNil)
			So(string(content), ShouldEqual, source)
		})

		Convey("when reading a missing program", func() {
			other := sha256.Sum256([]byte("missing."))
			_, err := fsys.ReadFile(publishedPath(publisher.String(), other))

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading an unpublished program for a publisher", func() {
			otherPublisher := sdk.AccAddress([]byte("publisher-address-02"))
			_, err := fsys.ReadFile(publishedPath(otherPublisher.String(), id))

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading with invalid file extension", func() {
			_, err := fsys.ReadFile(publisher.String() + "/abc")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading with invalid publisher", func() {
			_, err := fsys.ReadFile("not-a-bech32/" + idToPath(id))

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading without a publisher separator", func() {
			_, err := fsys.ReadFile(publisher.String())

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading with invalid id format", func() {
			_, err := fsys.ReadFile(publisher.String() + "/zzzz.pl")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading with invalid hexadecimal id content", func() {
			invalidID := strings.Repeat("z", sha256.Size*2)
			_, err := fsys.ReadFile(publisher.String() + "/" + invalidID + ".pl")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading with traversal path", func() {
			_, err := fsys.ReadFile("../secret.pl")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func TestUserFSOpen(t *testing.T) {
	Convey("Given a user library filesystem", t, func() {
		key := storetypes.NewKVStoreKey("test")
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

		source := "foo."
		id := sha256.Sum256([]byte(source))
		modTime := time.Date(2026, time.March, 9, 16, 0, 0, 0, time.UTC)
		publisher := sdk.AccAddress([]byte("publisher-address-03"))
		path := publishedPath(publisher.String(), id)

		fsys := NewFS(testCtx.Ctx, &testKeeper{
			programs: map[string]logictypes.StoredProgram{
				string(id[:]): {
					Source:     source,
					CreatedAt:  modTime.Add(-time.Hour).Unix(),
					SourceSize: uint64(len(source)),
				},
			},
			publications: map[string]logictypes.ProgramPublication{
				string(publisher) + ":" + string(id[:]): {
					PublishedAt: modTime.Unix(),
				},
			},
		})

		Convey("when opening an existing program file", func() {
			file, err := fsys.Open(path)
			So(err, ShouldBeNil)
			defer file.Close()

			info, err := file.Stat()
			So(err, ShouldBeNil)
			So(info.Name(), ShouldEqual, path)
			So(info.Size(), ShouldEqual, int64(len(source)))
			So(info.ModTime(), ShouldEqual, modTime)
		})
	})
}

func TestShareFSErrors(t *testing.T) {
	Convey("Given a shared logic filesystem", t, func() {
		key := storetypes.NewKVStoreKey("test")
		testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
		publisher := sdk.AccAddress([]byte("publisher-address-04"))
		source := "bar."
		id := sha256.Sum256([]byte(source))
		path := publishedPath(publisher.String(), id)

		Convey("when the filesystem has no keeper", func() {
			fsys := NewFS(testCtx.Ctx, nil)
			_, err := fsys.ReadFile(path)

			So(err, ShouldNotBeNil)
			So(errors.Is(err, errVFSUnavailable), ShouldBeTrue)
		})

		Convey("when opening an invalid path", func() {
			fsys := NewFS(testCtx.Ctx, nil)
			_, err := fsys.Open(path)

			So(err, ShouldNotBeNil)
			So(errors.Is(err, errVFSUnavailable), ShouldBeTrue)
		})

		Convey("when reading the root path", func() {
			fsys := NewFS(testCtx.Ctx, &testKeeper{})
			_, err := fsys.ReadFile(".")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when publication lookup fails", func() {
			publicationErr := errors.New("publication lookup failed")
			fsys := NewFS(testCtx.Ctx, &testKeeper{publicationErr: publicationErr})
			_, err := fsys.ReadFile(path)

			So(err, ShouldNotBeNil)
			So(errors.Is(err, publicationErr), ShouldBeTrue)
		})

		Convey("when stored program lookup fails", func() {
			programErr := errors.New("program lookup failed")
			fsys := NewFS(testCtx.Ctx, &testKeeper{
				publications: map[string]logictypes.ProgramPublication{
					string(publisher) + ":" + string(id[:]): {PublishedAt: 1},
				},
				programErr: programErr,
			})
			_, err := fsys.ReadFile(path)

			So(err, ShouldNotBeNil)
			So(errors.Is(err, programErr), ShouldBeTrue)
		})

		Convey("when a publication exists without the stored artifact", func() {
			fsys := NewFS(testCtx.Ctx, &testKeeper{
				publications: map[string]logictypes.ProgramPublication{
					string(publisher) + ":" + string(id[:]): {PublishedAt: 1},
				},
			})
			_, err := fsys.ReadFile(path)

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})
	})
}

func idToPath(id [sha256.Size]byte) string {
	return hex.EncodeToString(id[:]) + ".pl"
}

func publishedPath(publisher string, id [sha256.Size]byte) string {
	return publisher + "/" + idToPath(id)
}
