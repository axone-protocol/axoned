package types_test

import (
	"fmt"
	"testing"

	"github.com/okp4/okp4d/x/knowledge/types"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDataspaceKey(t *testing.T) {
	Convey("Considering DataspaceKey() function", t, func(c C) {
		cases := []struct {
			name     string
			id       string
			expected []byte
		}{
			{
				name:     "nominal",
				id:       "f21e8600-5a23-4214-81c8-26fb29ab8cb6",
				expected: append(append([]byte{}, 0x11), "f21e8600-5a23-4214-81c8-26fb29ab8cb6"...),
			},
			{
				name:     "pathological",
				id:       "",
				expected: append([]byte{}, 0x11),
			},
		}
		for n, c := range cases {
			Convey(fmt.Sprintf("When calling function with id %s (case %d)", c.id, n), func() {
				result := types.GetDataspaceKey(c.id)

				Convey(fmt.Sprintf("Then result shall be `%s`", c.expected), func() {
					So(result, ShouldResemble, c.expected)
				})
			})
		}
	})
}
