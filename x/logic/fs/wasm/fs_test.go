//nolint:lll
package wasm

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"testing"

	dbm "github.com/cosmos/cosmos-db"
	"go.uber.org/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"cosmossdk.io/log"
	"cosmossdk.io/store"
	"cosmossdk.io/store/metrics"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axone-protocol/axoned/v14/x/logic/testutil"
)

//nolint:gocognit
func TestWasmVFS(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Given a test cases", t, func() {
		contractAddress := "axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk"
		cases := []struct {
			contractAddress string
			query           []byte
			data            []byte
			uri             string
			fail            bool
			wantResult      []byte
			wantError       string
		}{
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `cosmwasm:cw-storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("calc(X) :-  X is 100 + 200."),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu"),
				uri:             `cosmwasm:cw-storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Sprintf("cosmwasm:cw-storage:%s?query=%%7B%%22object_data%%22%%3A%%7B%%22id%%22%%3A%%20%%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%%22%%7D%%7D: failed to unmarshal JSON WASM response to string: invalid character 'Y' looking for beginning of value", contractAddress),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `cosmwasm:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("calc(X) :-  X is 100 + 200."),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `axone:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantError:       fmt.Sprintf("axone:%s?query=%%7B%%22object_data%%22%%3A%%7B%%22id%%22%%3A%%20%%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%%22%%7D%%7D: invalid scheme, expected 'cosmwasm', got 'axone'", contractAddress),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:cw-storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Sprintf("cosmwasm:cw-storage:%s?query=%%7B%%22object_data%%22%%3A%%7B%%22id%%22%%3A%%20%%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%%22%%7D%%7D: failed to decode WASM base64 response: illegal base64 data at input byte 0", contractAddress),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:cw-storage?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Sprintf("cosmwasm:cw-storage?query=%%7B%%22object_data%%22%%3A%%7B%%22id%%22%%3A%%20%%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%%22%%7D%%7D: failed to convert path 'cw-storage' to contract address: decoding bech32 failed: invalid separator index -1"),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:cw-storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?wasm=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Sprintf("cosmwasm:cw-storage:%s?wasm=%%7B%%22object_data%%22%%3A%%7B%%22id%%22%%3A%%20%%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%%22%%7D%%7D: uri should contains `query` params", contractAddress),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Sprintf("cosmwasm:?query=%%7B%%22object_data%%22%%3A%%7B%%22id%%22%%3A%%20%%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%%22%%7D%%7D: empty path given, should be 'cosmwasm:{contractName}:{contractAddr}?query={query}'"),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("foo-bar"),
				uri:             `cosmwasm:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=false`,
				wantResult:      []byte("foo-bar"),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `cosmwasm:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=true`,
				wantResult:      []byte("calc(X) :-  X is 100 + 200."),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=foo`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Sprintf(`cosmwasm:%s?query=%%7B%%22object_data%%22%%3A%%7B%%22id%%22%%3A%%20%%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%%22%%7D%%7D&base64Decode=foo: failed to convert 'base64Decode' query value to boolean: strconv.ParseBool: parsing "foo": invalid syntax`, contractAddress),
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `% %`,
				wantResult:      []byte("\"\""),
				wantError:       "% %: invalid argument",
			},
			{
				contractAddress: contractAddress,
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `cosmwasm:cw-storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				fail:            true,
				wantError:       "cosmwasm:cw-storage:axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D: failed to query WASM contract axone15ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3tsrhsdrk: failed to query smart contract",
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the uri #%d: %s", nc, tc.uri), func() {
				Convey("and a wasm keeper initialized with the given values", func() {
					db := dbm.NewMemDB()
					stateStore := store.NewCommitMultiStore(db, log.NewNopLogger(), metrics.NewNoOpMetrics())
					wasmKeeper := testutil.NewMockWasmKeeper(ctrl)
					ctx := sdk.
						NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
					sdk.GetConfig().SetBech32PrefixForAccount("axone", "axonepub")

					wasmKeeper.EXPECT().
						QuerySmart(ctx, sdk.MustAccAddressFromBech32(tc.contractAddress), tc.query).
						AnyTimes().
						DoAndReturn(func(_, _, _ any) ([]byte, error) {
							if tc.fail {
								return nil, errors.New("failed to query smart contract")
							}

							return tc.data, nil
						})

					Convey("and a wasm file system under test", func() {
						vfs := NewFS(ctx, wasmKeeper)

						Convey(fmt.Sprintf(`when the open("%s") is called`, tc.uri), func() {
							file, err := vfs.Open(tc.uri)

							Convey("then the result should be as expected", func() {
								if tc.wantError != "" {
									So(err, ShouldNotBeNil)
									So(err.Error(), ShouldEqual, fmt.Sprintf("open %s", tc.wantError))
								} else {
									So(err, ShouldBeNil)

									defer file.Close()
									info, err := file.Stat()
									So(err, ShouldBeNil)

									So(info.Name(), ShouldEqual, tc.uri)
									So(info.Size(), ShouldEqual, int64(len(tc.wantResult)))
									So(info.ModTime(), ShouldEqual, ctx.BlockTime())
									So(info.IsDir(), ShouldBeFalse)
									So(info.Mode(), ShouldEqual, fs.ModeIrregular)
									So(info.Sys(), ShouldBeNil)

									data := make([]byte, info.Size())
									for {
										_, err := file.Read(data)
										if errors.Is(err, io.EOF) {
											break
										}
										continue
									}

									So(data, ShouldResemble, tc.wantResult)
								}
							})

							Convey(fmt.Sprintf(`when the readFile("%s") is called`, tc.uri), func() {
								result, err := vfs.ReadFile(tc.uri)

								Convey("then the result should be as expected", func() {
									if tc.wantError != "" {
										So(err, ShouldNotBeNil)
										So(err.Error(), ShouldEqual, fmt.Sprintf("readfile %s", tc.wantError))
									} else {
										So(err, ShouldBeNil)
										So(result, ShouldResemble, tc.wantResult)
									}
								})
							})
						})
					})
				})
			})
		}
	})
}
