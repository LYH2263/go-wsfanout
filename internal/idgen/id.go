package idgen

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync/atomic"
	"time"
)

var seq atomic.Uint64

// NewID 生成带前缀的唯一 ID。
func NewID(prefix string) string {
	n := seq.Add(1)
	var b [4]byte
	_, _ = rand.Read(b[:])
	return fmt.Sprintf("%s-%d-%d-%s", prefix, time.Now().UnixNano(), n, hex.EncodeToString(b[:]))
}

// Short 短 ID。
func Short() string {
	var b [6]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
