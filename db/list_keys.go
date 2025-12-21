package db

func ListDevKeys(emailFilter string) ([]DevKey, error) {
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
	var args []any
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
