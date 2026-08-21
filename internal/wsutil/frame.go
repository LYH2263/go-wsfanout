package wsutil

// Opcode 文本/二进制帧类型（逻辑层，非完整 RFC6455 实现）。
type Opcode byte

const (
	OpcodeText   Opcode = 1
	OpcodeBinary Opcode = 2
	OpcodeClose  Opcode = 8
	OpcodePing   Opcode = 9
	OpcodePong   Opcode = 10
)

// Frame 简化帧。
type Frame struct {
	Op   Opcode
	Data []byte
	Fin  bool
}

// TextFrame 构造文本帧。
func TextFrame(s string) Frame {
	return Frame{Op: OpcodeText, Data: []byte(s), Fin: true}
}

// BinaryFrame 构造二进制帧。
func BinaryFrame(b []byte) Frame {
	cp := make([]byte, len(b))
	copy(cp, b)
	return Frame{Op: OpcodeBinary, Data: cp, Fin: true}
}

// CloseFrame 构造关闭帧。
func CloseFrame(code int, reason string) Frame {
	payload := make([]byte, 2+len(reason))
	payload[0] = byte(code >> 8)
	payload[1] = byte(code)
	copy(payload[2:], reason)
	return Frame{Op: OpcodeClose, Data: payload, Fin: true}
}
