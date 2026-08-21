package wsfanout_test

import (
	"errors"
	"testing"

	"example.com/wsfanout"
)

func TestBug04_NilDefaultEncoderNoPanic(t *testing.T) {
	// New 不传 WithEncoder 时必须有缺省编码器，不得空指针
	h := wsfanout.New()
	defer h.Close()
	h.MustJoin("lobby", "c1")
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Broadcast with default encoder panicked: %v", r)
		}
	}()
	err := h.Broadcast("lobby", wsfanout.Message{Type: "chat", Body: []byte("hi")})
	if err != nil && !errors.Is(err, wsfanout.ErrNoEncoder) {
		// 允许返回 ErrNoEncoder，但绝不能 panic；缺省路径应成功
		t.Fatalf("unexpected err: %v", err)
	}
	if err != nil {
		t.Fatalf("default encoder missing: %v", err)
	}
	if h.Encoder() == nil {
		t.Fatal("Encoder() is nil after New()")
	}
}
