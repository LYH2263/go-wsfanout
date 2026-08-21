package backpressure

import (
	"sync"

	"example.com/wsfanout/internal/errs"
)

// Queue 有界发送队列。
type Queue struct {
	mu         sync.Mutex
	ch         chan []byte
	dropOldest bool
	closed     bool
	dropped    int
}

// NewQueue 创建队列。
func NewQueue(size int, dropOldest bool) *Queue {
	if size <= 0 {
		size = 16
	}
	return &Queue{
		ch:         make(chan []byte, size),
		dropOldest: dropOldest,
	}
}

// Push 入队；满时按策略丢最旧或返回 ErrQueueFull。
func (q *Queue) Push(data []byte) (dropped bool, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return false, errs.ErrClosed
	}
	select {
	case q.ch <- data:
		return false, nil
	default:
		if !q.dropOldest {
			return false, errs.ErrQueueFull
		}
		select {
		case <-q.ch:
			q.dropped++
		default:
		}
		select {
		case q.ch <- data:
			return true, nil
		default:
			return true, errs.ErrQueueFull
		}
	}
}

// Pop 出队；关闭且空时 ok=false。
func (q *Queue) Pop() ([]byte, bool) {
	data, ok := <-q.ch
	return data, ok
}

// Len 当前长度（近似）。
func (q *Queue) Len() int {
	return len(q.ch)
}

// Dropped 累计丢弃数。
func (q *Queue) Dropped() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.dropped
}

// Close 关闭队列。
func (q *Queue) Close() {
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.closed {
		return
	}
	q.closed = true
	close(q.ch)
}
