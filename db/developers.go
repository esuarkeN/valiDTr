package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Developer struct {
	ID        int64
	Email     string
	Name      sql.NullString
	AddedAt   time.Time
	RemovedAt sql.NullTime
}

func AddDeveloper(email, name string, addedAt time.Time) error {
	email = normalizeEmail(email)
	if email == "" {
		return fmt.Errorf("developer email is required")
	}

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

	var devID int64
	err = tx.QueryRow(
		`SELECT id FROM developers WHERE lower(email) = lower(?)`,
		email,
	).Scan(&devID)
	if err == sql.ErrNoRows {
		res, e := tx.Exec(
			`INSERT INTO developers(email, name, added_at, removed_at) VALUES(?, ?, ?, NULL)`,
			email, name, addedAt.UTC(),
		)
		if e != nil {
			return fmt.Errorf("add developer: %w", e)
		}
		devID, _ = res.LastInsertId()
	} else if err != nil {
		return fmt.Errorf("add developer: %w", err)
	} else {
		if name != "" {
			if _, e := tx.Exec(`UPDATE developers SET name = ?, removed_at = NULL WHERE id = ?`, name, devID); e != nil {
				return fmt.Errorf("add developer: %w", e)
			}
		} else {
			if _, e := tx.Exec(`UPDATE developers SET removed_at = NULL WHERE id = ?`, devID); e != nil {
				return fmt.Errorf("add developer: %w", e)
			}
		}
	}

	var dummy int
	err = tx.QueryRow(
		`SELECT 1 FROM developer_status WHERE developer_id = ? AND removed_at IS NULL LIMIT 1`,
		devID,
	).Scan(&dummy)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("add developer: %w", err)
	}
	if err != sql.ErrNoRows {
		return fmt.Errorf("developer already active: %s", email)
	}

	if _, err := tx.Exec(
		`INSERT INTO developer_status(developer_id, added_at, removed_at) VALUES(?, ?, NULL)`,
		devID, addedAt.UTC(),
	); err != nil {
		return fmt.Errorf("add developer: %w", err)
	}

	return tx.Commit()
}

func RemoveDeveloper(email string, removedAt time.Time) error {
	email = normalizeEmail(email)
	if email == "" {
		return fmt.Errorf("developer email is required")
	}

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

	var devID int64
	err = tx.QueryRow(`SELECT id FROM developers WHERE lower(email) = lower(?)`, email).Scan(&devID)
	if err == sql.ErrNoRows {
		return fmt.Errorf("developer not found: %s", email)
	}
	if err != nil {
		return fmt.Errorf("remove developer: %w", err)
	}

	res, err := tx.Exec(
		`UPDATE developer_status
         SET removed_at = ?
         WHERE developer_id = ? AND removed_at IS NULL`,
		removedAt.UTC(), devID,
	)
	if err != nil {
		return fmt.Errorf("remove developer: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("developer not active: %s", email)
	}

	if _, err := tx.Exec(
		`UPDATE developers SET removed_at = ? WHERE id = ?`,
		removedAt.UTC(), devID,
	); err != nil {
		return fmt.Errorf("remove developer: %w", err)
	}

	return tx.Commit()
}

func GetDeveloperByEmail(email string) (*Developer, error) {
	email = normalizeEmail(email)
	if email == "" {
		return nil, nil
	}

	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var d Developer
	err = db.QueryRow(
		`SELECT id, email, name, added_at, removed_at FROM developers WHERE lower(email) = lower(?)`,
		email,
	).Scan(&d.ID, &d.Email, &d.Name, &d.AddedAt, &d.RemovedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func IsDeveloperActiveAt(email string, t time.Time) (bool, error) {
	email = normalizeEmail(email)
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
		`SELECT 1 FROM developer_status
         WHERE developer_id = ?
           AND added_at <= ?
           AND (removed_at IS NULL OR removed_at > ?)
         ORDER BY added_at DESC
         LIMIT 1`,
		d.ID, t.UTC(), t.UTC(),
	).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func ListDevelopers() ([]Developer, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`
SELECT d.id, d.email, d.name, d.added_at,
       (SELECT removed_at FROM developer_status s
        WHERE s.developer_id = d.id AND s.removed_at IS NOT NULL
        ORDER BY s.removed_at DESC LIMIT 1) AS last_removed_at,
       EXISTS(SELECT 1 FROM developer_status s
              WHERE s.developer_id = d.id AND s.removed_at IS NULL) AS is_active
FROM developers d
ORDER BY d.email`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Developer
	for rows.Next() {
		var d Developer
		var lastRemoved sql.NullTime
		var isActive bool
		if err := rows.Scan(&d.ID, &d.Email, &d.Name, &d.AddedAt, &lastRemoved, &isActive); err != nil {
			return nil, err
		}
		if !isActive {
			d.RemovedAt = lastRemoved
		} else {
			d.RemovedAt = sql.NullTime{Valid: false}
		}
		out = append(out, d)
	}
	return out, nil
}

func AddDeveloperStatus(email string, addedAt time.Time, removedAt sql.NullTime) error {
	email = normalizeEmail(email)
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
		`INSERT INTO developer_status(developer_id, added_at, removed_at) VALUES(?, ?, ?)`,
		d.ID, addedAt.UTC(), removedAt,
	)
	if err != nil {
		return fmt.Errorf("add developer status: %w", err)
	}
	return nil
}

func HasActiveDeveloperStatus(email string) (bool, error) {
	email = normalizeEmail(email)
	d, err := GetDeveloperByEmail(email)
	if err != nil {
		return false, err
	}
	if d == nil {
		return false, fmt.Errorf("developer not found: %s", email)
	}

	db, err := openDB()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var dummy int
	err = db.QueryRow(
		`SELECT 1 FROM developer_status WHERE developer_id = ? AND removed_at IS NULL LIMIT 1`,
		d.ID,
	).Scan(&dummy)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
