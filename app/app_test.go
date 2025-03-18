package app

import (
	"fmt"
	"testing"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMaxWasmSizeParsing(t *testing.T) {
	const defaultMaxWasmSize = 42

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
						wasmtypes.MaxWasmSize = defaultMaxWasmSize
					})
				})
			}
		})
	})
}
