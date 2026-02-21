package comet

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	corecomet "cosmossdk.io/core/comet"
	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestSysCometVFSReadFile(t *testing.T) {
	Convey("Given a sys/comet VFS", t, func() {
		evidenceTime := time.Date(2024, 4, 9, 9, 10, 11, 0, time.UTC)
		expectedAt := []byte(
			"comet{evidence:[evidence{height:11,time:1712653811,total_voting_power:12,type:1,validator:validator{address:[7,8],power:9}}]," +
				"last_commit:commit_info{round:3,votes:[vote_info{block_id_flag:2,validator:validator{address:[10,11],power:12}}]}," +
				"proposer_address:[4,5,6],validators_hash:[1,2,3]}.\n")
		expectedEvidence := []byte(
			"[evidence{height:11,time:1712653811,total_voting_power:12,type:1,validator:validator{address:[7,8],power:9}}].\n")
		expectedLastCommit := []byte(
			"commit_info{round:3,votes:[vote_info{block_id_flag:2,validator:validator{address:[10,11],power:12}}]}.\n")

		vfs := NewFS(newTestContext(
			coreheader.Info{Time: time.Date(2024, 4, 10, 10, 44, 27, 0, time.UTC)},
			baseapp.NewBlockInfo(
				[]abci.Misbehavior{{
					Type: abci.MisbehaviorType_DUPLICATE_VOTE,
					Validator: abci.Validator{
						Address: []byte{7, 8},
						Power:   9,
					},
					Height:           11,
					Time:             evidenceTime,
					TotalVotingPower: 12,
				}},
				[]byte{1, 2, 3},
				[]byte{4, 5, 6},
				abci.CommitInfo{
					Round: 3,
					Votes: []abci.VoteInfo{{
						Validator: abci.Validator{
							Address: []byte{10, 11},
							Power:   12,
						},
						BlockIdFlag: cmtproto.BlockIDFlagCommit,
					}},
				},
			),
		))

		cases := []struct {
			path string
			want []byte
		}{
			{path: "@", want: expectedAt},
			{path: "validators_hash", want: []byte("[1,2,3].\n")},
			{path: "proposer_address", want: []byte("[4,5,6].\n")},
			{path: "evidence", want: expectedEvidence},
			{path: "last_commit", want: expectedLastCommit},
			{path: "last_commit/round", want: []byte("3.\n")},
			{path: "last_commit/votes", want: []byte("[vote_info{block_id_flag:2,validator:validator{address:[10,11],power:12}}].\n")},
		}

		for i, tc := range cases {
			Convey(fmt.Sprintf("when reading case #%d path %s", i, tc.path), func() {
				got, err := vfs.ReadFile(tc.path)

				So(err, ShouldBeNil)
				So(got, ShouldResemble, tc.want)
			})
		}
	})
}

func TestSysCometVFSOpen(t *testing.T) {
	Convey("Given a sys/comet VFS", t, func() {
		headerTime := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
		vfs := NewFS(newTestContext(
			coreheader.Info{Time: headerTime},
			nil,
		))

		Convey("when opening validators_hash", func() {
			f, err := vfs.Open("validators_hash")

			So(err, ShouldBeNil)
			defer f.Close()

			info, err := f.Stat()
			So(err, ShouldBeNil)
			So(info.Name(), ShouldEqual, "validators_hash")
			So(info.ModTime(), ShouldEqual, headerTime)

			content, err := io.ReadAll(f)
			So(err, ShouldBeNil)
			So(content, ShouldResemble, []byte("[].\n"))
		})
	})
}

func TestSysCometVFSErrors(t *testing.T) {
	Convey("Given a sys/comet VFS", t, func() {
		vfs := NewFS(newTestContext(coreheader.Info{}, nil))

		Convey("when reading an unknown path", func() {
			_, err := vfs.ReadFile("unknown")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrNotExist), ShouldBeTrue)
		})

		Convey("when reading a path escaping root", func() {
			_, err := vfs.ReadFile("../validators_hash")

			So(err, ShouldNotBeNil)
			So(errors.Is(err, fs.ErrPermission), ShouldBeTrue)
		})
	})
}

func newTestContext(headerInfo coreheader.Info, cometInfo corecomet.BlockInfo) sdk.Context {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())

	ctx := sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger()).WithHeaderInfo(headerInfo)
	if cometInfo != nil {
		ctx = ctx.WithCometInfo(cometInfo)
	}

	return ctx
}
