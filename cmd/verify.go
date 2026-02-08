package cmd

import (
	"fmt"
	"time"

	"github.com/esuarkeN/valiDTr/db"
	"github.com/esuarkeN/valiDTr/git"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:     "verify [commit-hash]",
	Short:   "Verify one commit against signature + developer/key policy",
	Args:    cobra.ExactArgs(1),
	PreRunE: initDB,
	RunE: func(cmd *cobra.Command, args []string) error {
		commit := args[0]

		if err := git.VerifyCommitSignature(commit); err != nil {
			return fmt.Errorf("%s: %w", commit, err)
		}

		keyID, err := git.CommitKeyID(commit)
		if err != nil {
			return fmt.Errorf("%s: %w", commit, err)
		}

		email, err := git.CommitEmail(commit, EmailMode())
		if err != nil {
			return fmt.Errorf("%s: %w", commit, err)
		}

		var t time.Time
		switch Policy() {
		case "current":
			t = time.Now().UTC()
		case "historical":
			t, err = git.CommitTimestamp(commit)
			if err != nil {
				return fmt.Errorf("%s: %w", commit, err)
			}
		default:
			return fmt.Errorf("unknown policy: %s (use current|historical)", Policy())
		}

		activeDev, err := db.IsDeveloperActiveAt(email, t)
		if err != nil {
			return err
		}
		if !activeDev {
			return fmt.Errorf("%s: developer not allowed (%s) under policy=%s", commit, email, Policy())
		}

		activeKey, err := db.IsKeyActiveForDeveloperAt(email, keyID, t)
		if err != nil {
			return err
		}
		if !activeKey {
			return fmt.Errorf("%s: key not allowed for %s (key=%s) under policy=%s", commit, email, keyID, Policy())
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
