package gpg

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

func VerifyGPGSignature(commitHash string) error {
	cmd := exec.Command("git", "show", "-s", "--format=%G?", commitHash)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return err
	}

	status := strings.TrimSpace(out.String())
	if status != "G" {
		return fmt.Errorf("signature not valid for %s (%%G?=%s)", commitHash, status)
	}
	return nil
}
