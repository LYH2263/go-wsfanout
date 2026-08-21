package backpressure

// Decide 根据利用率决定是否标记慢客户端。
func Decide(queued, capacity int, dropOldest bool) (slow bool, shouldKick bool) {
	if capacity <= 0 {
		return false, false
	}
	ratio := float64(queued) / float64(capacity)
	if ratio >= 0.9 {
		return true, !dropOldest
	}
	if ratio >= 0.7 {
		return true, false
	}
	return false, false
}

// Cap 规范化队列容量。
func Cap(n, min, max int) int {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}
