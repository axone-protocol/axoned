package device

import (
	"errors"
	"io"
	"io/fs"
	"testing"
	"time"
)

func TestHalfDuplexLifecycle(t *testing.T) {
	commitCalls := 0
	file := New(HalfDuplexConfig{
		Path:              "/v1/dev/echo",
		ModTime:           time.Unix(0, 0).UTC(),
		AllowEmptyRequest: false,
		Commit: func(request []byte) ([]byte, error) {
			commitCalls++
			if string(request) != "ping" {
				t.Fatalf("unexpected request: %q", string(request))
			}
			return []byte("pong"), nil
		},
	})

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("ping")); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	buf := make([]byte, 8)
	n, err := file.Read(buf)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if got := string(buf[:n]); got != "pong" {
		t.Fatalf("unexpected response: %q", got)
	}

	if _, err := file.Read(buf); !errors.Is(err, io.EOF) {
		t.Fatalf("expected EOF, got: %v", err)
	}

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("x")); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("expected write permission error after commit, got: %v", err)
	}

	if commitCalls != 1 {
		t.Fatalf("commit should be called once, got: %d", commitCalls)
	}
}

func TestHalfDuplexEmptyRequestValidation(t *testing.T) {
	commitCalled := false
	file := New(HalfDuplexConfig{
		Path:              "/v1/dev/test",
		ModTime:           time.Unix(0, 0).UTC(),
		AllowEmptyRequest: false,
		Commit: func(_ []byte) ([]byte, error) {
			commitCalled = true
			return nil, nil
		},
	})

	_, err := file.Read(make([]byte, 1))
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}
	if commitCalled {
		t.Fatal("commit should not be called when request is missing")
	}
}

func TestHalfDuplexAllowEmptyRequest(t *testing.T) {
	file := New(HalfDuplexConfig{
		Path:              "/v1/dev/test",
		ModTime:           time.Unix(0, 0).UTC(),
		AllowEmptyRequest: true,
		Commit: func(request []byte) ([]byte, error) {
			if len(request) != 0 {
				t.Fatalf("expected empty request, got: %v", request)
			}
			return []byte("ok"), nil
		},
	})

	buf := make([]byte, 2)
	n, err := file.Read(buf)
	if err != nil {
		t.Fatalf("read failed: %v", err)
	}
	if got := string(buf[:n]); got != "ok" {
		t.Fatalf("unexpected response: %q", got)
	}
}

func TestHalfDuplexMaxRequestBytes(t *testing.T) {
	file := New(HalfDuplexConfig{
		Path:              "/v1/dev/test",
		ModTime:           time.Unix(0, 0).UTC(),
		MaxRequestBytes:   4,
		AllowEmptyRequest: false,
		Commit: func(request []byte) ([]byte, error) {
			return request, nil
		},
	})

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("12345")); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("expected fs.ErrPermission, got: %v", err)
	}
}

func TestHalfDuplexClosed(t *testing.T) {
	file := New(HalfDuplexConfig{
		Path:              "/v1/dev/test",
		ModTime:           time.Unix(0, 0).UTC(),
		AllowEmptyRequest: false,
		Commit: func(request []byte) ([]byte, error) {
			return request, nil
		},
	})

	if err := file.Close(); err != nil {
		t.Fatalf("close failed: %v", err)
	}

	if _, err := file.Read(make([]byte, 1)); !errors.Is(err, fs.ErrClosed) {
		t.Fatalf("expected fs.ErrClosed on read, got: %v", err)
	}
	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("x")); !errors.Is(err, fs.ErrClosed) {
		t.Fatalf("expected fs.ErrClosed on write, got: %v", err)
	}
}

func TestHalfDuplexCommitErrorPropagation(t *testing.T) {
	expected := errors.New("boom")
	file := New(HalfDuplexConfig{
		Path:              "/v1/dev/test",
		ModTime:           time.Unix(0, 0).UTC(),
		AllowEmptyRequest: false,
		Commit: func(_ []byte) ([]byte, error) {
			return nil, expected
		},
	})

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("x")); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	if _, err := file.Read(make([]byte, 1)); !errors.Is(err, expected) {
		t.Fatalf("expected commit error, got: %v", err)
	}
}

func TestHalfDuplexMaxResponseBytes(t *testing.T) {
	file := New(HalfDuplexConfig{
		Path:              "/v1/dev/test",
		ModTime:           time.Unix(0, 0).UTC(),
		AllowEmptyRequest: false,
		MaxResponseBytes:  2,
		Commit: func(_ []byte) ([]byte, error) {
			return []byte("too-big"), nil
		},
	})

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("x")); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	if _, err := file.Read(make([]byte, 16)); !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("expected ErrResponseTooLarge, got: %v", err)
	}
}
