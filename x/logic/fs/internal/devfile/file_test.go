package devfile

import (
	"errors"
	"io"
	"io/fs"
	"testing"
	"time"
)

func TestHalfDuplexLifecycle(t *testing.T) {
	commitCalls := 0
	file, err := New(
		WithPath("/v1/dev/echo"),
		WithModTime(time.Unix(0, 0).UTC()),
		WithAllowEmptyRequest(false),
		WithCommit(func(request []byte) ([]byte, error) {
			commitCalls++
			if string(request) != "ping" {
				t.Fatalf("unexpected request: %q", string(request))
			}
			return []byte("pong"), nil
		}),
	)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

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
	file, err := New(
		WithPath("/v1/dev/test"),
		WithModTime(time.Unix(0, 0).UTC()),
		WithAllowEmptyRequest(false),
		WithCommit(func(_ []byte) ([]byte, error) {
			commitCalled = true
			return nil, nil
		}),
	)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	_, err = file.Read(make([]byte, 1))
	if !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected ErrInvalidRequest, got: %v", err)
	}

	if _, err := file.Read(make([]byte, 1)); !errors.Is(err, ErrInvalidRequest) {
		t.Fatalf("expected same ErrInvalidRequest on subsequent reads, got: %v", err)
	}

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("x")); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("expected fs.ErrPermission after failed commit, got: %v", err)
	}

	if commitCalled {
		t.Fatal("commit should not be called when request is missing")
	}
}

func TestHalfDuplexAllowEmptyRequest(t *testing.T) {
	file, err := New(
		WithPath("/v1/dev/test"),
		WithModTime(time.Unix(0, 0).UTC()),
		WithAllowEmptyRequest(true),
		WithCommit(func(request []byte) ([]byte, error) {
			if len(request) != 0 {
				t.Fatalf("expected empty request, got: %v", request)
			}
			return []byte("ok"), nil
		}),
	)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

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
	file, err := New(
		WithPath("/v1/dev/test"),
		WithModTime(time.Unix(0, 0).UTC()),
		WithMaxRequestBytes(4),
		WithAllowEmptyRequest(false),
		WithCommit(func(request []byte) ([]byte, error) {
			return request, nil
		}),
	)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("12345")); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("expected fs.ErrPermission, got: %v", err)
	}
}

func TestHalfDuplexClosed(t *testing.T) {
	file, err := New(
		WithPath("/v1/dev/test"),
		WithModTime(time.Unix(0, 0).UTC()),
		WithAllowEmptyRequest(false),
		WithCommit(func(request []byte) ([]byte, error) {
			return request, nil
		}),
	)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

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
	commitCalls := 0
	file, err := New(
		WithPath("/v1/dev/test"),
		WithModTime(time.Unix(0, 0).UTC()),
		WithAllowEmptyRequest(false),
		WithCommit(func(_ []byte) ([]byte, error) {
			commitCalls++
			return nil, expected
		}),
	)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("x")); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	if _, err := file.Read(make([]byte, 1)); !errors.Is(err, expected) {
		t.Fatalf("expected commit error, got: %v", err)
	}

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("y")); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("expected fs.ErrPermission after failed commit, got: %v", err)
	}

	if _, err := file.Read(make([]byte, 1)); !errors.Is(err, expected) {
		t.Fatalf("expected same commit error on subsequent reads, got: %v", err)
	}

	if commitCalls != 1 {
		t.Fatalf("commit should be called once, got: %d", commitCalls)
	}
}

func TestHalfDuplexMaxResponseBytes(t *testing.T) {
	commitCalls := 0
	file, err := New(
		WithPath("/v1/dev/test"),
		WithModTime(time.Unix(0, 0).UTC()),
		WithAllowEmptyRequest(false),
		WithMaxResponseBytes(2),
		WithCommit(func(_ []byte) ([]byte, error) {
			commitCalls++
			return []byte("too-big"), nil
		}),
	)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("x")); err != nil {
		t.Fatalf("write failed: %v", err)
	}

	if _, err := file.Read(make([]byte, 16)); !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("expected ErrResponseTooLarge, got: %v", err)
	}

	if _, err := file.(interface{ Write([]byte) (int, error) }).Write([]byte("y")); !errors.Is(err, fs.ErrPermission) {
		t.Fatalf("expected fs.ErrPermission after oversized response, got: %v", err)
	}

	if _, err := file.Read(make([]byte, 16)); !errors.Is(err, ErrResponseTooLarge) {
		t.Fatalf("expected same response-too-large error on subsequent reads, got: %v", err)
	}

	if commitCalls != 1 {
		t.Fatalf("commit should be called once, got: %d", commitCalls)
	}
}

func TestHalfDuplexMissingCommit(t *testing.T) {
	_, err := New(
		WithPath("/v1/dev/test"),
		WithModTime(time.Unix(0, 0).UTC()),
	)
	if !errors.Is(err, ErrMissingCommit) {
		t.Fatalf("expected ErrMissingCommit, got: %v", err)
	}
}
