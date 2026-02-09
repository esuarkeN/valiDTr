package db

import "strings"

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

func normalizeKeyID(keyID string) string {
	return strings.ToUpper(strings.TrimSpace(keyID))
}
