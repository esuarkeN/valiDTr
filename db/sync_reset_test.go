//go:build cgo

package db_test

import (
	"testing"
	"time"

	"github.com/esuarkeN/valiDTr/db"
)

func TestSyncFromYAMLReset_InsertsEpochActiveEntries(t *testing.T) {
	initTestDB(t)

	cfgPath := writeTempFile(t, "config.yml", `
developers:
  - email: dev@example.com
    name: Dev
    keys:
      - id: ABCDEF01
`)

	if err := db.SyncFromYAMLReset(cfgPath, true); err != nil {
		t.Fatalf("SyncFromYAMLReset: %v", err)
	}

	epoch := time.Unix(0, 0).UTC()

	active, err := db.IsDeveloperActiveAt("dev@example.com", epoch)
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt: %v", err)
	}
	if !active {
		t.Fatalf("expected developer active at epoch after reset sync")
	}

	active, err = db.IsKeyActiveForDeveloperAt("dev@example.com", "ABCDEF01", epoch)
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt: %v", err)
	}
	if !active {
		t.Fatalf("expected key active at epoch after reset sync")
	}
}

func TestSyncFromYAMLReset_IsAtomicOnParseOrValidationError(t *testing.T) {
	initTestDB(t)

	initialCfg := writeTempFile(t, "initial.yml", `
developers:
  - email: dev@example.com
    name: Dev
    keys:
      - id: ABCDEF01
`)
	if err := db.SyncFromYAMLReset(initialCfg, true); err != nil {
		t.Fatalf("SyncFromYAMLReset(initial): %v", err)
	}

	badCfg := writeTempFile(t, "bad.yml", `
developers:
  - email: dev@example.com
    name: Dev
    keys:
      - id: ABCDEF01
  - email: other@example.com
    name: Other
    keys:
      - id:
`)
	if err := db.SyncFromYAMLReset(badCfg, true); err == nil {
		t.Fatalf("expected SyncFromYAMLReset(bad) to fail")
	}

	epoch := time.Unix(0, 0).UTC()

	active, err := db.IsDeveloperActiveAt("dev@example.com", epoch)
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt(dev@example.com): %v", err)
	}
	if !active {
		t.Fatalf("expected original developer state to remain after failed reset sync")
	}

	active, err = db.IsKeyActiveForDeveloperAt("dev@example.com", "ABCDEF01", epoch)
	if err != nil {
		t.Fatalf("IsKeyActiveForDeveloperAt(dev@example.com, ABCDEF01): %v", err)
	}
	if !active {
		t.Fatalf("expected original key state to remain after failed reset sync")
	}

	active, err = db.IsDeveloperActiveAt("other@example.com", epoch)
	if err != nil {
		t.Fatalf("IsDeveloperActiveAt(other@example.com): %v", err)
	}
	if active {
		t.Fatalf("expected failed sync not to partially add other@example.com")
	}
}
