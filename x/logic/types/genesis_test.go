package types_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/axone-protocol/axoned/v14/x/logic/types"
)

func TestGenesisState_Validate(t *testing.T) {
	Convey("Given genesis state validation cases", t, func() {
		for _, tc := range []struct {
			desc     string
			genState *types.GenesisState
			valid    bool
		}{
			{
				desc:     "default is valid",
				genState: types.DefaultGenesis(),
				valid:    true,
			},
			{
				desc:     "valid genesis state",
				genState: &types.GenesisState{},
				valid:    true,
			},
		} {
			Convey(tc.desc, func() {
				err := tc.genState.Validate()
				if tc.valid {
					So(err, ShouldBeNil)
				} else {
					So(err, ShouldNotBeNil)
				}
			})
		}
	})
}
