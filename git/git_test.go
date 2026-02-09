package git_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/esuarkeN/valiDTr/git"
)

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not found in PATH")
	}
}

func run(t *testing.T, env []string, args ...string) string {
	t.Helper()

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), env...)

	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		t.Fatalf("%s failed: %v\nstderr: %s", strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
	}
	return strings.TrimSpace(out.String())
}

func TestGitHelpers_TimestampEmailRangeAndSignatureFailure(t *testing.T) {
	requireGit(t)

	repo := t.TempDir()
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd: %v", err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatalf("Chdir(%q): %v", repo, err)
	}
	t.Cleanup(func() { _ = os.Chdir(oldwd) })

	run(t, nil, "git", "init", "-b", "main")
	run(t, nil, "git", "config", "user.name", "Test User")
	run(t, nil, "git", "config", "user.email", "test@example.com")
	run(t, nil, "git", "config", "commit.gpgsign", "false")

	writeFile := func(name, content string) {
		t.Helper()
		p := filepath.Join(repo, name)
		if err := os.WriteFile(p, []byte(content), 0o600); err != nil {
			t.Fatalf("WriteFile(%q): %v", p, err)
		}
	}

	commit := func(msg, ts string, extra string) string {
		t.Helper()
		writeFile("file.txt", msg+extra+"\n")
		run(t, nil, "git", "add", "file.txt")
		run(t, []string{
			"GIT_AUTHOR_DATE=" + ts,
			"GIT_COMMITTER_DATE=" + ts,
		}, "git", "commit", "-m", msg)
		return run(t, nil, "git", "rev-parse", "HEAD")
	}

	c1 := commit("c1", "2020-01-01T00:00:00Z", "")
	c2 := commit("c2", "2020-01-02T00:00:00Z", "x")

	root, err := git.RootCommit("HEAD")
	if err != nil {
		t.Fatalf("RootCommit: %v", err)
	}
	if root != c1 {
		t.Fatalf("expected root %s, got %s", c1, root)
	}

	email, err := git.CommitEmail(c1, "committer")
	if err != nil {
		t.Fatalf("CommitEmail: %v", err)
	}
	if email != "test@example.com" {
		t.Fatalf("expected committer email test@example.com, got %q", email)
	}
	if _, err := git.CommitEmail(c1, "nope"); err == nil {
		t.Fatalf("expected invalid email mode to fail")
	}

	ts, err := git.CommitTimestamp(c1)
	if err != nil {
		t.Fatalf("CommitTimestamp: %v", err)
	}
	wantTS := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	if !ts.Equal(wantTS) {
		t.Fatalf("expected timestamp %s, got %s", wantTS.Format(time.RFC3339), ts.Format(time.RFC3339))
	}

	commits, err := git.CommitsInRange(c1, c2)
	if err != nil {
		t.Fatalf("CommitsInRange: %v", err)
	}
	if len(commits) != 1 || commits[0] != c2 {
		t.Fatalf("expected range commits [%s], got %v", c2, commits)
	}

	if err := git.VerifyCommitSignature(c1); err == nil {
		t.Fatalf("expected unsigned commit to fail signature verification")
	}
	if _, err := git.CommitKeyID(c1); err == nil {
		t.Fatalf("expected CommitKeyID to fail for unsigned commit")
	}
}
