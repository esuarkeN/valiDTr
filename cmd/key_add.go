package cmd

import (
	"time"

	"valiDTr/db"

	"github.com/spf13/cobra"
)

var keyAddSince string

func parseSince(s string) (time.Time, error) {
	if s == "" || s == "now" {
		return time.Now().UTC(), nil
	}
	if s == "epoch" {
		return time.Unix(0, 0).UTC(), nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, err
	}
	return t.UTC(), nil
}

var keyAddCmd = &cobra.Command{
	Use:   "add-key [developer-email] [key-id]",
	Short: "Assign a GPG key to a developer (active starting at --since)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		keyID := args[1]
		since, err := parseSince(keyAddSince)
		if err != nil {
			return err
		}
		return db.AddKeyToDeveloper(email, keyID, since)
	},
}

func init() {
	keyAddCmd.Flags().StringVar(&keyAddSince, "since", "now", "Start time: now|epoch|RFC3339")
	rootCmd.AddCommand(keyAddCmd)
}
