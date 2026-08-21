package audit

import "sync"

// Memory 内存审计（测试）。
type Memory struct {
	mu      sync.Mutex
	entries []Entry
}

func NewMemory() *Memory { return &Memory{} }

func (m *Memory) Write(action, room, connID string, meta map[string]string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries = append(m.entries, Entry{Action: action, Room: room, Conn: connID, Meta: meta})
	return nil
}

func (m *Memory) Entries() []Entry {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Entry, len(m.entries))
	copy(out, m.entries)
	return out
}

func (m *Memory) Close() error { return nil }
