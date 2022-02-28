package keeper_test

import (
	"testing"

	testkeeper "github.com/okp4/okp4d/testutil/keeper"
	"github.com/okp4/okp4d/x/knowledge/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx := testkeeper.KnowledgeKeeper(t)
	params := types.DefaultParams()

	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}
