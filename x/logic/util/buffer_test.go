package util

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBoundedBuffer(t *testing.T) {
	Convey("Given test cases", t, func() {
		tcs := []struct {
			size       int
			bytes      [][]byte
			wantResult string
		}{
			{20, [][]byte{[]byte("hello world")}, "hello world"},
			{11, [][]byte{[]byte("hello"), []byte(" "), []byte("world")}, "hello world"},
			{10, [][]byte{[]byte("hello"), []byte(" "), []byte("world")}, "ello world"},
			{5, [][]byte{[]byte("hello world")}, "world"},
			{1, [][]byte{[]byte("hello"), []byte(" "), []byte("world")}, "d"},
			{1, [][]byte{}, ""},
			{0, [][]byte{[]byte("hello world")}, ""},
		}

		for tn, tc := range tcs {
			Convey(fmt.Sprintf("Given a BoundedBuffer with size %d (%d)", tc.size, tn), func() {
				buffer, err := NewBoundedBuffer(tc.size)
				So(err, ShouldBeNil)

				Convey("When calling the Write() function", func() {
					wantCount := 0
					count := 0
					for _, b := range tc.bytes {
						n, err := buffer.Write(b)
						count += n
						wantCount += len(b)
						So(err, ShouldBeNil)
					}

					Convey("Then we should get %s", func() {
						So(buffer.String(), ShouldEqual, tc.wantResult)
					})
					Convey(fmt.Sprintf("And the number of bytes written should be %d", wantCount), func() {
						So(count, ShouldEqual, wantCount)
					})
				})
			})
		}
	})
}
