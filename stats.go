package wsfanout

// Stats 返回运行快照。
func (h *Hub) Stats() HubStats {
	h.mu.RLock()
	rooms := len(h.rooms)
	conns := len(h.conns)
	h.mu.RUnlock()
	return HubStats{
		Rooms:       rooms,
		Connections: conns,
		Broadcasts:  h.broadcasts.Load(),
		Dropped:     h.dropped.Load(),
		Kicked:      h.kicked.Load(),
		WriteErrors: h.writeErrors.Load(),
	}
}
