package git

import (
	"fmt"
	"os/exec"
)

func VerifyCommitSignature(commitHash string) error {
	cmd := exec.Command("git", "log", "--show-signature", "-1", commitHash)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
