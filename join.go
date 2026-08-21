package wsfanout

import (
	"example.com/wsfanout/internal/backpressure"
	"example.com/wsfanout/internal/conn"
	"example.com/wsfanout/internal/errs"
	"example.com/wsfanout/internal/idgen"
	"example.com/wsfanout/internal/policy"
	"example.com/wsfanout/internal/validate"
)

// Join 将连接加入房间；connID 为空则自动生成。
func (h *Hub) Join(roomName, connID string) (string, error) {
	if h.closed.Load() {
		return "", ErrClosed
	}
	if err := validate.RoomName(roomName); err != nil {
		return "", err
	}
	if connID == "" {
		connID = idgen.NewID("c")
	}
	if err := validate.ConnID(connID); err != nil {
		return "", err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed.Load() {
		return "", ErrClosed
	}
	if _, exists := h.conns[connID]; exists {
		return "", ErrAlreadyIn
	}
	r := h.ensureRoomLocked(roomName)
	if r.MemberCount() >= r.Policy().MaxMembers && r.Policy().MaxMembers > 0 {
		return "", ErrRoomFull
	}
	qsize := r.Policy().QueueSize
	if qsize <= 0 {
		qsize = h.defaultQueue
	}
	q := backpressure.NewQueue(qsize, r.Policy().DropOldest)
	c := conn.New(connID, roomName, q, h.clock.Now())
	if err := r.Add(connID); err != nil {
		return "", err
	}
	h.conns[connID] = c
	h.metrics.IncConns(1)
	if h.auditor != nil {
		_ = h.auditor.Write("join", roomName, connID, nil)
	}
	_ = policy.OnJoin(r.Policy(), connID)
	return connID, nil
}

// Leave 将连接移出房间并关闭其写队列。
func (h *Hub) Leave(connID string) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	c, ok := h.conns[connID]
	if !ok {
		return ErrNotFound
	}
	roomName := c.Room()
	if r, ok := h.rooms[roomName]; ok {
		r.Remove(connID)
		if r.MemberCount() == 0 {
			delete(h.rooms, roomName)
			h.metrics.IncRooms(-1)
		}
	}
	delete(h.conns, connID)
	h.metrics.IncConns(-1)
	c.Close()
	if h.auditor != nil {
		_ = h.auditor.Write("leave", roomName, connID, nil)
	}
	return nil
}

// MustJoin 同 Join，失败时 panic（仅测试夹具）。
func (h *Hub) MustJoin(roomName, connID string) string {
	id, err := h.Join(roomName, connID)
	if err != nil {
		panic(err)
	}
	return id
}

// AttachWriter 为连接安装写出回调（模拟 WebSocket 写）。
func (h *Hub) AttachWriter(connID string, w conn.Writer) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.conns[connID]
	if !ok {
		return errs.ErrNotFound
	}
	c.SetWriter(w)
	return nil
}
