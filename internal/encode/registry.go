package encode

import "sync"

// Registry 管理命名编码器。
type Registry struct {
	mu  sync.RWMutex
	by  map[string]Encoder
	def string
}

// NewRegistry 创建注册表。
func NewRegistry() *Registry {
	r := &Registry{by: make(map[string]Encoder), def: "json"}
	r.by["json"] = NewJSONEncoder()
	r.by["raw"] = &RawEncoder{}
	return r
}

func (r *Registry) Register(name string, enc Encoder) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.by[name] = enc
}

func (r *Registry) Get(name string) (Encoder, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	e, ok := r.by[name]
	return e, ok
}

func (r *Registry) Default() Encoder {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.by[r.def]
}

func (r *Registry) SetDefault(name string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.by[name]; !ok {
		return false
	}
	r.def = name
	return true
}
