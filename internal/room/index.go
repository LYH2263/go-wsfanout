package room

// IndexOf 线性查找成员下标；不存在返回 -1。
func (r *Room) IndexOf(id string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if i, ok := r.index[id]; ok {
		return i
	}
	return -1
}

// ClonePolicy 拷贝策略。
func ClonePolicy(p Policy) Policy {
	return p
}
