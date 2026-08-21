// Package wsfanout 提供进程内 WebSocket 房间扇出中心。
//
// 能力概览：
//   - Join / Leave：连接加入或离开命名房间
//   - Broadcast / BroadcastContext / SendTo：房间广播与点对点推送
//   - Snapshot：导出房间成员列表（须拷贝，避免与内部表共享）
//   - ApplyRoomConfig：热更新房间策略；持久化失败须回滚内存
//   - Close：先 Flush 各连接写队列，再关闭连接与房间表
//
// 载荷在 Broadcast 入口深拷贝，调用方事后修改不会污染在途消息。
// Close 之后再 Broadcast 返回 ErrClosed，不得 panic。
package wsfanout
