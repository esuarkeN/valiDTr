package cmd

import (
	"time"

	"github.com/esuarkeN/valiDTr/db"

	"github.com/spf13/cobra"
)

var devRemoveCmd = &cobra.Command{
	Use:     "remove-dev [email]",
	Short:   "Remove a developer (inactive starting now)",
	Args:    cobra.ExactArgs(1),
	PreRunE: initDB,
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		return db.RemoveDeveloper(email, time.Now().UTC())
	},
}

func init() {
	rootCmd.AddCommand(devRemoveCmd)
}
