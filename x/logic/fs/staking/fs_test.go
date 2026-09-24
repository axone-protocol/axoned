package staking

import (
	"context"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"

	coreheader "cosmossdk.io/core/header"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/axone-protocol/axoned/v15/x/logic/testutil"
	logictypes "github.com/axone-protocol/axoned/v15/x/logic/types"
)

const delegatorAddress = "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa"

func TestVFS(t *testing.T) {
	Convey("Given a staking VFS", t, func() {
		sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")
		ctrl := gomock.NewController(t)
		service := testutil.NewMockStakingQueryService(ctrl)
		ctx := context.WithValue(newTestContext(), logictypes.StakingQueryServiceContextKey, service)
		stakingFS := NewFS(ctx)

		Convey("when reading paginated delegations", func() {
			service.EXPECT().DelegatorDelegations(gomock.Any(), gomock.Any()).DoAndReturn(
				func(
					_ context.Context, request *stakingtypes.QueryDelegatorDelegationsRequest,
				) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
					if len(request.Pagination.Key) == 0 {
						return &stakingtypes.QueryDelegatorDelegationsResponse{
							DelegationResponses: []stakingtypes.DelegationResponse{{
								Delegation: stakingtypes.Delegation{ValidatorAddress: "axonevaloper1first"},
								Balance:    sdk.NewCoin("uaxone", math.NewInt(100)),
							}},
							Pagination: &query.PageResponse{NextKey: []byte{1}},
						}, nil
					}
					return &stakingtypes.QueryDelegatorDelegationsResponse{
						DelegationResponses: []stakingtypes.DelegationResponse{{
							Delegation: stakingtypes.Delegation{ValidatorAddress: "axonevaloper1second"},
							Balance:    sdk.NewCoin("uaxone", math.NewInt(200)),
						}},
						Pagination: &query.PageResponse{},
					}, nil
				},
			).Times(2)

			data, err := stakingFS.ReadFile(delegatorAddress + "/delegations/@")

			So(err, ShouldBeNil)
			So(string(data), ShouldContainSubstring, "delegation{balance:coin(uaxone,100),validator:axonevaloper1first}")
			So(string(data), ShouldContainSubstring, "delegation{balance:coin(uaxone,200),validator:axonevaloper1second}")
		})

		Convey("when reading unbonding delegations", func() {
			service.EXPECT().Params(gomock.Any(), gomock.Any()).Return(
				&stakingtypes.QueryParamsResponse{Params: stakingtypes.Params{BondDenom: "uaxone"}}, nil,
			)
			service.EXPECT().DelegatorUnbondingDelegations(gomock.Any(), gomock.Any()).Return(
				&stakingtypes.QueryDelegatorUnbondingDelegationsResponse{
					UnbondingResponses: []stakingtypes.UnbondingDelegation{{
						ValidatorAddress: "axonevaloper1validator",
						Entries: []stakingtypes.UnbondingDelegationEntry{{
							Balance:        math.NewInt(75),
							CompletionTime: time.Unix(1_800_000_000, 0).UTC(),
						}},
					}},
					Pagination: &query.PageResponse{},
				},
				nil,
			)

			data, err := stakingFS.ReadFile(delegatorAddress + "/unbonding_delegations/@")

			So(err, ShouldBeNil)
			So(string(data), ShouldContainSubstring, "unbonding_entry{balance:coin(uaxone,75),completion_time:1800000000}")
		})

		Convey("when reading an invalid path", func() {
			_, err := stakingFS.ReadFile(delegatorAddress + "/unknown/@")

			So(err, ShouldNotBeNil)
		})
	})
}

func TestVFSUnavailable(t *testing.T) {
	Convey("Given a staking VFS without a staking query service", t, func() {
		sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")
		stakingFS := NewFS(newTestContext())

		_, err := stakingFS.Open(delegatorAddress + "/delegations/@")

		So(err, ShouldNotBeNil)
	})
}

func newTestContext() sdk.Context {
	db := dbm.NewMemDB()
	stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
	return sdk.NewContext(stateStore, cmtproto.Header{}, false, log.NewNopLogger()).WithHeaderInfo(coreheader.Info{
		Height: 42,
		Time:   time.Date(2024, time.April, 10, 10, 44, 27, 0, time.UTC),
	})
}
