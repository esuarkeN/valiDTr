package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "gpgkeys.db"

func InitDB() {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS keys (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		key_id TEXT UNIQUE,
		status      TEXT NOT NULL,             
    	added_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    	revoked_at  TIMESTAMP 
	);
	ALTER TABLE keys ADD COLUMN added_at   TIMESTAMP;
	ALTER TABLE keys ADD COLUMN revoked_at TIMESTAMP;
	UPDATE keys SET added_at = COALESCE(timestamp, CURRENT_TIMESTAMP);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}
}

func AddGPGKey(keyID string) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
        INSERT INTO keys (key_id, status, added_at)
        VALUES (?, 'active', CURRENT_TIMESTAMP)
        ON CONFLICT(key_id) DO UPDATE SET
            status   = 'active',
            added_at = COALESCE(added_at, CURRENT_TIMESTAMP),
            revoked_at = NULL
    `, keyID)
	return err
}

func RevokeGPGKey(keyID string) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
        UPDATE keys
        SET status = 'revoked',
            revoked_at = COALESCE(revoked_at, CURRENT_TIMESTAMP)
        WHERE key_id = ?
    `, keyID)
	return err
}

func IsKeyActive(keyID string) (bool, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var status string
	err = db.QueryRow("SELECT status FROM keys WHERE key_id = ?", keyID).Scan(&status)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return status == "active", nil
}

func IsKeyActiveAt(keyID string, t time.Time) (bool, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var status string
	var addedAt time.Time
	var revokedAt sql.NullTime

	err = db.QueryRow(`
        SELECT status, added_at, revoked_at
        FROM keys
        WHERE key_id = ?
    `, keyID).Scan(&status, &addedAt, &revokedAt)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if status != "active" && status != "revoked" {
		return false, nil
	}

	if t.Before(addedAt) {
		return false, nil
	}

	if revokedAt.Valid && !t.Before(revokedAt.Time) {
		return false, nil
	}

	return true, nil
}
