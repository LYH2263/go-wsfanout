package wsfanout

import (
	"example.com/wsfanout/internal/audit"
	"example.com/wsfanout/internal/clock"
	"example.com/wsfanout/internal/encode"
	"example.com/wsfanout/internal/persist"
)

// Option 配置 Hub。
type Option func(*Hub)

// WithEncoder 设置载荷编码器；未设置时使用 JSON 缺省编码器。
func WithEncoder(enc encode.Encoder) Option {
	return func(h *Hub) { h.encoder = enc }
}

// WithStore 设置房间配置持久化。
func WithStore(s persist.Store) Option {
	return func(h *Hub) { h.store = s }
}

// WithAuditor 设置审计日志。
func WithAuditor(a *audit.Logger) Option {
	return func(h *Hub) { h.auditor = a }
}

// WithClock 注入时钟（测试用）。
func WithClock(c clock.Clock) Option {
	return func(h *Hub) { h.clock = c }
}

// WithDefaultQueue 设置缺省发送队列长度。
func WithDefaultQueue(n int) Option {
	return func(h *Hub) {
		if n > 0 {
			h.defaultQueue = n
		}
	}
}

// WithDefaultPolicy 设置新建房间缺省策略。
func WithDefaultPolicy(p RoomPolicy) Option {
	return func(h *Hub) { h.defaultPolicy = p }
}
