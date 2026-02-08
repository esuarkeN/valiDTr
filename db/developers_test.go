//go:build cgo

package db_test

import (
	"testing"
	"time"

	"github.com/esuarkeN/valiDTr/db"
)

func TestDeveloperLifecycleWindows(t *testing.T) {
	initTestDB(t)

	email := "dev@example.com"

	addedAt := mustUTC(t, "2020-01-01T00:00:00Z")
	removedAt := mustUTC(t, "2020-01-02T00:00:00Z")
	readdedAt := mustUTC(t, "2020-01-03T00:00:00Z")

	if err := db.AddDeveloper(email, "Dev", addedAt); err != nil {
		t.Fatalf("AddDeveloper: %v", err)
	}

	active, err := db.IsDeveloperActiveAt(email, addedAt.Add(-time.Second))
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt(before add): %v", err)
	}
	if active {
		t.Fatalf("expected inactive before added_at")
	}

	active, err = db.IsDeveloperActiveAt(email, addedAt)
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt(at add): %v", err)
	}
	if !active {
		t.Fatalf("expected active at added_at")
	}

	if err := db.RemoveDeveloper(email, removedAt); err != nil {
		t.Fatalf("RemoveDeveloper: %v", err)
	}

	active, err = db.IsDeveloperActiveAt(email, removedAt.Add(-time.Second))
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt(before remove): %v", err)
	}
	if !active {
		t.Fatalf("expected active before removed_at")
	}

	active, err = db.IsDeveloperActiveAt(email, removedAt)
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt(at remove): %v", err)
	}
	if active {
		t.Fatalf("expected inactive at removed_at")
	}

	if err := db.AddDeveloper(email, "", readdedAt); err != nil {
		t.Fatalf("AddDeveloper(readd): %v", err)
	}

	active, err = db.IsDeveloperActiveAt(email, removedAt.Add(time.Second))
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt(between remove and readd): %v", err)
	}
	if active {
		t.Fatalf("expected inactive between remove and re-add")
	}

	active, err = db.IsDeveloperActiveAt(email, readdedAt)
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt(at readd): %v", err)
	}
	if !active {
		t.Fatalf("expected active at re-added_at")
	}
}
