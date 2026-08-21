package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"example.com/wsfanout/internal/errs"
)

// RoomRecord 持久化房间策略。
type RoomRecord struct {
	Name         string    `json:"name"`
	MaxMembers   int       `json:"max_members"`
	QueueSize    int       `json:"queue_size"`
	DropOldest   bool      `json:"drop_oldest"`
	KickSlowMs   int       `json:"kick_slow_ms"`
	AllowAnon    bool      `json:"allow_anon"`
	RequireToken bool      `json:"require_token"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Store 房间配置存储。
type Store interface {
	SaveRoom(rec RoomRecord) error
	LoadRoom(name string) (RoomRecord, error)
	ListRooms() ([]RoomRecord, error)
}

// MemoryStore 内存实现；可注入失败。
type MemoryStore struct {
	mu   sync.Mutex
	data map[string]RoomRecord
	fail bool
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{data: make(map[string]RoomRecord)}
}

func (s *MemoryStore) SetFail(v bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.fail = v
}

func (s *MemoryStore) SaveRoom(rec RoomRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.fail {
		return errs.ErrPersist
	}
	s.data[rec.Name] = rec
	return nil
}

func (s *MemoryStore) LoadRoom(name string) (RoomRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	rec, ok := s.data[name]
	if !ok {
		return RoomRecord{}, errs.ErrNotFound
	}
	return rec, nil
}

func (s *MemoryStore) ListRooms() ([]RoomRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]RoomRecord, 0, len(s.data))
	for _, r := range s.data {
		out = append(out, r)
	}
	return out, nil
}

// FileStore JSON 文件存储。
type FileStore struct {
	path string
	mu   sync.Mutex
}

func NewFileStore(path string) *FileStore {
	return &FileStore{path: path}
}

func (s *FileStore) SaveRoom(rec RoomRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, _ := s.loadLocked()
	all[rec.Name] = rec
	return s.saveLocked(all)
}

func (s *FileStore) LoadRoom(name string) (RoomRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, err := s.loadLocked()
	if err != nil {
		return RoomRecord{}, err
	}
	rec, ok := all[name]
	if !ok {
		return RoomRecord{}, errs.ErrNotFound
	}
	return rec, nil
}

func (s *FileStore) ListRooms() ([]RoomRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	all, err := s.loadLocked()
	if err != nil {
		return nil, err
	}
	out := make([]RoomRecord, 0, len(all))
	for _, r := range all {
		out = append(out, r)
	}
	return out, nil
}

func (s *FileStore) loadLocked() (map[string]RoomRecord, error) {
	b, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]RoomRecord{}, nil
		}
		return nil, err
	}
	var all map[string]RoomRecord
	if err := json.Unmarshal(b, &all); err != nil {
		return nil, err
	}
	if all == nil {
		all = map[string]RoomRecord{}
	}
	return all, nil
}

func (s *FileStore) saveLocked(all map[string]RoomRecord) error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(all, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, b, 0o644)
}
