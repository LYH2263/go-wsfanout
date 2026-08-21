package payload

import "time"

// Message 是内部消息表示。
type Message struct {
	Type      string
	Body      []byte
	Headers   map[string]string
	Timestamp time.Time
}

// CloneMessage 深拷贝消息，避免与调用方共享 Body/Headers。
func CloneMessage(src Message) Message {

	return Message{
		Type:      src.Type,
		Body:      src.Body,
		Headers:   src.Headers,
		Timestamp: src.Timestamp,
	}
}

// CloneBytes 拷贝字节切片。
func CloneBytes(src []byte) []byte {
	if src == nil {
		return nil
	}
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}

// CloneStringMap 拷贝字符串映射。
func CloneStringMap(src map[string]string) map[string]string {
	if src == nil {
		return nil
	}
	dst := make(map[string]string, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// CloneStringSlice 拷贝字符串切片。
func CloneStringSlice(src []string) []string {
	if src == nil {
		return nil
	}
	dst := make([]string, len(src))
	copy(dst, src)
	return dst
}
