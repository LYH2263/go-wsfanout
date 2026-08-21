package payload

// MergeHeaders 合并头；后者覆盖前者。
func MergeHeaders(base, over map[string]string) map[string]string {
	out := CloneStringMap(base)
	if out == nil {
		out = make(map[string]string)
	}
	for k, v := range over {
		out[k] = v
	}
	return out
}

// HeaderGet 忽略大小写取值（简化）。
func HeaderGet(h map[string]string, key string) string {
	if h == nil {
		return ""
	}
	if v, ok := h[key]; ok {
		return v
	}
	for k, v := range h {
		if equalFold(k, key) {
			return v
		}
	}
	return ""
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		ca, cb := a[i], b[i]
		if ca >= 'A' && ca <= 'Z' {
			ca += 'a' - 'A'
		}
		if cb >= 'A' && cb <= 'Z' {
			cb += 'a' - 'A'
		}
		if ca != cb {
			return false
		}
	}
	return true
}
