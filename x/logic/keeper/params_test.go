package keeper_test

import (
	"testing"

	testkeeper "github.com/okp4/okp4d/testutil/keeper"
	"github.com/okp4/okp4d/x/logic/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.LogicKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
	require.EqualValues(t, params.Foo, k.Foo(ctx))
}
