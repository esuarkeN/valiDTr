package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

func placeholders(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.TrimRight(strings.Repeat("?,", n), ",")
}

func SyncFromYAMLReset(path string, reset bool) error {

	return SyncFromYAML(path, reset)
}

func SyncFromYAMLReconcile(path string) error {
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

	now := time.Now().UTC()
	epoch := time.Unix(0, 0).UTC()

	cfgEmails := make(map[string]struct{}, len(cfg.Developers))

	for _, d := range cfg.Developers {
		if d.Email == "" {
			return fmt.Errorf("developer missing email")
		}
		cfgEmails[d.Email] = struct{}{}
		var devID int64
		err := tx.QueryRow(`SELECT id FROM developers WHERE email = ?`, d.Email).Scan(&devID)
		if err == sql.ErrNoRows {
			res, e := tx.Exec(
				`INSERT INTO developers(email, name, added_at, removed_at) VALUES(?, ?, ?, NULL)`,
				d.Email, nullIfEmpty(d.Name), epoch,
			)
			if e != nil {
				return fmt.Errorf("insert developer %s: %w", d.Email, e)
			}
			devID, _ = res.LastInsertId()
			if _, e := tx.Exec(
				`INSERT INTO developer_status(developer_id, added_at, removed_at)
				 VALUES(?, ?, NULL)`,
				devID, epoch,
			); e != nil {
				return fmt.Errorf("insert developer status %s: %w", d.Email, e)
			}
		} else if err != nil {
			return fmt.Errorf("select developer %s: %w", d.Email, err)
		} else {
			_, e := tx.Exec(
				`UPDATE developers SET name = ?, removed_at = NULL WHERE id = ?`,
				nullIfEmpty(d.Name), devID,
			)
			if e != nil {
				return fmt.Errorf("update developer %s: %w", d.Email, e)
			}
			var dummy int
			err = tx.QueryRow(
				`SELECT 1 FROM developer_status
				 WHERE developer_id = ? AND removed_at IS NULL
				 LIMIT 1`,
				devID,
			).Scan(&dummy)
			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("check developer status %s: %w", d.Email, err)
			}
			if err == sql.ErrNoRows {
				if _, e := tx.Exec(
					`INSERT INTO developer_status(developer_id, added_at, removed_at)
					 VALUES(?, ?, NULL)`,
					devID, now,
				); e != nil {
					return fmt.Errorf("reactivate developer %s: %w", d.Email, e)
				}
			}
		}

		keySet := map[string]struct{}{}
		for _, k := range d.Keys {
			if k.ID == "" {
				return fmt.Errorf("developer %s has a key with empty id", d.Email)
			}
			keySet[k.ID] = struct{}{}
			var dummy int
			err = tx.QueryRow(
				`SELECT 1 FROM developer_keys
				 WHERE developer_id = ? AND key_id = ? AND revoked_at IS NULL
				 LIMIT 1`,
				devID, k.ID,
			).Scan(&dummy)
			if err != nil && err != sql.ErrNoRows {
				return fmt.Errorf("check active key %s for %s: %w", k.ID, d.Email, err)
			}
			if err == sql.ErrNoRows {
				// If the key existed before and was revoked, re-add with "now" to preserve the revoke window.
				addedAt := epoch
				err = tx.QueryRow(
					`SELECT 1 FROM developer_keys
					 WHERE developer_id = ? AND key_id = ?
					 LIMIT 1`,
					devID, k.ID,
				).Scan(&dummy)
				if err != nil && err != sql.ErrNoRows {
					return fmt.Errorf("check prior key %s for %s: %w", k.ID, d.Email, err)
				}
				if err != sql.ErrNoRows {
					addedAt = now
				}
				if _, e2 := tx.Exec(
					`INSERT INTO developer_keys(developer_id, key_id, added_at, revoked_at)
					 VALUES(?, ?, ?, NULL)`,
					devID, k.ID, addedAt,
				); e2 != nil {
					return fmt.Errorf("insert key %s for %s: %w", k.ID, d.Email, e2)
				}
			}
		}

		if len(keySet) == 0 {
			_, e := tx.Exec(
				`UPDATE developer_keys SET revoked_at = ?
				 WHERE developer_id = ? AND revoked_at IS NULL`,
				now, devID,
			)
			if e != nil {
				return fmt.Errorf("revoke missing keys for %s: %w", d.Email, e)
			}
		} else {
			keys := make([]string, 0, len(keySet))
			for k := range keySet {
				keys = append(keys, k)
			}
			args := make([]any, 0, 2+len(keys))
			args = append(args, now, devID)
			for _, k := range keys {
				args = append(args, k)
			}
			q := fmt.Sprintf(
				`UPDATE developer_keys SET revoked_at = ?
				 WHERE developer_id = ? AND revoked_at IS NULL
				   AND key_id NOT IN (%s)`,
				placeholders(len(keys)),
			)
			if _, e := tx.Exec(q, args...); e != nil {
				return fmt.Errorf("revoke missing keys for %s: %w", d.Email, e)
			}
		}
	}

	if len(cfgEmails) == 0 {
		if _, err := tx.Exec(`UPDATE developer_status SET removed_at = ? WHERE removed_at IS NULL`, now); err != nil {
			return fmt.Errorf("remove all devs: %w", err)
		}
		if _, err := tx.Exec(`UPDATE developers SET removed_at = ? WHERE removed_at IS NULL`, now); err != nil {
			return fmt.Errorf("remove all devs: %w", err)
		}
	} else {
		emails := make([]string, 0, len(cfgEmails))
		for e := range cfgEmails {
			emails = append(emails, e)
		}
		args := make([]any, 0, 1+len(emails))
		args = append(args, now)
		for _, e := range emails {
			args = append(args, e)
		}

		q := fmt.Sprintf(
			`UPDATE developer_status SET removed_at = ?
			 WHERE removed_at IS NULL AND developer_id IN (
			   SELECT id FROM developers WHERE email NOT IN (%s)
			 )`,
			placeholders(len(emails)),
		)
		if _, err := tx.Exec(q, args...); err != nil {
			return fmt.Errorf("remove missing devs: %w", err)
		}

		q = fmt.Sprintf(
			`UPDATE developers SET removed_at = ?
			 WHERE removed_at IS NULL AND email NOT IN (%s)`,
			placeholders(len(emails)),
		)
		if _, err := tx.Exec(q, args...); err != nil {
			return fmt.Errorf("remove missing devs: %w", err)
		}
	}

	return tx.Commit()
}

func nullIfEmpty(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}
