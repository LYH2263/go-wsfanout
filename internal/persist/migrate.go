package persist

import "time"

// NormalizeRecord 填充缺省字段。
func NormalizeRecord(rec RoomRecord) RoomRecord {
	if rec.MaxMembers <= 0 {
		rec.MaxMembers = 256
	}
	if rec.QueueSize <= 0 {
		rec.QueueSize = 64
	}
	if rec.UpdatedAt.IsZero() {
		rec.UpdatedAt = time.Now().UTC()
	}
	return rec
}

// DiffPolicy 比较两个记录是否策略相同。
func DiffPolicy(a, b RoomRecord) bool {
	return a.MaxMembers == b.MaxMembers &&
		a.QueueSize == b.QueueSize &&
		a.DropOldest == b.DropOldest &&
		a.KickSlowMs == b.KickSlowMs &&
		a.AllowAnon == b.AllowAnon &&
		a.RequireToken == b.RequireToken
}
