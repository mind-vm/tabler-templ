package contacts

import (
	"context"
	"errors"
	"testing"
)

func TestStore(t *testing.T) {
	ctx := context.Background()
	s := NewStore()
	all, _ := s.List(ctx, Query{})
	if len(all) != len(seed) || all[0].Name != "Ada Lovelace" {
		t.Fatalf("list: %d, first %q", len(all), all[0].Name)
	}
	bell, _ := s.List(ctx, Query{Search: "bell labs", Status: Active})
	if len(bell) != 1 || bell[0].Name != "Ken Thompson" {
		t.Fatalf("search: %+v", bell)
	}
	if _, err := s.Create(ctx, Contact{Email: "x", Status: "?"}); len(err.(FieldErrors)) != 3 {
		t.Fatalf("validation: %v", err)
	}
	c, _ := s.Get(ctx, "1")
	c.Notes = "a"
	if _, err := s.Update(ctx, c); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Update(ctx, c); !errors.Is(err, ErrConflict) {
		t.Fatalf("second update with the old version: %v", err)
	}
	if _, err := s.Delete(ctx, "1"); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Get(ctx, "1"); !errors.Is(err, ErrNotFound) {
		t.Fatal("deleted contact still there")
	}
	if _, err := s.Restore(ctx, "1"); err != nil {
		t.Fatal(err)
	}
	ev, _ := s.Activity(ctx, "1")
	if ev[0].Text != "Restored after deletion" {
		t.Fatalf("activity: %+v", ev)
	}
}
