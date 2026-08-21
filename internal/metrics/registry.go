package metrics

import "sync/atomic"

// Registry 简单计数器集合。
type Registry struct {
	rooms      atomic.Int64
	conns      atomic.Int64
	broadcasts atomic.Int64
	dropped    atomic.Int64
	kicked     atomic.Int64
	writes     atomic.Int64
	errors     atomic.Int64
}

func NewRegistry() *Registry { return &Registry{} }

func (r *Registry) IncRooms(d int64)      { r.rooms.Add(d) }
func (r *Registry) IncConns(d int64)      { r.conns.Add(d) }
func (r *Registry) IncBroadcasts(d int64) { r.broadcasts.Add(d) }
func (r *Registry) IncDropped(d int64)    { r.dropped.Add(d) }
func (r *Registry) IncKicked(d int64)     { r.kicked.Add(d) }
func (r *Registry) IncWrites(d int64)     { r.writes.Add(d) }
func (r *Registry) IncErrors(d int64)     { r.errors.Add(d) }

// Snapshot 导出。
type Snapshot struct {
	Rooms      int64 `json:"rooms"`
	Conns      int64 `json:"conns"`
	Broadcasts int64 `json:"broadcasts"`
	Dropped    int64 `json:"dropped"`
	Kicked     int64 `json:"kicked"`
	Writes     int64 `json:"writes"`
	Errors     int64 `json:"errors"`
}

func (r *Registry) Snapshot() Snapshot {
	return Snapshot{
		Rooms:      r.rooms.Load(),
		Conns:      r.conns.Load(),
		Broadcasts: r.broadcasts.Load(),
		Dropped:    r.dropped.Load(),
		Kicked:     r.kicked.Load(),
		Writes:     r.writes.Load(),
		Errors:     r.errors.Load(),
	}
}
