package wsfanout

import "time"

// Message 是对外广播/点对点推送的文本载荷封装。
type Message struct {
	Type      string
	Body      []byte
	Headers   map[string]string
	Timestamp time.Time
}

// RoomInfo 是 Snapshot 导出的房间视图。
type RoomInfo struct {
	Name      string
	Members   []string
	CreatedAt time.Time
	Policy    RoomPolicy
}

// RoomPolicy 控制背压与慢客户端踢出。
type RoomPolicy struct {
	MaxMembers   int
	QueueSize    int
	DropOldest   bool
	KickSlowMs   int
	AllowAnon    bool
	RequireToken bool
}

// HubStats 汇总运行指标。
type HubStats struct {
	Rooms       int
	Connections int
	Broadcasts  int64
	Dropped     int64
	Kicked      int64
	WriteErrors int64
}

// ConnInfo 描述管理页展示的连接。
type ConnInfo struct {
	ID       string
	Room     string
	JoinedAt time.Time
	Queued   int
	Slow     bool
}
