package keeper_test

import (
	"context"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/okp4/okp4d/testutil/keeper"
	"github.com/okp4/okp4d/x/knowledge/keeper"
	"github.com/okp4/okp4d/x/knowledge/types"
)

func setupMsgServer(t testing.TB) (types.MsgServer, context.Context) {
	k, ctx := keepertest.KnowledgeKeeper(t)
	return keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)
}
