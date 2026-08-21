package audit

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"example.com/wsfanout/internal/errs"
)

// Entry 审计条目。
type Entry struct {
	At     time.Time         `json:"at"`
	Action string            `json:"action"`
	Room   string            `json:"room,omitempty"`
	Conn   string            `json:"conn,omitempty"`
	Meta   map[string]string `json:"meta,omitempty"`
}

// Logger 文件审计日志，支持轮转。
type Logger struct {
	mu       sync.Mutex
	path     string
	f        *os.File
	maxBytes int64
	written  int64
}

// Open 打开审计日志。
func Open(path string, maxBytes int64) (*Logger, error) {
	if maxBytes <= 0 {
		maxBytes = 1 << 20
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil && filepath.Dir(path) != "." {
		// dir may be .
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	st, _ := f.Stat()
	var written int64
	if st != nil {
		written = st.Size()
	}
	return &Logger{path: path, f: f, maxBytes: maxBytes, written: written}, nil
}

// Write 写入一条审计。
func (l *Logger) Write(action, room, connID string, meta map[string]string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return errs.ErrAudit
	}
	ent := Entry{At: time.Now().UTC(), Action: action, Room: room, Conn: connID, Meta: meta}
	b, err := json.Marshal(ent)
	if err != nil {
		return err
	}
	b = append(b, '\n')
	n, err := l.f.Write(b)
	l.written += int64(n)
	if err != nil {
		return errs.WrapAudit(err)
	}
	if l.written >= l.maxBytes {
		if err := l.rotateLocked(); err != nil {
			return err
		}
	}
	return nil
}

// Rotate 强制轮转。
func (l *Logger) Rotate() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.rotateLocked()
}

func (l *Logger) rotateLocked() error {
	if l.f != nil {
		_ = l.f.Sync()
		_ = l.f.Close()
		l.f = nil
	}
	bak := l.path + ".1"
	_ = os.Remove(bak)
	if err := os.Rename(l.path, bak); err != nil && !os.IsNotExist(err) {
		return errs.WrapAudit(err)
	}
	f, err := os.OpenFile(l.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return errs.WrapAudit(err)
	}
	l.f = f
	l.written = 0
	return nil
}

// Close 关闭底层文件。
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.f == nil {
		return nil
	}
	err := l.f.Close()
	l.f = nil
	return err
}

// Path 返回日志路径。
func (l *Logger) Path() string { return l.path }
