package keeper_test

import (
	"testing"

	"google.golang.org/protobuf/encoding/protowire"

	. "github.com/smartystreets/goconvey/convey"

	storetypes "cosmossdk.io/store/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/testutil"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axone-protocol/axoned/v14/x/logic"
	"github.com/axone-protocol/axoned/v14/x/logic/keeper"
	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

func TestMigrator_Migrate4to5(t *testing.T) {
	Convey("Given legacy params payloads", t, func() {
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
			Convey(tc.name, func() {
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
					types.DefaultGasPolicy(),
				)

				expectedBz, err := encCfg.Codec.Marshal(&expectedParams)
				So(err, ShouldBeNil)

				legacyBz := legacyParamsBytes(encCfg.Codec, expectedParams.Limits, tc.populateLegacyFields)
				So(legacyBz, ShouldNotResemble, expectedBz)

				store := testCtx.Ctx.KVStore(key)
				store.Set(types.ParamsKey, legacyBz)

				err = keeper.NewMigrator(*logicKeeper).Migrate4to5(testCtx.Ctx)
				So(err, ShouldBeNil)

				So(store.Get(types.ParamsKey), ShouldResemble, expectedBz)
				So(logicKeeper.GetParams(testCtx.Ctx), ShouldResemble, expectedParams)
			})
		}
	})
}

func legacyParamsBytes(codec codec.BinaryCodec, limits types.Limits, populateLegacyFields bool) []byte {
	var interpreter []byte
	if populateLegacyFields {
		interpreter = appendLegacyFilter(interpreter, 1, []string{"consult/1"}, []string{"open/4"})
		interpreter = protowire.AppendTag(interpreter, 3, protowire.BytesType)
		interpreter = protowire.AppendString(interpreter, "user_bootstrap.")
		interpreter = appendLegacyFilter(interpreter, 4, []string{"cosmwasm:"}, []string{"https://"})
	}

	var gasPolicy []byte
	if populateLegacyFields {
		gasPolicy = protowire.AppendTag(gasPolicy, 1, protowire.VarintType)
		gasPolicy = protowire.AppendVarint(gasPolicy, 19)
		gasPolicy = protowire.AppendTag(gasPolicy, 2, protowire.VarintType)
		gasPolicy = protowire.AppendVarint(gasPolicy, 23)
		gasPolicy = appendLegacyPredicateCost(gasPolicy, "consult/1", 29)
	}

	limitsBz, err := codec.Marshal(&limits)
	if err != nil {
		panic(err)
	}

	legacy := protowire.AppendTag(nil, 1, protowire.BytesType)
	legacy = protowire.AppendBytes(legacy, interpreter)
	legacy = protowire.AppendTag(legacy, 2, protowire.BytesType)
	legacy = protowire.AppendBytes(legacy, limitsBz)
	legacy = protowire.AppendTag(legacy, 3, protowire.BytesType)
	legacy = protowire.AppendBytes(legacy, gasPolicy)

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

func appendLegacyPredicateCost(dst []byte, predicate string, cost uint64) []byte {
	var predicateCost []byte
	predicateCost = protowire.AppendTag(predicateCost, 1, protowire.BytesType)
	predicateCost = protowire.AppendString(predicateCost, predicate)
	predicateCost = protowire.AppendTag(predicateCost, 2, protowire.VarintType)
	predicateCost = protowire.AppendVarint(predicateCost, cost)

	dst = protowire.AppendTag(dst, 3, protowire.BytesType)
	dst = protowire.AppendBytes(dst, predicateCost)

	return dst
}
