package wsfanout

import (
	"sync"

	"example.com/wsfanout/internal/conn"
)

// Close 关闭 Hub：先 Flush 各连接写队列，再拆房间表并关闭连接。
func (h *Hub) Close() error {
	if !h.closed.CompareAndSwap(false, true) {
		return nil
	}

	h.mu.Lock()
	conns := make([]*conn.Conn, 0, len(h.conns))
	for _, c := range h.conns {
		conns = append(conns, c)
	}

	for id, c := range h.conns {
		c.Close()
		delete(h.conns, id)
	}
	for name := range h.rooms {
		delete(h.rooms, name)
	}
	h.mu.Unlock()

	var wg sync.WaitGroup
	for _, c := range conns {
		wg.Add(1)
		go func(c *conn.Conn) {
			defer wg.Done()
			_ = c.Flush()
		}(c)
	}
	wg.Wait()

	h.mu.Lock()
	defer h.mu.Unlock()
	if h.auditor != nil {
		_ = h.auditor.Close()
	}
	return nil
}

// IsClosed 报告 Hub 是否已关闭。
func (h *Hub) IsClosed() bool {
	return h.closed.Load()
}
