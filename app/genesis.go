package app

import (
	"encoding/json"
	"testing"

	"cosmossdk.io/log"

	dbm "github.com/cosmos/cosmos-db"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"

	"github.com/cosmos/cosmos-sdk/codec"
)

// The genesis state of the blockchain is represented here as a map of raw json
// messages key'd by a identifier string.
// The identifier is used to determine which module genesis information belongs
// to so it may be appropriately routed during init chain.
// Within this application default genesis information is retrieved from
// the ModuleBasicManager which populates json from each BasicModule
// object provided to it during init.
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application, for testing purpose.
func NewDefaultGenesisState(t testing.TB, cdc codec.JSONCodec) GenesisState {
	tempApp := New(log.NewNopLogger(), dbm.NewMemDB(), nil, true, simtestutil.NewAppOptionsWithFlagHome(t.TempDir()))
	return tempApp.BasicModuleManager.DefaultGenesis(cdc)
}
