package wsfanout

import (
	"sync"
	"sync/atomic"

	"example.com/wsfanout/internal/audit"
	"example.com/wsfanout/internal/clock"
	"example.com/wsfanout/internal/conn"
	"example.com/wsfanout/internal/encode"
	"example.com/wsfanout/internal/metrics"
	"example.com/wsfanout/internal/persist"
	"example.com/wsfanout/internal/room"
)

// Hub 是房间扇出中心。
type Hub struct {
	mu            sync.RWMutex
	rooms         map[string]*room.Room
	conns         map[string]*conn.Conn
	closed        atomic.Bool
	encoder       encode.Encoder
	store         persist.Store
	auditor       *audit.Logger
	clock         clock.Clock
	metrics       *metrics.Registry
	defaultQueue  int
	defaultPolicy RoomPolicy
	broadcasts    atomic.Int64
	dropped       atomic.Int64
	kicked        atomic.Int64
	writeErrors   atomic.Int64
}

// New 创建 Hub；未指定 Encoder 时安装缺省 JSON 编码器。
func New(opts ...Option) *Hub {
	h := &Hub{
		rooms:        make(map[string]*room.Room),
		conns:        make(map[string]*conn.Conn),
		metrics:      metrics.NewRegistry(),
		clock:        clock.Real{},
		defaultQueue: 64,
		defaultPolicy: RoomPolicy{
			MaxMembers: 256,
			QueueSize:  64,
			DropOldest: true,
			KickSlowMs: 5000,
			AllowAnon:  true,
		},
	}
	for _, o := range opts {
		o(h)
	}

	if h.store == nil {
		h.store = persist.NewMemoryStore()
	}
	return h
}

func (h *Hub) ensureRoomLocked(name string) *room.Room {
	r, ok := h.rooms[name]
	if ok {
		return r
	}
	pol := room.Policy{
		MaxMembers:   h.defaultPolicy.MaxMembers,
		QueueSize:    h.defaultPolicy.QueueSize,
		DropOldest:   h.defaultPolicy.DropOldest,
		KickSlowMs:   h.defaultPolicy.KickSlowMs,
		AllowAnon:    h.defaultPolicy.AllowAnon,
		RequireToken: h.defaultPolicy.RequireToken,
	}
	r = room.New(name, pol, h.clock.Now())
	h.rooms[name] = r
	h.metrics.IncRooms(1)
	return r
}

// Encoder 返回当前编码器（测试/管理页用）。
func (h *Hub) Encoder() encode.Encoder {
	return h.encoder
}

// Metrics 返回指标注册表。
func (h *Hub) Metrics() *metrics.Registry {
	return h.metrics
}
