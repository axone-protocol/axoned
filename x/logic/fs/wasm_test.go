//nolint:lll
package fs

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"testing"

	"github.com/golang/mock/gomock"

	. "github.com/smartystreets/goconvey/convey"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/okp4/okp4d/x/logic/testutil"
)

func TestWasmHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Given a test cases", t, func() {
		cases := []struct {
			contractAddress string
			query           []byte
			data            []byte
			uri             string
			wantResult      []byte
			wantError       error
		}{
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `cosmwasm:cw-storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("calc(X) :-  X is 100 + 200."),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu"),
				uri:             `cosmwasm:cw-storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Errorf("failed unmarshal json wasm response to string: invalid character 'Y' looking for beginning of value"),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `cosmwasm:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("calc(X) :-  X is 100 + 200."),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `okp4:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantError:       fmt.Errorf("invalid scheme"),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:cw-storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Errorf("failed decode wasm base64 respone: illegal base64 data at input byte 0"),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:cw-storage?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Errorf("failed convert path 'cw-storage' to contract address: decoding bech32 failed: invalid separator index -1"),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:cw-storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?wasm=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Errorf("uri should contains `query` params"),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Errorf("emtpy path given, should be 'cosmwasm:{contractName}:{contractAddr}?query={query}'"),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("foo-bar"),
				uri:             `cosmwasm:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=false`,
				wantResult:      []byte("foo-bar"),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"Y2FsYyhYKSA6LSAgWCBpcyAxMDAgKyAyMDAu\""),
				uri:             `cosmwasm:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=true`,
				wantResult:      []byte("calc(X) :-  X is 100 + 200."),
			},
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           []byte("{\"object_data\":{\"id\": \"4cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05\"}}"),
				data:            []byte("\"hey\""),
				uri:             `cosmwasm:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D&base64Decode=foo`,
				wantResult:      []byte("\"\""),
				wantError:       fmt.Errorf("failed convert 'base64Decode' query value to boolean: strconv.ParseBool: parsing \"foo\": invalid syntax"),
			},
		}
		for nc, tc := range cases {
			Convey(fmt.Sprintf("Given the uri #%d: %s", nc, tc.uri), func() {
				Convey("and a wasm keeper initialized with the given values", func() {
					db := tmdb.NewMemDB()
					stateStore := store.NewCommitMultiStore(db)
					wasmKeeper := testutil.NewMockWasmKeeper(ctrl)
					ctx := sdk.
						NewContext(stateStore, tmproto.Header{}, false, log.NewNopLogger())
					sdk.GetConfig().SetBech32PrefixForAccount("okp4", "okp4pub")

					wasmKeeper.EXPECT().
						QuerySmart(ctx, sdk.MustAccAddressFromBech32(tc.contractAddress), tc.query).
						AnyTimes().
						Return(tc.data, nil)

					Convey("and wasm handler", func() {
						handler := NewWasmHandler(wasmKeeper)

						Convey("When ask handler if it can open uri", func() {
							uri, err := url.Parse(tc.uri)

							So(err, ShouldBeNil)

							Convey("Then handler response should be as expected", func() {
								file, err := handler.Open(ctx, uri)

								if tc.wantError != nil {
									So(err, ShouldNotBeNil)
									So(err.Error(), ShouldEqual, tc.wantError.Error())
								} else {
									So(err, ShouldBeNil)

									defer file.Close()
									info, _ := file.Stat()
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
						})
					})
				})
			})
		}
	})
}
