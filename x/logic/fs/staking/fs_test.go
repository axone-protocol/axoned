package staking

import (
	"context"
	"errors"
	"testing"
	"time"

	dbm "github.com/cosmos/cosmos-db"

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

	logictypes "github.com/axone-protocol/axoned/v15/x/logic/types"
)

const (
	delegatorAddress = "axone1ffd5wx65l407yvm478cxzlgygw07h79sw4jwpa"
	bondDenom        = "uaxone"
)

func TestVFS(t *testing.T) {
	Convey("Given a staking VFS", t, func() {
		sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")
		service := &stakingQueryServiceStub{}
		ctx := context.WithValue(newTestContext(), logictypes.StakingQueryServiceContextKey, service)
		stakingFS := NewFS(ctx)

		Convey("when reading paginated delegations", func() {
			largeAmount, ok := math.NewIntFromString("9223372036854775808")
			So(ok, ShouldBeTrue)
			service.delegations = func(
				_ context.Context, request *stakingtypes.QueryDelegatorDelegationsRequest,
			) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
				So(request.DelegatorAddr, ShouldEqual, delegatorAddress)
				So(request.Pagination.Limit, ShouldEqual, pageSize)
				switch string(request.Pagination.Key) {
				case "":
					return &stakingtypes.QueryDelegatorDelegationsResponse{
						DelegationResponses: []stakingtypes.DelegationResponse{{
							Delegation: stakingtypes.Delegation{ValidatorAddress: "axonevaloper1first"},
							Balance:    sdk.NewCoin("uaxone", math.NewInt(100)),
						}},
						Pagination: &query.PageResponse{NextKey: []byte{1}},
					}, nil
				case "\x01":
					return &stakingtypes.QueryDelegatorDelegationsResponse{
						Pagination: &query.PageResponse{NextKey: []byte{2}},
					}, nil
				case "\x02":
					return &stakingtypes.QueryDelegatorDelegationsResponse{
						DelegationResponses: []stakingtypes.DelegationResponse{{
							Delegation: stakingtypes.Delegation{ValidatorAddress: "axonevaloper1second"},
							Balance:    sdk.NewCoin("uaxone", largeAmount),
						}},
						Pagination: &query.PageResponse{},
					}, nil
				default:
					return nil, errors.New("unexpected pagination key")
				}
			}

			data, err := stakingFS.ReadFile(delegatorAddress + "/delegations/@")

			So(err, ShouldBeNil)
			So(string(data), ShouldContainSubstring, "delegation{balance:coin(uaxone,100),validator:axonevaloper1first}")
			So(string(data), ShouldContainSubstring, "delegation{balance:coin(uaxone,'9223372036854775808'),validator:axonevaloper1second}")
		})

		Convey("when reading an empty delegation collection", func() {
			service.delegations = func(
				context.Context, *stakingtypes.QueryDelegatorDelegationsRequest,
			) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
				return &stakingtypes.QueryDelegatorDelegationsResponse{Pagination: &query.PageResponse{}}, nil
			}

			data, err := stakingFS.ReadFile(delegatorAddress + "/delegations/@")

			So(err, ShouldBeNil)
			So(data, ShouldBeEmpty)
		})

		Convey("when reading unbonding delegations", func() {
			service.params = func(context.Context, *stakingtypes.QueryParamsRequest) (*stakingtypes.QueryParamsResponse, error) {
				return &stakingtypes.QueryParamsResponse{Params: stakingtypes.Params{BondDenom: bondDenom}}, nil
			}
			service.unbondingDelegations = func(
				_ context.Context, request *stakingtypes.QueryDelegatorUnbondingDelegationsRequest,
			) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
				So(request.DelegatorAddr, ShouldEqual, delegatorAddress)
				So(request.Pagination.Limit, ShouldEqual, pageSize)
				return &stakingtypes.QueryDelegatorUnbondingDelegationsResponse{
					UnbondingResponses: []stakingtypes.UnbondingDelegation{{
						ValidatorAddress: "axonevaloper1validator",
						Entries: []stakingtypes.UnbondingDelegationEntry{{
							Balance:        math.NewInt(75),
							CompletionTime: time.Unix(1_800_000_000, 0).UTC(),
						}},
					}, {ValidatorAddress: "axonevaloper1empty"}},
					Pagination: &query.PageResponse{},
				}, nil
			}

			data, err := stakingFS.ReadFile(delegatorAddress + "/unbonding_delegations/@")

			So(err, ShouldBeNil)
			So(string(data), ShouldContainSubstring, "unbonding_entry{balance:coin(uaxone,75),completion_time:1800000000}")
			So(string(data), ShouldContainSubstring, "entries:[]")
		})

		Convey("when reading redelegations", func() {
			service.params = func(context.Context, *stakingtypes.QueryParamsRequest) (*stakingtypes.QueryParamsResponse, error) {
				return &stakingtypes.QueryParamsResponse{Params: stakingtypes.Params{BondDenom: bondDenom}}, nil
			}
			service.redelegations = func(
				_ context.Context, request *stakingtypes.QueryRedelegationsRequest,
			) (*stakingtypes.QueryRedelegationsResponse, error) {
				So(request.DelegatorAddr, ShouldEqual, delegatorAddress)
				So(request.Pagination.Limit, ShouldEqual, pageSize)
				return &stakingtypes.QueryRedelegationsResponse{
					RedelegationResponses: []stakingtypes.RedelegationResponse{{
						Redelegation: stakingtypes.Redelegation{
							ValidatorSrcAddress: "axonevaloper1source",
							ValidatorDstAddress: "axonevaloper1destination",
						},
						Entries: []stakingtypes.RedelegationEntryResponse{{
							RedelegationEntry: stakingtypes.RedelegationEntry{CompletionTime: time.Unix(1_800_000_001, 0).UTC()},
							Balance:           math.NewInt(55),
						}},
					}, {
						Redelegation: stakingtypes.Redelegation{
							ValidatorSrcAddress: "axonevaloper1empty",
							ValidatorDstAddress: "axonevaloper1emptyto",
						},
					}},
					Pagination: &query.PageResponse{},
				}, nil
			}

			data, err := stakingFS.ReadFile(delegatorAddress + "/redelegations/@")

			So(err, ShouldBeNil)
			So(string(data), ShouldContainSubstring, "redelegation_entry{balance:coin(uaxone,55),completion_time:1800000001}")
			So(string(data), ShouldContainSubstring, "source_validator:axonevaloper1source")
			So(string(data), ShouldContainSubstring, "destination_validator:axonevaloper1destination")
			So(string(data), ShouldContainSubstring, "entries:[]")
		})

		Convey("when reading an invalid path", func() {
			for _, name := range []string{
				"../" + delegatorAddress + "/delegations/@",
				"not-an-address/delegations/@",
				delegatorAddress + "/delegations",
				delegatorAddress + "/unknown/@",
			} {
				_, err := stakingFS.ReadFile(name)
				So(err, ShouldNotBeNil)
			}
		})
	})
}

func TestVFSQueryErrors(t *testing.T) {
	Convey("Given staking query failures", t, func() {
		sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")
		queryErr := errors.New("query failed")
		service := &stakingQueryServiceStub{
			delegations: func(
				context.Context, *stakingtypes.QueryDelegatorDelegationsRequest,
			) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
				return nil, queryErr
			},
			unbondingDelegations: func(
				context.Context, *stakingtypes.QueryDelegatorUnbondingDelegationsRequest,
			) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
				return nil, queryErr
			},
			redelegations: func(
				context.Context, *stakingtypes.QueryRedelegationsRequest,
			) (*stakingtypes.QueryRedelegationsResponse, error) {
				return nil, queryErr
			},
			params: func(
				context.Context, *stakingtypes.QueryParamsRequest,
			) (*stakingtypes.QueryParamsResponse, error) {
				return &stakingtypes.QueryParamsResponse{Params: stakingtypes.Params{BondDenom: bondDenom}}, nil
			},
		}
		stakingFS := NewFS(context.WithValue(newTestContext(), logictypes.StakingQueryServiceContextKey, service))

		_, err := stakingFS.ReadFile(delegatorAddress + "/delegations/@")
		So(errors.Is(err, queryErr), ShouldBeTrue)

		_, err = stakingFS.ReadFile(delegatorAddress + "/unbonding_delegations/@")
		So(errors.Is(err, queryErr), ShouldBeTrue)

		_, err = stakingFS.ReadFile(delegatorAddress + "/redelegations/@")
		So(errors.Is(err, queryErr), ShouldBeTrue)

		paramsErr := errors.New("params query failed")
		service.params = func(context.Context, *stakingtypes.QueryParamsRequest) (*stakingtypes.QueryParamsResponse, error) {
			return nil, paramsErr
		}
		_, err = stakingFS.Open(delegatorAddress + "/unbonding_delegations/@")
		So(errors.Is(err, paramsErr), ShouldBeTrue)
		_, err = stakingFS.Open(delegatorAddress + "/redelegations/@")
		So(errors.Is(err, paramsErr), ShouldBeTrue)
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

type stakingQueryServiceStub struct {
	delegations func(
		context.Context, *stakingtypes.QueryDelegatorDelegationsRequest,
	) (*stakingtypes.QueryDelegatorDelegationsResponse, error)
	unbondingDelegations func(
		context.Context, *stakingtypes.QueryDelegatorUnbondingDelegationsRequest,
	) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error)
	redelegations func(
		context.Context, *stakingtypes.QueryRedelegationsRequest,
	) (*stakingtypes.QueryRedelegationsResponse, error)
	params func(
		context.Context, *stakingtypes.QueryParamsRequest,
	) (*stakingtypes.QueryParamsResponse, error)
}

func (s *stakingQueryServiceStub) DelegatorDelegations(
	ctx context.Context, request *stakingtypes.QueryDelegatorDelegationsRequest,
) (*stakingtypes.QueryDelegatorDelegationsResponse, error) {
	if s.delegations == nil {
		return nil, errors.New("unexpected delegations query")
	}
	return s.delegations(ctx, request)
}

func (s *stakingQueryServiceStub) DelegatorUnbondingDelegations(
	ctx context.Context, request *stakingtypes.QueryDelegatorUnbondingDelegationsRequest,
) (*stakingtypes.QueryDelegatorUnbondingDelegationsResponse, error) {
	if s.unbondingDelegations == nil {
		return nil, errors.New("unexpected unbonding delegations query")
	}
	return s.unbondingDelegations(ctx, request)
}

func (s *stakingQueryServiceStub) Redelegations(
	ctx context.Context, request *stakingtypes.QueryRedelegationsRequest,
) (*stakingtypes.QueryRedelegationsResponse, error) {
	if s.redelegations == nil {
		return nil, errors.New("unexpected redelegations query")
	}
	return s.redelegations(ctx, request)
}

func (s *stakingQueryServiceStub) Params(
	ctx context.Context, request *stakingtypes.QueryParamsRequest,
) (*stakingtypes.QueryParamsResponse, error) {
	if s.params == nil {
		return nil, errors.New("unexpected params query")
	}
	return s.params(ctx, request)
}
