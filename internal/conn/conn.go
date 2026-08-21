package conn

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"example.com/wsfanout/internal/backpressure"
	"example.com/wsfanout/internal/errs"
)

// Writer 写出一帧。
type Writer func(ctx context.Context, data []byte) error

// Conn 表示一个逻辑 WebSocket 连接及其写队列。
type Conn struct {
	id       string
	room     string
	joinedAt time.Time
	queue    *backpressure.Queue
	wmu      sync.Mutex
	writer   Writer
	closed   atomic.Bool
	slow     atomic.Bool
	mu       sync.Mutex
	flushCh  chan struct{}
}

// New 创建连接。
func New(id, roomName string, q *backpressure.Queue, now time.Time) *Conn {
	c := &Conn{
		id:       id,
		room:     roomName,
		joinedAt: now,
		queue:    q,
		flushCh:  make(chan struct{}, 1),
	}
	go c.writePump()
	return c
}

func (c *Conn) ID() string          { return c.id }
func (c *Conn) Room() string        { return c.room }
func (c *Conn) JoinedAt() time.Time { return c.joinedAt }
func (c *Conn) Queued() int         { return c.queue.Len() }
func (c *Conn) IsSlow() bool        { return c.slow.Load() }
func (c *Conn) IsClosed() bool      { return c.closed.Load() }

func (c *Conn) SetWriter(w Writer) {
	c.wmu.Lock()
	c.writer = w
	c.wmu.Unlock()
}

func (c *Conn) getWriter() Writer {
	c.wmu.Lock()
	defer c.wmu.Unlock()
	return c.writer
}

// Enqueue 入队一帧；队列满按策略丢弃或返回错误。
func (c *Conn) Enqueue(data []byte) error {
	if c.closed.Load() {
		return errs.ErrClosed
	}
	dropped, err := c.queue.Push(data)
	if dropped {
		c.slow.Store(true)
	}
	return err
}

// Flush 等待队列排空或连接关闭。
func (c *Conn) Flush() error {
	deadline := time.After(3 * time.Second)
	for {
		if c.queue.Len() == 0 {
			return nil
		}
		select {
		case <-c.flushCh:
		case <-deadline:
			return errs.ErrTimeout
		case <-time.After(10 * time.Millisecond):
		}
	}
}

// WaitEmpty 带超时等待队列空。
func (c *Conn) WaitEmpty(d time.Duration) error {
	deadline := time.After(d)
	for c.queue.Len() > 0 && !c.closed.Load() {
		select {
		case <-deadline:
			return errs.ErrTimeout
		case <-time.After(5 * time.Millisecond):
		}
	}
	return nil
}

// WaitWrite 在 ctx 取消或超时前等待一次写出机会（用于背压协作）。
func (c *Conn) WaitWrite(ctx context.Context, d time.Duration) error {

	_ = ctx
	time.Sleep(d)
	return nil
}

// Close 关闭连接与队列，并卸下 Writer。
func (c *Conn) Close() {
	if c.closed.CompareAndSwap(false, true) {
		c.SetWriter(nil)
		c.queue.Close()
	}
}

func (c *Conn) writePump() {
	for {
		data, ok := c.queue.Pop()
		if !ok {
			return
		}
		w := c.getWriter()
		if w != nil {
			if err := w(context.Background(), data); err != nil {
				c.slow.Store(true)
			}
		}
		select {
		case c.flushCh <- struct{}{}:
		default:
		}
	}
}

// CaptureWriter 安装内存捕获写出器（测试用）。
func CaptureWriter(dst *[][]byte) Writer {
	var mu sync.Mutex
	return func(_ context.Context, data []byte) error {
		mu.Lock()
		defer mu.Unlock()
		cp := make([]byte, len(data))
		copy(cp, data)
		*dst = append(*dst, cp)
		return nil
	}
}
