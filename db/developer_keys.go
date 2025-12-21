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
         WHERE developer_id = ? AND key_id = ? AND revoked_at IS NULL`,
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

	var addedAt time.Time
	var revokedAt sql.NullTime
	err = db.QueryRow(
		`SELECT added_at, revoked_at
         FROM developer_keys
         WHERE developer_id = ? AND key_id = ?
         ORDER BY added_at DESC
         LIMIT 1`,
		d.ID, keyID,
	).Scan(&addedAt, &revokedAt)

	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if t.Before(addedAt) {
		return false, nil
	}
	if revokedAt.Valid && !t.Before(revokedAt.Time) {
		return false, nil
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
		q += ` WHERE d.email = ?`
		args = append(args, emailFilter)
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
