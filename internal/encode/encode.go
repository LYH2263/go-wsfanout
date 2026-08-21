package encode

import (
	"encoding/json"

	"example.com/wsfanout/internal/payload"
)

// Encoder 将内部消息编码为出站帧。
type Encoder interface {
	Encode(m payload.Message) ([]byte, error)
	Name() string
}

// JSONEncoder 缺省 JSON 编码器。
type JSONEncoder struct{}

// NewJSONEncoder 创建缺省编码器。
func NewJSONEncoder() *JSONEncoder { return &JSONEncoder{} }

func (e *JSONEncoder) Name() string { return "json" }

func (e *JSONEncoder) Encode(m payload.Message) ([]byte, error) {
	type wire struct {
		Type      string            `json:"type"`
		Body      string            `json:"body"`
		Headers   map[string]string `json:"headers,omitempty"`
		Timestamp int64             `json:"ts"`
	}
	w := wire{
		Type:      m.Type,
		Body:      string(m.Body),
		Headers:   m.Headers,
		Timestamp: m.Timestamp.UnixMilli(),
	}
	return json.Marshal(w)
}

// RawEncoder 直接输出 Body。
type RawEncoder struct{}

func (e *RawEncoder) Name() string { return "raw" }

func (e *RawEncoder) Encode(m payload.Message) ([]byte, error) {

	return m.Body, nil
}
