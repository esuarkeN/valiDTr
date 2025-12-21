package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var dbFile string

func InitDB(path string) error {
	dbFile = path
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		return err
	}

	schema := `
CREATE TABLE IF NOT EXISTS developers (
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  email      TEXT NOT NULL UNIQUE,
  name       TEXT,
  added_at   TIMESTAMP NOT NULL,
  removed_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS developer_keys (
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  developer_id INTEGER NOT NULL,
  key_id       TEXT NOT NULL,
  added_at     TIMESTAMP NOT NULL,
  revoked_at   TIMESTAMP,
  FOREIGN KEY(developer_id) REFERENCES developers(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_devkeys_dev ON developer_keys(developer_id);
CREATE INDEX IF NOT EXISTS idx_devkeys_key ON developer_keys(key_id);
`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("init schema: %w", err)
	}

	return nil
}

func openDB() (*sql.DB, error) {
	if dbFile == "" {
		dbFile = "validtr.db"
	}
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(`PRAGMA foreign_keys = ON;`); err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
