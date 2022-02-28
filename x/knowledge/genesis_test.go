package knowledge_test

import (
	"testing"

	keepertest "github.com/okp4/okp4d/testutil/keeper"
	"github.com/okp4/okp4d/testutil/nullify"
	"github.com/okp4/okp4d/x/knowledge"
	"github.com/okp4/okp4d/x/knowledge/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.KnowledgeKeeper(t)
	knowledge.InitGenesis(ctx, *k, genesisState)
	got := knowledge.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
