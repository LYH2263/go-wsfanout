package buffer

import "sync"

// Pool 简单字节缓冲池。
type Pool struct {
	p sync.Pool
}

// NewPool 创建池。
func NewPool(size int) *Pool {
	if size <= 0 {
		size = 512
	}
	return &Pool{p: sync.Pool{New: func() any {
		b := make([]byte, 0, size)
		return &b
	}}}
}

// Get 取缓冲。
func (p *Pool) Get() *[]byte {
	return p.p.Get().(*[]byte)
}

// Put 还缓冲。
func (p *Pool) Put(b *[]byte) {
	if b == nil {
		return
	}
	*b = (*b)[:0]
	p.p.Put(b)
}

// CopyTo 拷贝到新切片。
func CopyTo(dst []byte, src []byte) []byte {
	need := len(src)
	if cap(dst) < need {
		dst = make([]byte, need)
	} else {
		dst = dst[:need]
	}
	copy(dst, src)
	return dst
}
