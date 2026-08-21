package wsfanout

import (
	"context"
	"time"

	"example.com/wsfanout/internal/errs"
	"example.com/wsfanout/internal/hubcore"
	"example.com/wsfanout/internal/payload"
	"example.com/wsfanout/internal/validate"
)

// Broadcast 向房间内所有成员扇出消息；载荷在入口深拷贝。
func (h *Hub) Broadcast(roomName string, msg Message) error {
	return h.BroadcastContext(context.Background(), roomName, msg)
}

// BroadcastContext 支持取消的房间广播。
func (h *Hub) BroadcastContext(ctx context.Context, roomName string, msg Message) error {
	if err := ctx.Err(); err != nil {
		return errs.WrapCanceled(err)
	}
	// 关停流程会把 h.rooms/h.conns 置 nil；Close 后补发广播必须返回 ErrClosed
	// 而非命中 nil map 写入导致 panic。先做一次无锁快检，避免无谓的编码开销。
	if h.closed.Load() {
		return ErrClosed
	}

	if err := validate.RoomName(roomName); err != nil {
		return err
	}

	cloned := payload.CloneMessage(payload.Message{
		Type:      msg.Type,
		Body:      msg.Body,
		Headers:   msg.Headers,
		Timestamp: msg.Timestamp,
	})
	if cloned.Timestamp.IsZero() {
		cloned.Timestamp = h.clock.Now()
	}

	encoded, err := h.encodeLocked(cloned)
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()
	// 持锁再检：Close 在本路径与编码之间完成关停时，rooms/conns 可能已置 nil，
	// 此处拦截以避免 ensureRoomLocked 对 nil map 写入而 panic。
	if h.closed.Load() {
		return ErrClosed
	}

	r := h.ensureRoomLocked(roomName)
	members := r.Members()
	targets := make([]*hubcore.Target, 0, len(members))
	for _, id := range members {
		c, ok := h.conns[id]
		if !ok {
			continue
		}
		targets = append(targets, &hubcore.Target{ID: id, Conn: c})
	}
	res := hubcore.Fanout(ctx, targets, encoded, func() error {
		return ctx.Err()
	})
	h.broadcasts.Add(1)
	h.metrics.IncBroadcasts(1)
	h.dropped.Add(int64(res.Dropped))
	h.writeErrors.Add(int64(res.WriteErrors))
	if h.auditor != nil {
		_ = h.auditor.Write("broadcast", roomName, "", map[string]string{
			"type": cloned.Type,
			"n":    itoa(len(targets)),
		})
	}
	if res.Canceled {
		return errs.WrapCanceled(context.Canceled)
	}
	if res.LastWriteErr != nil {
		return res.LastWriteErr
	}
	return nil
}

// SendTo 向单个连接同步推送（走 WriteDirect，写出错误经 WrapWrite 上浮）。
func (h *Hub) SendTo(connID string, msg Message) error {
	if h.closed.Load() {
		return ErrClosed
	}
	cloned := payload.CloneMessage(payload.Message{
		Type:      msg.Type,
		Body:      msg.Body,
		Headers:   msg.Headers,
		Timestamp: msg.Timestamp,
	})
	if cloned.Timestamp.IsZero() {
		cloned.Timestamp = h.clock.Now()
	}
	encoded, err := h.encodeLocked(cloned)
	if err != nil {
		return err
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	c, ok := h.conns[connID]
	if !ok {
		return ErrNotFound
	}
	return c.WriteDirect(context.Background(), encoded)
}

func (h *Hub) encodeLocked(m payload.Message) ([]byte, error) {
	enc := h.encoder
	if enc == nil {
		return nil, ErrNoEncoder
	}
	b, err := enc.Encode(m)
	if err != nil {
		return nil, errs.WrapEncode(err)
	}
	return b, nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for n > 0 {
		i--
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[i:])
}

// DrainConn 阻塞直到连接写队列被消费或超时（测试用）。
func (h *Hub) DrainConn(connID string, timeout time.Duration) error {
	h.mu.RLock()
	c, ok := h.conns[connID]
	h.mu.RUnlock()
	if !ok {
		return ErrNotFound
	}
	return c.WaitEmpty(timeout)
}
