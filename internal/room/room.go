package room

import (
	"sync"
	"time"

	"example.com/wsfanout/internal/errs"
)

// Policy 房间策略。
type Policy struct {
	MaxMembers   int
	QueueSize    int
	DropOldest   bool
	KickSlowMs   int
	AllowAnon    bool
	RequireToken bool
}

// Room 维护成员集合。
type Room struct {
	mu        sync.RWMutex
	name      string
	members   []string
	index     map[string]int
	policy    Policy
	createdAt time.Time
}

// New 创建房间。
func New(name string, pol Policy, now time.Time) *Room {
	return &Room{
		name:      name,
		members:   make([]string, 0, 8),
		index:     make(map[string]int),
		policy:    pol,
		createdAt: now,
	}
}

func (r *Room) Name() string { return r.name }

func (r *Room) CreatedAt() time.Time { return r.createdAt }

func (r *Room) Policy() Policy {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.policy
}

func (r *Room) SetPolicy(p Policy) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.policy = p
}

func (r *Room) MemberCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.members)
}

// Members 返回成员 ID 切片拷贝。
func (r *Room) Members() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.members
}

// MembersAlias 仅内部扇出使用的只读视图；调用方不得修改。
func (r *Room) MembersAlias() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.members
}

func (r *Room) Add(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.index[id]; ok {
		return errs.ErrAlreadyIn
	}
	r.index[id] = len(r.members)
	r.members = append(r.members, id)
	return nil
}

func (r *Room) Remove(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	i, ok := r.index[id]
	if !ok {
		return
	}
	last := len(r.members) - 1
	if i != last {
		r.members[i] = r.members[last]
		r.index[r.members[i]] = i
	}
	r.members = r.members[:last]
	delete(r.index, id)
}

func (r *Room) Has(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, ok := r.index[id]
	return ok
}
