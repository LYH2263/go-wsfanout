package wsfanout_test

import (
	"context"
	"testing"

	"example.com/wsfanout"
)

func TestBug07_BroadcastContextHonorsCancel(t *testing.T) {
	h := wsfanout.New()
	defer h.Close()
	h.MustJoin("lobby", "c1")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := h.BroadcastContext(ctx, "lobby", wsfanout.Message{Type: "x", Body: []byte("y")})
	if err == nil {
		t.Fatal("BroadcastContext ignored canceled context")
	}
}
