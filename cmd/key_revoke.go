package cmd

import (
	"time"

	"github.com/esuarkeN/valiDTr/db"

	"github.com/spf13/cobra"
)

var keyRevokeCmd = &cobra.Command{
	Use:     "revoke-key [developer-email] [key-id]",
	Short:   "Revoke a developer's key (inactive starting now)",
	Args:    cobra.ExactArgs(2),
	PreRunE: initDB,
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		keyID := args[1]
		return db.RevokeDeveloperKey(email, keyID, time.Now().UTC())
	},
}

func init() {
	rootCmd.AddCommand(keyRevokeCmd)
}
