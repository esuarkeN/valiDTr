package cmd

import (
	"time"

	"valiDTr/db"

	"github.com/spf13/cobra"
)

var devAddName string
var devAddSince string

func parseDevSince(s string) (time.Time, error) {
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

var devAddCmd = &cobra.Command{
	Use:   "add-dev [email]",
	Short: "Add a developer (active starting at --since)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		since, err := parseDevSince(devAddSince)
		if err != nil {
			return err
		}
		return db.AddDeveloper(email, devAddName, since)
	},
}

func init() {
	devAddCmd.Flags().StringVar(&devAddName, "name", "", "Optional display name")
	devAddCmd.Flags().StringVar(&devAddSince, "since", "now", "Start time: now|epoch|RFC3339")
	rootCmd.AddCommand(devAddCmd)
}
