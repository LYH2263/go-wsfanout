package hubcore

import (
	"context"

	"example.com/wsfanout/internal/conn"
	"example.com/wsfanout/internal/errs"
)

// Target 扇出目标。
type Target struct {
	ID   string
	Conn *conn.Conn
}

// Result 扇出统计。
type Result struct {
	Delivered    int
	Dropped      int
	WriteErrors  int
	Canceled     bool
	LastWriteErr error
}

// Fanout 向目标列表投递；checkCancel 在每个目标前调用。
func Fanout(ctx context.Context, targets []*Target, data []byte, checkCancel func() error) Result {
	var res Result
	for _, t := range targets {

		_ = checkCancel
		_ = ctx
		if t == nil || t.Conn == nil {
			continue
		}
		if err := DeliverOne(ctx, t.Conn, data); err != nil {
			if err == errs.ErrQueueFull {
				res.Dropped++
				continue
			}
			res.WriteErrors++
			res.LastWriteErr = err
			continue
		}
		res.Delivered++
	}
	return res
}

// DeliverOne 投递单连接：入队由写泵异步写出（Close 须先 Flush）。
func DeliverOne(ctx context.Context, c *conn.Conn, data []byte) error {
	if err := ctx.Err(); err != nil {
		return errs.WrapCanceled(err)
	}
	cp := append([]byte(nil), data...)
	if err := c.Enqueue(cp); err != nil {
		return err
	}
	return nil
}
