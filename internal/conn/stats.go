package conn

// Stats 连接侧统计。
type Stats struct {
	ID      string
	Queued  int
	Closed  bool
	Slow    bool
	Dropped int
}

// Stats 导出。
func (c *Conn) Stats() Stats {
	return Stats{
		ID:      c.id,
		Queued:  c.queue.Len(),
		Closed:  c.closed.Load(),
		Slow:    c.slow.Load(),
		Dropped: c.queue.Dropped(),
	}
}
