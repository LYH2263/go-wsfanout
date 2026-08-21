# go-wsfanout

进程内 WebSocket 房间扇出中心：Join / Leave / Broadcast / SendTo / Close，配套背压、审计、房间配置持久化与管理页。

## 构建与测试

```bash
go build ./...
go test ./... -count=1
```

## 管理守护进程

```bash
go run ./cmd/wsd -addr :8102 -web web
```

打开 http://127.0.0.1:8102/ 可查看房间、在线连接并试广播。

## 库能力摘要

- `Hub.Join` / `Leave`：房间成员管理
- `Hub.Broadcast` / `BroadcastContext`：房间扇出（载荷深拷贝）
- `Hub.SendTo`：按连接 ID 点对点推送
- `Hub.Snapshot`：导出房间成员快照（独立切片）
- `Hub.ApplyRoomConfig`：热更新房间策略；持久化失败不留下半应用状态
- `Hub.Close`：先 Flush 写队列再拆房间
