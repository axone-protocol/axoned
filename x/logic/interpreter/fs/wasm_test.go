//nolint:gocognit,lll
package fs

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/mock/gomock"
	"github.com/okp4/okp4d/x/logic/testutil"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmdb "github.com/tendermint/tm-db"
)

func TestWasmHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	Convey("Given a test cases", t, func() {
		cases := []struct {
			contractAddress string
			query           string
			data            []byte
			canOpen         bool
			uri             string
			wantResult      []byte
			wantError       error
		}{
			{
				contractAddress: "okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht",
				query:           "",
				data:            []byte(""),
				canOpen:         true,
				uri:             `cosmwasm:cw-storage:okp415ekvz3qdter33mdnk98v8whv5qdr53yusksnfgc08xd26fpdn3ts8gddht?query=%7B%22object_data%22%3A%7B%22id%22%3A%20%224cbe36399aabfcc7158ee7a66cbfffa525bb0ceab33d1ff2cff08759fe0a9b05%22%7D%7D`,
				wantResult:      []byte("\"\""),
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

					wasmKeeper.EXPECT().
						QuerySmart(ctx, tc.contractAddress, tc.query).
						AnyTimes().
						Return(tc.data, nil)

					Convey("and wasm handler", func() {
						handler := NewWasmFS(wasmKeeper)

						Convey("When ask handler if it can open uri", func() {
							uri, err := url.Parse(tc.uri)

							So(err, ShouldBeNil)

							result := handler.CanOpen(ctx, uri)

							So(result, ShouldEqual, tc.canOpen)

							if result {
								Convey("Then handler response should be as expected", func() {
									data, err := handler.Open(ctx, uri)

									if tc.wantError != nil {
										So(err, ShouldNotBeNil)
										So(err, ShouldEqual, tc.wantError.Error())
									} else {
										So(err, ShouldBeNil)
										So(data, ShouldEqual, tc.wantResult)
									}
								})
							}
						})

					})
				})
			})
		}
	})
}
