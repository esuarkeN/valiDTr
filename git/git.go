package git

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
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
func GetCommitKeyID(commitHash string) (string, error) {
	cmd := exec.Command("git", "log", "--show-signature", "-1", commitHash)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output := out.String()
	// Look for "using RSA key <KEYID>"
	re := regexp.MustCompile(`using \w+ key ([A-F0-9]+)`)
	match := re.FindStringSubmatch(output)
	if len(match) < 2 {
		return "", errors.New("could not extract key ID from signature")
	}
	return strings.TrimSpace(match[1]), nil
}
func GetCommitTimestamp(commitHash string) (time.Time, error) {
	cmd := exec.Command("git", "show", "-s", "--format=%ct", commitHash)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return time.Time{}, err
	}

	s := strings.TrimSpace(out.String())
	secs, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(secs, 0).UTC(), nil
}
