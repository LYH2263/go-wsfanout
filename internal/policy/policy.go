package policy

import "example.com/wsfanout/internal/room"

// OnJoin 加入时策略检查钩子。
func OnJoin(pol room.Policy, connID string) error {
	if pol.RequireToken && connID == "" {
		return errToken
	}
	if !pol.AllowAnon && looksAnon(connID) {
		return errAnon
	}
	return nil
}

// ShouldKick 根据慢写标记与策略决定是否踢出。
func ShouldKick(pol room.Policy, slow bool, slowForMs int) bool {
	if !slow {
		return false
	}
	if pol.KickSlowMs <= 0 {
		return false
	}
	return slowForMs >= pol.KickSlowMs
}

// Normalize 规范化策略缺省值。
func Normalize(pol room.Policy) room.Policy {
	if pol.MaxMembers <= 0 {
		pol.MaxMembers = 256
	}
	if pol.QueueSize <= 0 {
		pol.QueueSize = 64
	}
	if pol.KickSlowMs < 0 {
		pol.KickSlowMs = 0
	}
	return pol
}

var errToken = errStr("wsfanout: token required")
var errAnon = errStr("wsfanout: anon not allowed")

type errStr string

func (e errStr) Error() string { return string(e) }

func looksAnon(id string) bool {
	return len(id) == 0 || id == "anon" || id == "guest"
}
