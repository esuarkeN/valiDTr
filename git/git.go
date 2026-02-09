package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runGit(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("git %v failed: %w: %s", args, err, strings.TrimSpace(errb.String()))
	}
	return strings.TrimSpace(out.String()), nil
}
func SignatureStatus(commit string) (string, error) {
	return runGit("show", "-s", "--format=%G?", commit)
}

func CommitKeyID(commit string) (string, error) {
	key, err := runGit("show", "-s", "--format=%GK", commit)
	if err != nil {
		return "", err
	}
	if key == "" {
		return "", fmt.Errorf("no signing key found (unsigned commit)")
	}
	return strings.ToUpper(key), nil
}

func CommitTimestamp(commit string) (time.Time, error) {
	s, err := runGit("show", "-s", "--format=%ct", commit)
	if err != nil {
		return time.Time{}, err
	}
	secs, err := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid commit timestamp %q: %w", s, err)
	}
	return time.Unix(secs, 0).UTC(), nil
}

func CommitEmail(commit, mode string) (string, error) {
	var format string
	switch mode {
	case "committer":
		format = "%ce"
	case "author":
		format = "%ae"
	default:
		return "", fmt.Errorf("unknown email mode: %s (use committer|author)", mode)
	}
	email, err := runGit("show", "-s", "--format="+format, commit)
	if err != nil {
		return "", err
	}
	if email == "" {
		return "", fmt.Errorf("could not read %s email", mode)
	}
	return strings.ToLower(email), nil
}

func VerifyCommitSignature(commit string) error {
	st, err := SignatureStatus(commit)
	if err != nil {
		return err
	}
	// Accept good signatures even if the key isn't trusted in the local keyring.
	if st != "G" && st != "U" {
		return fmt.Errorf("signature not valid: %%G?=%s", st)
	}
	return nil
}

func CommitsInRange(from, to string) ([]string, error) {
	out, err := runGit("rev-list", "--reverse", from+".."+to)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(out) == "" {
		return []string{}, nil
	}
	return strings.Split(out, "\n"), nil
}

func RootCommit(of string) (string, error) {
	return runGit("rev-list", "--max-parents=0", of)
}
