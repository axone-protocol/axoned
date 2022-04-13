package simulation

import (
	"math/rand"

	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simtypes "github.com/cosmos/cosmos-sdk/types/simulation"
	"github.com/okp4/okp4d/x/knowledge/keeper"
	"github.com/okp4/okp4d/x/knowledge/types"
)

func SimulateMsgTriggerService(
	ak types.AccountKeeper,
	bk types.BankKeeper,
	k keeper.Keeper,
) simtypes.Operation {
	return func(r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context, accs []simtypes.Account, chainID string,
	) (simtypes.OperationMsg, []simtypes.FutureOperation, error) {
		simAccount, _ := simtypes.RandomAcc(r, accs)
		msg := &types.MsgTriggerService{
			Creator: simAccount.Address.String(),
		}

		// TODO: Handling the TriggerService simulation

		return simtypes.NoOpMsg(types.ModuleName, msg.Type(), "TriggerService simulation not implemented"), nil, nil
	}
}
