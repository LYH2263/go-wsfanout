package wsfanout_test

import (
	"errors"
	"testing"

	"example.com/wsfanout"
)

func TestBug03_BroadcastAfterCloseNoPanic(t *testing.T) {
	h := wsfanout.New()
	h.MustJoin("lobby", "c1")
	if err := h.Close(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Broadcast after Close panicked: %v", r)
		}
	}()
	err := h.Broadcast("lobby", wsfanout.Message{Type: "x", Body: []byte("y")})
	if !errors.Is(err, wsfanout.ErrClosed) {
		t.Fatalf("want ErrClosed, got %v", err)
	}
}
