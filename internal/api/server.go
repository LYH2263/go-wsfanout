package api

import (
	"encoding/json"
	"net/http"
	"time"

	"example.com/wsfanout"
)

// Server 管理 HTTP API。
type Server struct {
	hub    *wsfanout.Hub
	webDir string
	mux    *http.ServeMux
}

// New 创建 API 服务。
func New(hub *wsfanout.Hub, webDir string) *Server {
	s := &Server{hub: hub, webDir: webDir, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/health", s.handleHealth)
	s.mux.HandleFunc("/api/stats", s.handleStats)
	s.mux.HandleFunc("/api/rooms", s.handleRooms)
	s.mux.HandleFunc("/api/conns", s.handleConns)
	s.mux.HandleFunc("/api/broadcast", s.handleBroadcast)
	s.mux.HandleFunc("/api/join", s.handleJoin)
	s.mux.HandleFunc("/api/leave", s.handleLeave)
	if s.webDir != "" {
		s.mux.Handle("/", http.FileServer(http.Dir(s.webDir)))
	}
}

// Handler 返回 http.Handler。
func (s *Server) Handler() http.Handler { return s.mux }

// ListenAndServe 启动。
func (s *Server) ListenAndServe(addr string) error {
	srv := &http.Server{Addr: addr, Handler: s.mux, ReadHeaderTimeout: 5 * time.Second}
	return srv.ListenAndServe()
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, map[string]any{"ok": !s.hub.IsClosed()})
}

func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.hub.Stats())
}

func (s *Server) handleRooms(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.hub.Snapshot())
}

func (s *Server) handleConns(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, 200, s.hub.ListConns())
}

func (s *Server) handleJoin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", 405)
		return
	}
	var req struct {
		Room string `json:"room"`
		ID   string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	id, err := s.hub.Join(req.Room, req.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, 200, map[string]string{"id": id})
}

func (s *Server) handleLeave(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", 405)
		return
	}
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	if err := s.hub.Leave(req.ID); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, 200, map[string]string{"ok": "1"})
}

func (s *Server) handleBroadcast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", 405)
		return
	}
	var req struct {
		Room string `json:"room"`
		Type string `json:"type"`
		Body string `json:"body"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err := s.hub.Broadcast(req.Room, wsfanout.Message{Type: req.Type, Body: []byte(req.Body)})
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	writeJSON(w, 200, map[string]string{"ok": "1"})
}
