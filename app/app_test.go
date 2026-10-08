package app

import (
	"fmt"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	dbm "github.com/cosmos/cosmos-db"

	. "github.com/smartystreets/goconvey/convey"

	"cosmossdk.io/log/v2"

	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
)

func TestMaxWasmSizeParsing(t *testing.T) {
	const defaultMaxWasmSize = 42
	originalWasmMaxSize := wasmtypes.MaxWasmSize
	t.Cleanup(func() { wasmtypes.MaxWasmSize = originalWasmMaxSize })
	originalMaxWasmSize := MaxWasmSize

	Convey("Given a test cases", t, func() {
		cases := []struct {
			name          string
			maxWasmSize   string
			expectedPanic bool
			expectedValue int
		}{
			{"empty string", "", false, defaultMaxWasmSize},
			{"valid number", "1048576", false, 1048576},
			{"invalid input", "not-a-number", true, defaultMaxWasmSize},
		}

		Convey(fmt.Sprintf("With default MaxWasmSize set to %d", defaultMaxWasmSize), func() {
			wasmtypes.MaxWasmSize = defaultMaxWasmSize

			for _, tc := range cases {
				Convey(fmt.Sprintf("When MaxWasmSize is '%s'", tc.maxWasmSize), func() {
					MaxWasmSize = tc.maxWasmSize

					Convey("Calling initialization should match expectations", func() {
						if tc.expectedPanic {
							So(mustConfigureWasmExtensionPoints, ShouldPanic)
						} else {
							mustConfigureWasmExtensionPoints()
							So(wasmtypes.MaxWasmSize, ShouldEqual, tc.expectedValue)
						}
					})

					Reset(func() {
						MaxWasmSize = originalMaxWasmSize
						wasmtypes.MaxWasmSize = originalWasmMaxSize
					})
				})
			}
		})
	})
}

func TestMakeEncodingConfigConstructsApp(t *testing.T) {
	Convey("Given the test encoding configuration", t, func() {
		encodingConfig := MakeEncodingConfig(t)

		So(encodingConfig.Codec, ShouldNotBeNil)
		So(encodingConfig.InterfaceRegistry, ShouldNotBeNil)
		So(encodingConfig.TxConfig, ShouldNotBeNil)
	})
}

func TestRetiredModulesAreDisabled(t *testing.T) {
	Convey("Given an initialized application", t, func() {
		application := New(
			log.NewNopLogger(),
			dbm.NewMemDB(),
			nil,
			true,
			simtestutil.NewAppOptionsWithFlagHome(t.TempDir()),
		)

		Convey("The group and capability modules should not be executable", func() {
			_, hasGroupModule := application.ModuleManager.Modules[legacyGroupStoreKey]
			_, hasCapabilityModule := application.ModuleManager.Modules[legacyCapabilityStoreKey]

			So(hasGroupModule, ShouldBeFalse)
			So(hasCapabilityModule, ShouldBeFalse)
		})

		Convey("Their legacy stores should remain mounted for data recovery", func() {
			So(application.GetKey(legacyGroupStoreKey), ShouldNotBeNil)
			So(application.GetKey(legacyCapabilityStoreKey), ShouldNotBeNil)
		})
	})
}
