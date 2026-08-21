package conn

import (
	"context"
	"testing"
	"time"

	"example.com/wsfanout/internal/backpressure"
)

func TestBug08_WriteWaitHonorsContext(t *testing.T) {
	q := backpressure.NewQueue(4, true)
	c := New("c1", "r", q, time.Now())
	defer c.Close()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := time.Now()
	err := c.WaitWrite(ctx, 2*time.Second)
	if err == nil {
		t.Fatal("WaitWrite ignored canceled ctx")
	}
	if time.Since(start) > 500*time.Millisecond {
		t.Fatalf("WaitWrite blocked too long: %v", time.Since(start))
	}
}
