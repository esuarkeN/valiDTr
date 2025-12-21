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
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(
		`INSERT INTO developers(email, name, added_at) VALUES(?, ?, ?)`,
		email, name, addedAt.UTC(),
	)
	if err != nil {
		return fmt.Errorf("add developer: %w", err)
	}
	return nil
}

func RemoveDeveloper(email string, removedAt time.Time) error {
	db, err := openDB()
	if err != nil {
		return err
	}
	defer db.Close()

	res, err := db.Exec(
		`UPDATE developers SET removed_at = ? WHERE email = ? AND removed_at IS NULL`,
		removedAt.UTC(), email,
	)
	if err != nil {
		return fmt.Errorf("remove developer: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("developer not found or already removed: %s", email)
	}
	return nil
}

func GetDeveloperByEmail(email string) (*Developer, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var d Developer
	err = db.QueryRow(
		`SELECT id, email, name, added_at, removed_at FROM developers WHERE email = ?`,
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
	d, err := GetDeveloperByEmail(email)
	if err != nil {
		return false, err
	}
	if d == nil {
		return false, nil
	}

	if t.Before(d.AddedAt) {
		return false, nil
	}
	if d.RemovedAt.Valid && !t.Before(d.RemovedAt.Time) {
		return false, nil
	}
	return true, nil
}

func ListDevelopers() ([]Developer, error) {
	db, err := openDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, email, name, added_at, removed_at FROM developers ORDER BY email`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Developer
	for rows.Next() {
		var d Developer
		if err := rows.Scan(&d.ID, &d.Email, &d.Name, &d.AddedAt, &d.RemovedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}
