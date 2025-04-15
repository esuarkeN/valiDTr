package db

import (
	"database/sql"
	"log"

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
		status TEXT
	);
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

	_, err = db.Exec("INSERT INTO keys (key_id, status) VALUES (?, ?)", keyID, "active")
	return err
}

func RevokeGPGKey(keyID string) error {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("UPDATE keys SET status = ? WHERE key_id = ?", "revoked", keyID)
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
