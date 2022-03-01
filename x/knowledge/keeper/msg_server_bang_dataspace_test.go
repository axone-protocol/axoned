package keeper_test

import (
	"fmt"

	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/okp4/okp4d/testutil/keeper"
	"github.com/okp4/okp4d/x/knowledge/keeper"
	"github.com/okp4/okp4d/x/knowledge/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBangDataspace(t *testing.T) {
	Convey("Given a knowledge keeper", t, func(c C) {
		k, ctx := keepertest.KnowledgeKeeper(t)
		srv, goCtx := keeper.NewMsgServerImpl(*k), sdk.WrapSDKContext(ctx)

		Convey("Given a Bangdataspace message", func() {
			msg := &types.MsgBangDataspace{
				Creator:     "x",
				Id:          "8134bf5b-d81b-4260-9591-767fda243828",
				Name:        "My dataspace",
				Description: "A dataspace for testing...",
			}

			Convey("When calling BangDataspace() function", func() {
				response, err := srv.BangDataspace(goCtx, msg)

				Convey("Then no error shall occur and a response shall be provided", func() {
					So(err, ShouldBeNil)
					So(response, ShouldNotBeNil)
					So(response, ShouldResemble, &types.MsgBangDataspaceResponse{})
				})
			})
		})

		Convey("Given a dataspace already created message", func() {
			id := "8b0d4fcd-4f83-4b00-b941-62ed4d454815"
			msg := &types.MsgBangDataspace{
				Creator:     "x",
				Id:          id,
				Name:        "My dataspace",
				Description: "A dataspace for testing...",
			}

			response, err := srv.BangDataspace(goCtx, msg)

			So(err, ShouldBeNil)

			Convey("When calling BangDataspace() function", func(c C) {
				response, err = srv.BangDataspace(goCtx, msg)

				Convey("Then an error shall occur and no response shall be provided", func() {
					So(err, ShouldNotBeNil)
					So(err.Error(), ShouldEqual, fmt.Sprintf("dataspace %s: entity already exists", id))
					So(response, ShouldBeNil)
				})
			})
		})
	})
}
