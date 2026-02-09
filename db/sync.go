package db

import (
	"database/sql"
	"fmt"
)

func ResetAll() error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM developer_keys;`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM developer_status;`); err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM developers;`); err != nil {
		return err
	}

	return tx.Commit()
}

func UpsertDeveloper(email, name string, addedAt, removedAt sql.NullTime) error {
	email = normalizeEmail(email)
	if email == "" {
		return fmt.Errorf("developer email is required")
	}

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
INSERT INTO developers(email, name, added_at, removed_at)
VALUES(?, ?, ?, ?)
ON CONFLICT(email) DO UPDATE SET
  name = excluded.name,
  removed_at = excluded.removed_at
`, email, name, addedAt, removedAt)

	if err != nil {
		return fmt.Errorf("upsert developer: %w", err)
	}
	return nil
}

func InsertDeveloperKey(email, keyID string, addedAt, revokedAt sql.NullTime) error {
	email = normalizeEmail(email)
	if email == "" {
		return fmt.Errorf("developer email is required")
	}
	keyID = normalizeKeyID(keyID)
	if keyID == "" {
		return fmt.Errorf("key id is required")
	}

	d, err := GetDeveloperByEmail(email)
	if err != nil {
		return err
	}
	if d == nil {
		return fmt.Errorf("developer not found during key insert: %s", email)
	}

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
INSERT INTO developer_keys(developer_id, key_id, added_at, revoked_at)
VALUES(?, ?, ?, ?)
`, d.ID, keyID, addedAt, revokedAt)

	if err != nil {
		return fmt.Errorf("insert developer_key: %w", err)
	}
	return nil
}
