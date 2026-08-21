package wsfanout_test

import (
	"testing"

	"example.com/wsfanout"
	"example.com/wsfanout/internal/persist"
)

func TestBug06_PersistFailureDoesNotApply(t *testing.T) {
	store := persist.NewMemoryStore()
	h := wsfanout.New(wsfanout.WithStore(store))
	defer h.Close()
	h.MustJoin("lobby", "c1")
	before, err := h.SnapshotRoom("lobby")
	if err != nil {
		t.Fatal(err)
	}
	store.SetFail(true)
	err = h.ApplyRoomConfig("lobby", wsfanout.RoomPolicy{
		MaxMembers: 3,
		QueueSize:  8,
		DropOldest: false,
		KickSlowMs: 100,
		AllowAnon:  false,
	})
	if err == nil {
		t.Fatal("expected persist failure")
	}
	after, err := h.SnapshotRoom("lobby")
	if err != nil {
		t.Fatal(err)
	}
	if after.Policy.MaxMembers != before.Policy.MaxMembers {
		t.Fatalf("policy applied despite persist fail: %+v -> %+v", before.Policy, after.Policy)
	}
	if after.Policy.AllowAnon != before.Policy.AllowAnon {
		t.Fatalf("AllowAnon half-applied: %+v", after.Policy)
	}
}
