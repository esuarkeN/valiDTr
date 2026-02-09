package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type yamlConfig struct {
	Developers []yamlDeveloper `yaml:"developers"`
}

type yamlDeveloper struct {
	Email string    `yaml:"email"`
	Name  string    `yaml:"name"`
	Keys  []yamlKey `yaml:"keys"`
}

type yamlKey struct {
	ID string `yaml:"id"`
}

// Policy note:
//   - For easy CI configuration, we treat config entries as "active now".
//   - added_at defaults to Unix epoch so historical policy can validate old commits
//     *as long as you don't rely on commit dates for security-critical gating*.
func SyncFromYAML(path string, reset bool) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg yamlConfig
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return fmt.Errorf("parse yaml: %w", err)
	}

	dbh, err := openDB()
	if err != nil {
		return err
	}
	defer dbh.Close()

	tx, err := dbh.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	epoch := sql.NullTime{Time: time.Unix(0, 0).UTC(), Valid: true}
	none := sql.NullTime{Valid: false}

	if reset {
		if _, err := tx.Exec(`DELETE FROM developer_keys;`); err != nil {
			return err
		}
		if _, err := tx.Exec(`DELETE FROM developer_status;`); err != nil {
			return err
		}
		if _, err := tx.Exec(`DELETE FROM developers;`); err != nil {
			return err
		}
	}

	for _, d := range cfg.Developers {
		email := normalizeEmail(d.Email)
		if email == "" {
			return fmt.Errorf("developer missing email")
		}

		_, err = tx.Exec(`
INSERT INTO developers(email, name, added_at, removed_at)
VALUES(?, ?, ?, ?)
ON CONFLICT(email) DO UPDATE SET
  name = excluded.name,
  removed_at = excluded.removed_at
`, email, nullIfEmpty(d.Name), epoch, none)
		if err != nil {
			return fmt.Errorf("upsert developer: %w", err)
		}

		var devID int64
		if err := tx.QueryRow(
			`SELECT id FROM developers WHERE lower(email) = lower(?)`,
			email,
		).Scan(&devID); err != nil {
			return fmt.Errorf("select developer after upsert: %w", err)
		}

		var activeDummy int
		err = tx.QueryRow(
			`SELECT 1 FROM developer_status WHERE developer_id = ? AND removed_at IS NULL LIMIT 1`,
			devID,
		).Scan(&activeDummy)
		if err != nil && err != sql.ErrNoRows {
			return fmt.Errorf("check active developer status: %w", err)
		}
		if err == sql.ErrNoRows {
			if _, err := tx.Exec(
				`INSERT INTO developer_status(developer_id, added_at, removed_at) VALUES(?, ?, NULL)`,
				devID, epoch.Time,
			); err != nil {
				return fmt.Errorf("insert developer status: %w", err)
			}
		}

		for _, k := range d.Keys {
			keyID := normalizeKeyID(k.ID)
			if keyID == "" {
				return fmt.Errorf("developer %s has key with empty id", email)
			}
			if _, err := tx.Exec(
				`INSERT INTO developer_keys(developer_id, key_id, added_at, revoked_at) VALUES(?, ?, ?, ?)`,
				devID, keyID, epoch, none,
			); err != nil {
				return fmt.Errorf("insert developer key: %w", err)
			}
		}
	}

	return tx.Commit()
}
