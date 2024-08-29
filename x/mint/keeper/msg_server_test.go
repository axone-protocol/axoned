package keeper_test

import (
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v10/x/mint/types"
)

func (s *IntegrationTestSuite) TestUpdateParams() {
	validInfMin := math.LegacyNewDecWithPrec(3, 2)
	validInfMax := math.LegacyNewDecWithPrec(20, 2)
	invalidInfMin := math.LegacyNewDecWithPrec(-3, 2)
	invalidInfMax := math.LegacyNewDecWithPrec(-20, 2)

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
			name: "set invalid params for inflation coef (negative value)",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:     sdk.DefaultBondDenom,
					InflationCoef: math.LegacyNewDecWithPrec(-73, 2),
					BlocksPerYear: uint64(60 * 60 * 8766 / 5),
				},
			},
			expectErr: true,
		},
		{
			name: "set invalid params for inflation max (negative value)",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:     sdk.DefaultBondDenom,
					InflationCoef: math.LegacyNewDecWithPrec(73, 2),
					InflationMax:  &invalidInfMax,
					InflationMin:  nil,
					BlocksPerYear: uint64(60 * 60 * 8766 / 5),
				},
			},
			expectErr: true,
		},
		{
			name: "set invalid params for inflation min (negative value)",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:     sdk.DefaultBondDenom,
					InflationCoef: math.LegacyNewDecWithPrec(73, 2),
					InflationMax:  nil,
					InflationMin:  &invalidInfMin,
					BlocksPerYear: uint64(60 * 60 * 8766 / 5),
				},
			},
			expectErr: true,
		},
		{
			name: "set invalid params for inflation min & max (min > max)",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:     sdk.DefaultBondDenom,
					InflationCoef: math.LegacyNewDecWithPrec(73, 2),
					InflationMax:  &validInfMin,
					InflationMin:  &validInfMax,
					BlocksPerYear: uint64(60 * 60 * 8766 / 5),
				},
			},
			expectErr: true,
		},
		{
			name: "set full valid params with boundaries",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:     sdk.DefaultBondDenom,
					InflationCoef: math.LegacyNewDecWithPrec(73, 2),
					InflationMax:  &validInfMax,
					InflationMin:  &validInfMin,
					BlocksPerYear: uint64(60 * 60 * 8766 / 5),
				},
			},
			expectErr: false,
		},
		{
			name: "set full valid params without boundaries",
			request: &types.MsgUpdateParams{
				Authority: s.mintKeeper.GetAuthority(),
				Params: types.Params{
					MintDenom:     sdk.DefaultBondDenom,
					InflationCoef: math.LegacyNewDecWithPrec(73, 2),
					InflationMax:  nil,
					InflationMin:  nil,
					BlocksPerYear: uint64(60 * 60 * 8766 / 5),
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
