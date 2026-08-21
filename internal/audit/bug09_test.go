package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBug09_AuditFileClosedBeforeRotate(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.jsonl")
	l, err := Open(path, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Write("join", "lobby", "c1", nil); err != nil {
		t.Fatal(err)
	}
	if err := l.Rotate(); err != nil {
		t.Fatalf("rotate failed (handle still held?): %v", err)
	}
	bak := path + ".1"
	if _, err := os.Stat(bak); err != nil {
		t.Fatalf("expected rotated file: %v", err)
	}
	// 轮转后应能再次写入新文件
	if err := l.Write("leave", "lobby", "c1", nil); err != nil {
		t.Fatal(err)
	}
	_ = l.Close()
}
