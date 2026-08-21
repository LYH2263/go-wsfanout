package wsfanout_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"example.com/wsfanout"
	"example.com/wsfanout/internal/encode"
)

func TestBug01_BroadcastBodySliceAlias(t *testing.T) {
	h := wsfanout.New(wsfanout.WithEncoder(&encode.RawEncoder{}))
	defer h.Close()
	id := h.MustJoin("lobby", "c1")
	var got [][]byte
	if err := h.AttachWriter(id, func(_ context.Context, data []byte) error {
		got = append(got, data) // 保留出站切片本身，便于发现别名
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	body := []byte("hello-buffer")
	if err := h.Broadcast("lobby", wsfanout.Message{Type: "chat", Body: body}); err != nil {
		t.Fatal(err)
	}
	deadline := time.After(2 * time.Second)
	for len(got) == 0 {
		select {
		case <-deadline:
			t.Fatal("no frame delivered")
		case <-time.After(5 * time.Millisecond):
		}
	}
	before := append([]byte(nil), got[0]...)
	body[0] = 'X'
	if !bytes.Equal(got[0], before) {
		t.Fatalf("broadcast body aliased with caller buffer: before=%q after=%q", before, got[0])
	}
}

// TestBug01_BroadcastHeadersMapAlias 覆盖 Headers 映射深拷贝：
// Broadcast 返回后改写调用方 Headers 不得污染在途帧。
func TestBug01_BroadcastHeadersMapAlias(t *testing.T) {
	// JSONEncoder 才会把 Headers 编进出站帧，用于探测映射别名。
	h := wsfanout.New()
	defer h.Close()
	id := h.MustJoin("lobby", "c1")
	var got []byte
	if err := h.AttachWriter(id, func(_ context.Context, data []byte) error {
		cp := make([]byte, len(data))
		copy(cp, data)
		got = cp
		return nil
	}); err != nil {
		t.Fatal(err)
	}
	headers := map[string]string{"tag": "v1"}
	if err := h.Broadcast("lobby", wsfanout.Message{Type: "chat", Body: []byte("hi"), Headers: headers}); err != nil {
		t.Fatal(err)
	}
	deadline := time.After(2 * time.Second)
	for got == nil {
		select {
		case <-deadline:
			t.Fatal("no frame delivered")
		case <-time.After(5 * time.Millisecond):
		}
	}
	before := append([]byte(nil), got...)
	headers["tag"] = "mutated"
	headers["sneaky"] = "extra"
	if !bytes.Equal(got, before) {
		t.Fatalf("broadcast headers aliased with caller map: before=%q after=%q", before, got)
	}
}
