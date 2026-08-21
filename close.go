package wsfanout

import (
	"sync"

	"example.com/wsfanout/internal/conn"
)

// Close 关闭 Hub：先 Flush 各连接写队列（写出器仍在），待在途帧全部写完成，
// 再拆房间表并关闭连接，避免队列里/写出中的帧被丢掉。
func (h *Hub) Close() error {
	if !h.closed.CompareAndSwap(false, true) {
		return nil
	}

	// 仅快照连接引用；此时不关闭、不删除，保证写出器与写队列仍可用。
	h.mu.Lock()
	conns := make([]*conn.Conn, 0, len(h.conns))
	for _, c := range h.conns {
		conns = append(conns, c)
	}
	h.mu.Unlock()

	// 先 Flush：并发等待各连接写队列排空。closed 已置位，Broadcast 不会再入队，
	// 故此处只需等已在队列/写出中的帧排空，不会丢帧。
	var wg sync.WaitGroup
	for _, c := range conns {
		wg.Add(1)
		go func(c *conn.Conn) {
			defer wg.Done()
			_ = c.Flush()
		}(c)
	}
	wg.Wait()

	// 再拆房间/关连接。此时在途帧已写完，卸下写出器、关队列不会丢帧。
	h.mu.Lock()
	defer h.mu.Unlock()
	for id, c := range h.conns {
		c.Close()
		delete(h.conns, id)
	}
	for name := range h.rooms {
		delete(h.rooms, name)
	}
	if h.auditor != nil {
		_ = h.auditor.Close()
	}
	return nil
}

// IsClosed 报告 Hub 是否已关闭。
func (h *Hub) IsClosed() bool {
	return h.closed.Load()
}
