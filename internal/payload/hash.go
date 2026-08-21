package payload

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashBody 计算消息体摘要（审计/去重）。
func HashBody(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

// EqualBody 比较两段 body。
func EqualBody(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
