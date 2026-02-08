//go:build cgo

package db_test

import (
	"testing"
	"time"

	"github.com/esuarkeN/valiDTr/db"
)

func TestKeyLifecycleWindows_ReaddPreservesOldCommits(t *testing.T) {
	initTestDB(t)

	email := "dev@example.com"
	keyID := "DEADBEEF"

	devAddedAt := mustUTC(t, "2020-01-01T00:00:00Z")
	keyAddedAt := mustUTC(t, "2020-01-02T00:00:00Z")
	keyRevokedAt := mustUTC(t, "2020-01-03T00:00:00Z")
	keyReaddedAt := mustUTC(t, "2020-01-04T00:00:00Z")

	if err := db.AddDeveloper(email, "Dev", devAddedAt); err != nil {
		t.Fatalf("AddDeveloper: %v", err)
	}
	if err := db.AddKeyToDeveloper(email, keyID, keyAddedAt); err != nil {
		t.Fatalf("AddKeyToDeveloper: %v", err)
	}

	active, err := db.IsKeyActiveForDeveloperAt(email, keyID, keyAddedAt.Add(-time.Second))
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(before add): %v", err)
	}
	if active {
		t.Fatalf("expected inactive before added_at")
	}

	if err := db.RevokeDeveloperKey(email, keyID, keyRevokedAt); err != nil {
		t.Fatalf("RevokeDeveloperKey: %v", err)
	}

	if err := db.AddKeyToDeveloper(email, keyID, keyReaddedAt); err != nil {
		t.Fatalf("AddKeyToDeveloper(readd): %v", err)
	}

	active, err = db.IsKeyActiveForDeveloperAt(email, keyID, keyRevokedAt.Add(-time.Second))
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(old commit): %v", err)
	}
	if !active {
		t.Fatalf("expected active before revocation even after re-add")
	}

	active, err = db.IsKeyActiveForDeveloperAt(email, keyID, keyRevokedAt)
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(at revoke): %v", err)
	}
	if active {
		t.Fatalf("expected inactive at revoked_at")
	}

	active, err = db.IsKeyActiveForDeveloperAt(email, keyID, keyRevokedAt.Add(time.Second))
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(between revoke and readd): %v", err)
	}
	if active {
		t.Fatalf("expected inactive between revocation and re-add")
	}

	active, err = db.IsKeyActiveForDeveloperAt(email, keyID, keyReaddedAt)
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(at readd): %v", err)
	}
	if !active {
		t.Fatalf("expected active at re-added_at")
	}
}
