package wsfanout_test

import (
	"context"
	"errors"
	"testing"

	"example.com/wsfanout"
)

var errWriteSentinel = errors.New("wsfanout-test: write sentinel")

func TestBug05_WriteErrorWrapsSentinel(t *testing.T) {
	h := wsfanout.New()
	defer h.Close()
	id := h.MustJoin("lobby", "c1")
	if err := h.AttachWriter(id, func(ctx context.Context, data []byte) error {
		return errWriteSentinel
	}); err != nil {
		t.Fatal(err)
	}
	err := h.SendTo(id, wsfanout.Message{Type: "x", Body: []byte("y")})
	if err == nil {
		t.Fatal("expected write error")
	}
	if !errors.Is(err, wsfanout.ErrWriteFailed) {
		t.Fatalf("want ErrWriteFailed, got %v", err)
	}
	if !errors.Is(err, errWriteSentinel) {
		t.Fatalf("want sentinel via %%w chain, got %v", err)
	}
}
