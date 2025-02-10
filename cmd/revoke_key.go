package cmd

import (
	"fmt"
	"valiDTr/db"

	"github.com/spf13/cobra"
)

var revokeKeyCmd = &cobra.Command{
	Use:   "revoke-key [key-id]",
	Short: "Revoke a developer's GPG key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyID := args[0]
		err := db.RevokeGPGKey(keyID)
		if err != nil {
			fmt.Println("Error revoking key:", err)
		} else {
			fmt.Println("Key revoked successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(revokeKeyCmd)
}
