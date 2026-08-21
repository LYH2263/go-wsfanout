package wsfanout

import (
	"example.com/wsfanout/internal/payload"
	"example.com/wsfanout/internal/room"
)

// Snapshot 导出全部房间视图；成员切片为独立拷贝。
func (h *Hub) Snapshot() []RoomInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]RoomInfo, 0, len(h.rooms))
	for _, r := range h.rooms {
		out = append(out, roomToInfo(r))
	}
	return out
}

// SnapshotRoom 导出单个房间；不存在返回 nil, ErrNotFound。
func (h *Hub) SnapshotRoom(name string) (*RoomInfo, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	r, ok := h.rooms[name]
	if !ok {
		return nil, ErrNotFound
	}
	info := roomToInfo(r)
	return &info, nil
}

// ListConns 列出连接信息。
func (h *Hub) ListConns() []ConnInfo {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]ConnInfo, 0, len(h.conns))
	for _, c := range h.conns {
		out = append(out, ConnInfo{
			ID:       c.ID(),
			Room:     c.Room(),
			JoinedAt: c.JoinedAt(),
			Queued:   c.Queued(),
			Slow:     c.IsSlow(),
		})
	}
	return out
}

func roomToInfo(r *room.Room) RoomInfo {
	members := payload.CloneStringSlice(r.Members())
	pol := r.Policy()
	return RoomInfo{
		Name:      r.Name(),
		Members:   members,
		CreatedAt: r.CreatedAt(),
		Policy: RoomPolicy{
			MaxMembers:   pol.MaxMembers,
			QueueSize:    pol.QueueSize,
			DropOldest:   pol.DropOldest,
			KickSlowMs:   pol.KickSlowMs,
			AllowAnon:    pol.AllowAnon,
			RequireToken: pol.RequireToken,
		},
	}
}
