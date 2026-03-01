package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protowire"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v14/x/logic"
	"github.com/axone-protocol/axoned/v14/x/logic/keeper"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

func TestMigrator_Migrate4to5(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		populateLegacyFields bool
	}{
		{
			name:                 "legacy interpreter fields empty",
			populateLegacyFields: false,
		},
		{
			name:                 "legacy interpreter fields populated",
			populateLegacyFields: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			encCfg := moduletestutil.MakeTestEncodingConfig(logic.AppModuleBasic{})
			key := storetypes.NewKVStoreKey(types.StoreKey)
			testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))

			logicKeeper := keeper.NewKeeper(
				encCfg.Codec,
				encCfg.InterfaceRegistry,
				key,
				key,
				authtypes.NewModuleAddress(govtypes.ModuleName),
				nil,
				nil,
				nil,
				nil,
			)

			expectedParams := types.NewParams(
				types.NewLimits(
					types.WithMaxSize(11),
					types.WithMaxResultCount(7),
					types.WithMaxUserOutputSize(13),
					types.WithMaxVariables(17),
				),
				types.NewGasPolicy(
					types.WithWeightingFactor(19),
					types.WithDefaultPredicateCost(23),
					types.WithPredicateCosts([]types.PredicateCost{
						{Predicate: "consult/1", Cost: 29},
					}),
				),
			)

			expectedBz, err := encCfg.Codec.Marshal(&expectedParams)
			require.NoError(t, err)

			legacyBz := legacyParamsBytes(expectedBz, tc.populateLegacyFields)
			require.NotEqual(t, expectedBz, legacyBz)

			store := testCtx.Ctx.KVStore(key)
			store.Set(types.ParamsKey, legacyBz)

			err = keeper.NewMigrator(*logicKeeper).Migrate4to5(testCtx.Ctx)
			require.NoError(t, err)

			require.Equal(t, expectedBz, store.Get(types.ParamsKey))
			require.Equal(t, expectedParams, logicKeeper.GetParams(testCtx.Ctx))
		})
	}
}

func legacyParamsBytes(canonical []byte, populateLegacyFields bool) []byte {
	var interpreter []byte
	if populateLegacyFields {
		interpreter = appendLegacyFilter(interpreter, 1, []string{"consult/1"}, []string{"open/4"})
		interpreter = protowire.AppendTag(interpreter, 3, protowire.BytesType)
		interpreter = protowire.AppendString(interpreter, "user_bootstrap.")
		interpreter = appendLegacyFilter(interpreter, 4, []string{"cosmwasm:"}, []string{"https://"})
	}

	legacy := protowire.AppendTag(nil, 1, protowire.BytesType)
	legacy = protowire.AppendBytes(legacy, interpreter)
	legacy = append(legacy, canonical...)

	return legacy
}

func appendLegacyFilter(dst []byte, fieldNum protowire.Number, whitelist, blacklist []string) []byte {
	var filter []byte
	for _, item := range whitelist {
		filter = protowire.AppendTag(filter, 1, protowire.BytesType)
		filter = protowire.AppendString(filter, item)
	}
	for _, item := range blacklist {
		filter = protowire.AppendTag(filter, 2, protowire.BytesType)
		filter = protowire.AppendString(filter, item)
	}

	dst = protowire.AppendTag(dst, fieldNum, protowire.BytesType)
	dst = protowire.AppendBytes(dst, filter)

	return dst
}
