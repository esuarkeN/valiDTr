package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// GPGKey represents a stored key in the database
type GPGKey struct {
	KeyID  string
	Status string
}

// ListGPGKeys fetches all stored GPG keys from the database
func ListGPGKeys() ([]GPGKey, error) {
	db, err := sql.Open("sqlite3", "gpgkeys.db")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT key_id, status FROM keys")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []GPGKey
	for rows.Next() {
		var key GPGKey
		if err := rows.Scan(&key.KeyID, &key.Status); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	if len(keys) == 0 {
		fmt.Println("No GPG keys found in the database.")
	}

	return keys, nil
}
