package conn

import (
	"context"

	"example.com/wsfanout/internal/errs"
)

// WriteDirect 直接调用 Writer（绕过队列，测试/紧急路径）。
func (c *Conn) WriteDirect(ctx context.Context, data []byte) error {
	if c.closed.Load() {
		return errs.ErrClosed
	}
	if err := ctx.Err(); err != nil {
		return errs.WrapCanceled(err)
	}
	w := c.getWriter()
	if w == nil {
		// 无写出器时入队
		return c.Enqueue(data)
	}
	if err := w(ctx, data); err != nil {
		return errs.WrapWrite(err)
	}
	return nil
}
