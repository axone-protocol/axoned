package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/okp4/okp4d/x/mint/types"
)

func (s *IntegrationTestSuite) TestUpdateParams() {
	testCases := []struct {
		name      string
		request   *types.MsgUpdateParams
		expectErr bool
	}{
		{
			name: "set invalid authority",
			request: &types.MsgUpdateParams{
				Authority: "foo",
			},
			expectErr: true,
		},
		{
			name: "set invalid params",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:             sdk.DefaultBondDenom,
					AnnualReductionFactor: sdk.NewDecWithPrec(-13, 2),
					BlocksPerYear:         uint64(60 * 60 * 8766 / 5),
				},
			},
			expectErr: true,
		},
		{
			name: "set full valid params",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:             sdk.DefaultBondDenom,
					AnnualReductionFactor: sdk.NewDecWithPrec(8, 2),
					BlocksPerYear:         uint64(60 * 60 * 8766 / 5),
				},
			},
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			_, err := s.msgServer.UpdateParams(s.ctx, tc.request)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}
