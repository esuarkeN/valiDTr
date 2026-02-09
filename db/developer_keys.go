package db

import (
	"database/sql"
	"fmt"
	"time"
)

type DevKey struct {
	DeveloperEmail string
	KeyID          string
	AddedAt        time.Time
	RevokedAt      sql.NullTime
}

func AddKeyToDeveloper(email, keyID string, addedAt time.Time) error {
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
		return fmt.Errorf("developer not found: %s", email)
	}

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO developer_keys(developer_id, key_id, added_at) VALUES(?, ?, ?)`,
		d.ID, keyID, addedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("add key to developer: %w", err)
	}
	return nil
}

func RevokeDeveloperKey(email, keyID string, revokedAt time.Time) error {
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
		return fmt.Errorf("developer not found: %s", email)
	}

	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec(
		`UPDATE developer_keys
         SET revoked_at = ?
         WHERE developer_id = ? AND upper(key_id) = upper(?) AND revoked_at IS NULL`,
		revokedAt.UTC(), d.ID, keyID,
	)
	if err != nil {
		return fmt.Errorf("revoke key: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("no active key mapping found for %s -> %s", email, keyID)
	}
	return nil
}

func IsKeyActiveForDeveloperAt(email, keyID string, t time.Time) (bool, error) {
	email = normalizeEmail(email)
	keyID = normalizeKeyID(keyID)
	d, err := GetDeveloperByEmail(email)
	if err != nil {
		return false, err
	}
	if d == nil {
		return false, nil
	}

	db, err := openDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var dummy int
	err = db.QueryRow(
		`SELECT 1
         FROM developer_keys
         WHERE developer_id = ?
           AND upper(key_id) = upper(?)
           AND added_at <= ?
           AND (revoked_at IS NULL OR revoked_at > ?)
         ORDER BY added_at DESC
         LIMIT 1`,
		d.ID, keyID, t.UTC(), t.UTC(),
	).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func ListKeys(emailFilter string) ([]DevKey, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	q := `
SELECT d.email, k.key_id, k.added_at, k.revoked_at
FROM developer_keys k
JOIN developers d ON d.id = k.developer_id
`
	args := []any{}
	if emailFilter != "" {
		q += ` WHERE lower(d.email) = lower(?)`
		args = append(args, normalizeEmail(emailFilter))
	}
	q += ` ORDER BY d.email, k.key_id, k.added_at`

	rows, err := db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []DevKey
	for rows.Next() {
		var dk DevKey
		if err := rows.Scan(&dk.DeveloperEmail, &dk.KeyID, &dk.AddedAt, &dk.RevokedAt); err != nil {
			return nil, err
		}
		out = append(out, dk)
	}
	return out, nil
}
