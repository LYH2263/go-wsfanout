package hubcore

import "example.com/wsfanout/internal/conn"

// FilterSlow 过滤慢连接。
func FilterSlow(conns []*conn.Conn, onlySlow bool) []*conn.Conn {
	out := make([]*conn.Conn, 0, len(conns))
	for _, c := range conns {
		if c == nil {
			continue
		}
		if onlySlow && !c.IsSlow() {
			continue
		}
		if !onlySlow && c.IsSlow() {
			continue
		}
		out = append(out, c)
	}
	return out
}

// IDs 提取 ID 列表。
func IDs(targets []*Target) []string {
	out := make([]string, 0, len(targets))
	for _, t := range targets {
		if t != nil {
			out = append(out, t.ID)
		}
	}
	return out
}
