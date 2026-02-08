//go:build cgo

package db_test

import (
	"testing"
	"time"

	"github.com/esuarkeN/valiDTr/db"
)

func TestSyncFromYAMLReconcile_KeyRevokeAndReaddKeepsOldWindowValid(t *testing.T) {
	initTestDB(t)

	cfg1 := writeTempFile(t, "config1.yml", `
developers:
  - email: dev@example.com
    name: Dev
    keys:
      - id: ABCDEF01
`)
	if err := db.SyncFromYAMLReconcile(cfg1); err != nil {
		t.Fatalf("SyncFromYAMLReconcile(cfg1): %v", err)
	}

	cfg2 := writeTempFile(t, "config2.yml", `
developers:
  - email: dev@example.com
    name: Dev
    keys: []
`)
	if err := db.SyncFromYAMLReconcile(cfg2); err != nil {
		t.Fatalf("SyncFromYAMLReconcile(cfg2): %v", err)
	}

	keys, err := db.ListKeys("dev@example.com")
	if err != nil {
		t.Fatalf("ListKeys: %v", err)
	}
	if len(keys) != 1 {
		t.Fatalf("expected 1 key mapping after revoke, got %d", len(keys))
	}
	if keys[0].KeyID != "ABCDEF01" {
		t.Fatalf("expected key ABCDEF01, got %q", keys[0].KeyID)
	}
	if !keys[0].RevokedAt.Valid {
		t.Fatalf("expected key to be revoked after config removal")
	}
	revokedAt := keys[0].RevokedAt.Time.UTC()

	time.Sleep(50 * time.Millisecond)

	cfg3 := writeTempFile(t, "config3.yml", `
developers:
  - email: dev@example.com
    name: Dev
    keys:
      - id: ABCDEF01
`)
	if err := db.SyncFromYAMLReconcile(cfg3); err != nil {
		t.Fatalf("SyncFromYAMLReconcile(cfg3): %v", err)
	}

	keys, err = db.ListKeys("dev@example.com")
	if err != nil {
		t.Fatalf("ListKeys(after readd): %v", err)
	}
	if len(keys) != 2 {
		t.Fatalf("expected 2 key mappings after re-add, got %d", len(keys))
	}

	var readdAddedAt time.Time
	for _, k := range keys {
		if k.KeyID == "ABCDEF01" && !k.RevokedAt.Valid && !k.AddedAt.Equal(time.Unix(0, 0).UTC()) {
			readdAddedAt = k.AddedAt.UTC()
		}
	}
	if readdAddedAt.IsZero() {
		t.Fatalf("expected a new active key mapping after re-add")
	}
	if !readdAddedAt.After(revokedAt) {
		t.Fatalf("expected re-add added_at > revoked_at, got %s <= %s", readdAddedAt.Format(time.RFC3339Nano), revokedAt.Format(time.RFC3339Nano))
	}

	oldCommitTime := revokedAt.Add(-time.Second)
	between := revokedAt.Add(readdAddedAt.Sub(revokedAt) / 2)
	newCommitTime := readdAddedAt.Add(time.Millisecond)

	active, err := db.IsKeyActiveForDeveloperAt("dev@example.com", "ABCDEF01", oldCommitTime)
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(old): %v", err)
	}
	if !active {
		t.Fatalf("expected active before revocation even after re-add")
	}

	active, err = db.IsKeyActiveForDeveloperAt("dev@example.com", "ABCDEF01", between)
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(between): %v", err)
	}
	if active {
		t.Fatalf("expected inactive between revoke and re-add")
	}

	active, err = db.IsKeyActiveForDeveloperAt("dev@example.com", "ABCDEF01", newCommitTime)
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(new): %v", err)
	}
	if !active {
		t.Fatalf("expected active after re-add")
	}
}
