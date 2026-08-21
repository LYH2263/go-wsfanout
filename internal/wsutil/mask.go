package wsutil

// Mask 应用 RFC6455 掩码（客户端→服务端）。
func Mask(key [4]byte, data []byte) []byte {
	out := make([]byte, len(data))
	for i := range data {
		out[i] = data[i] ^ key[i%4]
	}
	return out
}

// Unmask 与 Mask 相同。
func Unmask(key [4]byte, data []byte) []byte {
	return Mask(key, data)
}
