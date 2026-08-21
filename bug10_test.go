package wsfanout_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"example.com/wsfanout"
)

func TestBug10_CloseFlushesBeforeDropRooms(t *testing.T) {
	h := wsfanout.New(wsfanout.WithDefaultQueue(16))
	id := h.MustJoin("lobby", "c1")

	var mu sync.Mutex
	var frames [][]byte
	release := make(chan struct{})
	entered := make(chan struct{}, 1)
	if err := h.AttachWriter(id, func(ctx context.Context, data []byte) error {
		select {
		case entered <- struct{}{}:
		default:
		}
		<-release
		mu.Lock()
		frames = append(frames, append([]byte(nil), data...))
		mu.Unlock()
		return nil
	}); err != nil {
		t.Fatal(err)
	}

	// 第一帧进入写泵后阻塞；第二帧留在队列。Close 必须先 Flush 再拆连接。
	if err := h.Broadcast("lobby", wsfanout.Message{Type: "chat", Body: []byte("frame-1")}); err != nil {
		t.Fatal(err)
	}
	select {
	case <-entered:
	case <-time.After(2 * time.Second):
		t.Fatal("writer did not see first frame")
	}
	if err := h.Broadcast("lobby", wsfanout.Message{Type: "chat", Body: []byte("frame-2")}); err != nil {
		t.Fatal(err)
	}
	time.Sleep(40 * time.Millisecond)

	closed := make(chan error, 1)
	go func() { closed <- h.Close() }()

	select {
	case <-closed:
		t.Fatal("Close returned before flush finished (likely dropped queued frames)")
	case <-time.After(80 * time.Millisecond):
	}
	close(release)
	select {
	case err := <-closed:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Close hung after release")
	}
	mu.Lock()
	n := len(frames)
	mu.Unlock()
	if n < 2 {
		t.Fatalf("Close dropped queued frames before flush; got %d frames", n)
	}
}
