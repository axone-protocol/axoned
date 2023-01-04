package context

import (
	goctx "context"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestWithLimit(t *testing.T) {
	cases := []struct {
		limit      uint64
		increments []struct {
			delta        uint64
			shouldCancel bool
		}
	}{
		{100, []struct {
			delta        uint64
			shouldCancel bool
		}{{10, false}, {10, false}, {80, false}}},
		{10, []struct {
			delta        uint64
			shouldCancel bool
		}{{5, false}, {5, false}, {1, true}, {1, true}}},
		{0, []struct {
			delta        uint64
			shouldCancel bool
		}{{0, false}, {1, true}, {0, true}}},
	}
	Convey("Given a parent context", t, func() {
		parent := goctx.Background()
		for ntc, tc := range cases {
			Convey(fmt.Sprintf("Given a parent context and a limit of %d (case %d)", tc.limit, ntc), func() {
				limit := tc.limit

				Convey("When WithLimit is called with the parent context and the limit", func() {
					ctx, inc := WithLimit(parent, limit)

					Convey("The returned increment function should correctly update the count value", func() {
						expectedCount := uint64(0)
						for ni, increment := range tc.increments {
							count := inc(increment.delta)
							expectedCount += increment.delta
							So(count, ShouldEqual, expectedCount)

							Convey(fmt.Sprintf("And the returned context's Done channel should be %s (case %d)", func() string {
								if increment.shouldCancel {
									return "closed"
								}

								return "open"
							}(), ni), func() {
								select {
								case <-ctx.Done():
									So(increment.shouldCancel, ShouldBeTrue)
								default:
									So(increment.shouldCancel, ShouldBeFalse)
								}
							})
						}
					})
				})
			})
		}
	})
}
