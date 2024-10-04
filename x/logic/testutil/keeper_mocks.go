package testutil

import (
	context "context"
	"fmt"
	"strconv"

	"github.com/golang/mock/gomock"
	"github.com/samber/lo"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func MockAuthQueryServiceWithAddresses(mock *MockAuthQueryService, addresses []string) {
	total := len(addresses)
	mock.
		EXPECT().
		Accounts(gomock.Any(), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(_ context.Context, req *authtypes.QueryAccountsRequest) (*authtypes.QueryAccountsResponse, error) {
			start := 0
			limit := 5
			toCursor := func(idx int) []byte { return []byte(fmt.Sprintf("%d", idx)) }
			fromCursor := func(k []byte) int {
				idx, err := strconv.Atoi(string(k))
				if err != nil {
					panic(err)
				}

				return idx
			}

			if req.Pagination != nil {
				if req.Pagination.Key != nil {
					start = fromCursor(req.Pagination.Key)
				}
				if req.Pagination.Limit != 0 {
					limit = int(req.Pagination.GetLimit())
				}
			}
			accounts := lo.Map(
				lo.Slice(addresses, start, start+limit),
				func(acc string, _ int) *codectypes.Any {
					addr, err := sdk.AccAddressFromBech32(acc)
					if err != nil {
						panic(err)
					}

					accI := authtypes.ProtoBaseAccount()
					err = accI.SetAddress(addr)
					if err != nil {
						panic(err)
					}

					anyV, err := codectypes.NewAnyWithValue(accI)
					if err != nil {
						panic(err)
					}

					return anyV
				})

			return &authtypes.QueryAccountsResponse{
				Accounts: accounts,
				Pagination: &query.PageResponse{
					NextKey: lo.If(start+limit < total, toCursor(start+1)).Else(nil),
					Total:   uint64(total),
				},
			}, nil
		})
}

func MockAuthQueryServiceWithError(mock *MockAuthQueryService, err error) {
	mock.
		EXPECT().
		Accounts(gomock.Any(), gomock.Any()).
		AnyTimes().
		DoAndReturn(func(_ context.Context, _ *authtypes.QueryAccountsRequest) (*authtypes.QueryAccountsResponse, error) {
			return nil, err
		})
}
