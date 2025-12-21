package cmd

import (
	"time"

	"valiDTr/db"

	"github.com/spf13/cobra"
)

var addKeyCmd = &cobra.Command{
	Use:   "add-key [developer-email] [key-id]",
	Short: "Assign a GPG key to a developer (active starting now)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		keyID := args[1]
		return db.AddKeyToDeveloper(email, keyID, time.Now().UTC())
	},
}

func init() {
	rootCmd.AddCommand(addKeyCmd)
}
