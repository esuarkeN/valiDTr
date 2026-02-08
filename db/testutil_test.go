//go:build cgo

package db_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/esuarkeN/valiDTr/db"
)

func initTestDB(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	dbPath := filepath.Join(dir, "validtr.db")
	if err := db.InitDB(dbPath); err != nil {
		t.Fatalf("InitDB(%q): %v", dbPath, err)
	}
	return dbPath
}

func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()

	dir := t.TempDir()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
		t.Fatalf("WriteFile(%q): %v", p, err)
	}
	return p
}

func mustUTC(t *testing.T, ts string) time.Time {
	t.Helper()

	tt, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		t.Fatalf("Parse(%q): %v", ts, err)
	}
	return tt.UTC()
}
