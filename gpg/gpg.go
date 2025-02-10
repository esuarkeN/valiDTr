package gpg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CheckGPGKeyExists checks if a given GPG key is in the keyring
func CheckGPGKeyExists(keyID string) (bool, error) {
	cmd := exec.Command("gpg", "--list-keys", "--with-colons", keyID)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false, err
	}

	// If the output contains "pub", the key exists
	return strings.Contains(out.String(), "pub"), nil
}

// GetGPGKeyFingerprint retrieves the fingerprint of a GPG key
func GetGPGKeyFingerprint(keyID string) (string, error) {
	cmd := exec.Command("gpg", "--fingerprint", keyID)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// Extract the fingerprint from output
	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "Key fingerprint") {
			return strings.TrimSpace(strings.Split(line, "=")[1]), nil
		}
	}

	return "", fmt.Errorf("fingerprint not found for key %s", keyID)
}

// VerifyGPGSignature verifies if a commit is correctly signed
func VerifyGPGSignature(commitHash string) error {
	cmd := exec.Command("git", "log", "--show-signature", "-1", commitHash)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println(out.String())
	if strings.Contains(out.String(), "Good signature") {
		return nil
	}

	return fmt.Errorf("invalid or missing signature for commit %s", commitHash)
}
