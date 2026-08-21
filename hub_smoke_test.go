package wsfanout_test

import (
	"testing"
	"time"

	"example.com/wsfanout"
	"example.com/wsfanout/internal/conn"
)

func TestHubSmokeJoinBroadcast(t *testing.T) {
	h := wsfanout.New()
	defer h.Close()
	id, err := h.Join("lobby", "c1")
	if err != nil {
		t.Fatal(err)
	}
	var got [][]byte
	if err := h.AttachWriter(id, conn.CaptureWriter(&got)); err != nil {
		t.Fatal(err)
	}
	if err := h.Broadcast("lobby", wsfanout.Message{Type: "chat", Body: []byte("hi")}); err != nil {
		t.Fatal(err)
	}
	deadline := time.After(2 * time.Second)
	for len(got) == 0 {
		select {
		case <-deadline:
			t.Fatal("no frame")
		case <-time.After(10 * time.Millisecond):
		}
	}
}
