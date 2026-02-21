package prologterm

import (
	"testing"

	"github.com/axone-protocol/prolog/v3/engine"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRender(t *testing.T) {
	Convey("Given a Prolog term renderer", t, func() {
		Convey("when quoted is enabled", func() {
			content, err := Render(engine.NewAtom("axone-testchain-1"), true)

			So(err, ShouldBeNil)
			So(content, ShouldResemble, []byte("'axone-testchain-1'.\n"))
		})

		Convey("when quoted is disabled", func() {
			content, err := Render(engine.NewAtom("axone-testchain-1"), false)

			So(err, ShouldBeNil)
			So(content, ShouldResemble, []byte("axone-testchain-1.\n"))
		})
	})
}
