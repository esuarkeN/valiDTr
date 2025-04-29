// valiDTr_test.go
package db_test

import (
	"os"
	"testing"
	"valiDTr/db"
)

const testKeyID = "ABCDEF123456"

func setupTestDB() {
	_ = os.Remove("gpgkeys.db")
	db.InitDB()
}

func TestAddAndListKey(t *testing.T) {
	setupTestDB()
	err := db.AddGPGKey(testKeyID)
	if err != nil {
		t.Fatalf("Failed to add key: %v", err)
	}

	keys, err := db.ListGPGKeys()
	if err != nil {
		t.Fatalf("Failed to list keys: %v", err)
	}
	if len(keys) != 1 || keys[0].KeyID != testKeyID {
		t.Errorf("Expected one key with ID %s, got: %+v", testKeyID, keys)
	}
}

func TestDuplicateKey(t *testing.T) {
	setupTestDB()
	_ = db.AddGPGKey(testKeyID)
	err := db.AddGPGKey(testKeyID)
	if err == nil {
		t.Error("Expected error on duplicate key, got nil")
	}
}

func TestRevokeKey(t *testing.T) {
	setupTestDB()
	_ = db.AddGPGKey(testKeyID)
	err := db.RevokeGPGKey(testKeyID)
	if err != nil {
		t.Fatalf("Failed to revoke key: %v", err)
	}

	active, err := db.IsKeyActive(testKeyID)
	if err != nil {
		t.Fatalf("Error checking key status: %v", err)
	}
	if active {
		t.Error("Expected key to be revoked, but it is still active")
	}
}

func TestIsKeyActive_UnknownKey(t *testing.T) {
	setupTestDB()
	active, err := db.IsKeyActive("NONEXISTENT")
	if err != nil {
		t.Fatalf("Error for unknown key: %v", err)
	}
	if active {
		t.Error("Expected unknown key to be inactive")
	}
}
