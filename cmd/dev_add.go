package cmd

import (
	"github.com/esuarkeN/valiDTr/db"

	"github.com/spf13/cobra"
)

var devAddName string
var devAddSince string

var devAddCmd = &cobra.Command{
	Use:     "add-dev [email]",
	Short:   "Add a developer (active starting at --since)",
	Args:    cobra.ExactArgs(1),
	PreRunE: initDB,
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		since, err := parseSince(devAddSince)
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
