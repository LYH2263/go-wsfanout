package wsfanout_test

import (
	"testing"

	"example.com/wsfanout"
)

func TestBug02_SnapshotMembersSliceAlias(t *testing.T) {
	h := wsfanout.New()
	defer h.Close()
	h.MustJoin("r1", "a")
	h.MustJoin("r1", "b")
	info, err := h.SnapshotRoom("r1")
	if err != nil {
		t.Fatal(err)
	}
	if len(info.Members) < 2 {
		t.Fatalf("want members, got %#v", info.Members)
	}
	orig0 := info.Members[0]
	info.Members[0] = "TAMPERED"
	info2, err := h.SnapshotRoom("r1")
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range info2.Members {
		if m == "TAMPERED" {
			t.Fatalf("snapshot members shared with internal table; orig=%q now=%v", orig0, info2.Members)
		}
	}
}
