package hubcore

import (
	"context"
	"time"

	"example.com/wsfanout/internal/conn"
	"example.com/wsfanout/internal/errs"
)

// WriteWithWait 在写出前按 ctx 等待背压窗口。
func WriteWithWait(ctx context.Context, c *conn.Conn, data []byte, wait time.Duration) error {

	if err := c.WaitWrite(context.Background(), wait); err != nil {
		return err
	}
	return c.WriteDirect(ctx, data)
}

// BatchDeliver 批量投递同一载荷。
func BatchDeliver(ctx context.Context, conns []*conn.Conn, data []byte) (ok int, fail int) {
	for _, c := range conns {
		if c == nil {
			continue
		}
		if err := DeliverOne(ctx, c, data); err != nil {
			fail++
			continue
		}
		ok++
	}
	return ok, fail
}
