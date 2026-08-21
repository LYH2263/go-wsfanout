package wsfanout

import (
	"example.com/wsfanout/internal/errs"
	"example.com/wsfanout/internal/persist"
	"example.com/wsfanout/internal/room"
	"example.com/wsfanout/internal/validate"
)

// ApplyRoomConfig 热更新房间策略；持久化失败时回滚内存变更。
func (h *Hub) ApplyRoomConfig(name string, pol RoomPolicy) error {
	if h.closed.Load() {
		return ErrClosed
	}
	if err := validate.RoomName(name); err != nil {
		return err
	}
	if pol.MaxMembers < 0 || pol.QueueSize < 0 {
		return errs.ErrInvalidRoom
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	if h.closed.Load() {
		return ErrClosed
	}
	r := h.ensureRoomLocked(name)
	prev := r.Policy()
	next := room.Policy{
		MaxMembers:   pol.MaxMembers,
		QueueSize:    pol.QueueSize,
		DropOldest:   pol.DropOldest,
		KickSlowMs:   pol.KickSlowMs,
		AllowAnon:    pol.AllowAnon,
		RequireToken: pol.RequireToken,
	}
	r.SetPolicy(next)

	rec := persist.RoomRecord{
		Name:         name,
		MaxMembers:   next.MaxMembers,
		QueueSize:    next.QueueSize,
		DropOldest:   next.DropOldest,
		KickSlowMs:   next.KickSlowMs,
		AllowAnon:    next.AllowAnon,
		RequireToken: next.RequireToken,
		UpdatedAt:    h.clock.Now(),
	}
	if err := h.store.SaveRoom(rec); err != nil {
		r.SetPolicy(prev)
		return errs.WrapPersist(err)
	}
	if h.auditor != nil {
		_ = h.auditor.Write("apply", name, "", map[string]string{"ok": "1"})
	}
	return nil
}

// LoadRoomConfig 从持久化加载策略（若存在）。
func (h *Hub) LoadRoomConfig(name string) (RoomPolicy, error) {
	rec, err := h.store.LoadRoom(name)
	if err != nil {
		return RoomPolicy{}, err
	}
	return RoomPolicy{
		MaxMembers:   rec.MaxMembers,
		QueueSize:    rec.QueueSize,
		DropOldest:   rec.DropOldest,
		KickSlowMs:   rec.KickSlowMs,
		AllowAnon:    rec.AllowAnon,
		RequireToken: rec.RequireToken,
	}, nil
}
