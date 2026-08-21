package metrics

import (
	"fmt"
	"strings"
)

// RenderText 渲染简易文本指标。
func (r *Registry) RenderText() string {
	s := r.Snapshot()
	var b strings.Builder
	fmt.Fprintf(&b, "wsfanout_rooms %d\n", s.Rooms)
	fmt.Fprintf(&b, "wsfanout_conns %d\n", s.Conns)
	fmt.Fprintf(&b, "wsfanout_broadcasts_total %d\n", s.Broadcasts)
	fmt.Fprintf(&b, "wsfanout_dropped_total %d\n", s.Dropped)
	fmt.Fprintf(&b, "wsfanout_kicked_total %d\n", s.Kicked)
	fmt.Fprintf(&b, "wsfanout_writes_total %d\n", s.Writes)
	fmt.Fprintf(&b, "wsfanout_errors_total %d\n", s.Errors)
	return b.String()
}
