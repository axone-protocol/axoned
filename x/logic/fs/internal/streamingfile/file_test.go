package streamingfile

import (
	"errors"
	"io"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRead(t *testing.T) {
	Convey("Given a streaming file", t, func() {
		read := func(f *File[int]) error {
			_, err := f.Read(make([]byte, 16))
			return err
		}

		Convey("When EOF is reached and cleanup fails", func() {
			stopErr := errors.New("stop failed")
			f := New[int](
				"test",
				time.Unix(0, 0),
				func() (Next[int], Stop, error) {
					next := func() (int, bool, error) {
						return 0, false, nil
					}
					stop := func() error { return stopErr }
					return next, stop, nil
				},
				func(_ int) ([]byte, error) {
					return []byte("unused"), nil
				},
			)

			err := read(f)

			Convey("Then it should return cleanup error", func() {
				So(errors.Is(err, stopErr), ShouldBeTrue)
			})
		})

		Convey("When next fails and cleanup fails", func() {
			nextErr := errors.New("next failed")
			stopErr := errors.New("stop failed")
			f := New[int](
				"test",
				time.Unix(0, 0),
				func() (Next[int], Stop, error) {
					next := func() (int, bool, error) {
						return 0, false, nextErr
					}
					stop := func() error { return stopErr }
					return next, stop, nil
				},
				func(_ int) ([]byte, error) {
					return []byte("unused"), nil
				},
			)

			err := read(f)

			Convey("Then it should preserve both errors", func() {
				So(errors.Is(err, nextErr), ShouldBeTrue)
				So(errors.Is(err, stopErr), ShouldBeTrue)
			})
		})

		Convey("When render fails and cleanup fails", func() {
			renderErr := errors.New("render failed")
			stopErr := errors.New("stop failed")
			f := New[int](
				"test",
				time.Unix(0, 0),
				func() (Next[int], Stop, error) {
					next := func() (int, bool, error) {
						return 1, true, nil
					}
					stop := func() error { return stopErr }
					return next, stop, nil
				},
				func(_ int) ([]byte, error) {
					return nil, renderErr
				},
			)

			err := read(f)

			Convey("Then it should preserve both errors", func() {
				So(errors.Is(err, renderErr), ShouldBeTrue)
				So(errors.Is(err, stopErr), ShouldBeTrue)
			})
		})

		Convey("When EOF is reached and cleanup succeeds", func() {
			f := New[int](
				"test",
				time.Unix(0, 0),
				func() (Next[int], Stop, error) {
					next := func() (int, bool, error) {
						return 0, false, nil
					}
					stop := func() error { return nil }
					return next, stop, nil
				},
				func(_ int) ([]byte, error) {
					return []byte("unused"), nil
				},
			)

			err := read(f)

			Convey("Then it should return EOF", func() {
				So(errors.Is(err, io.EOF), ShouldBeTrue)
			})
		})
	})
}
