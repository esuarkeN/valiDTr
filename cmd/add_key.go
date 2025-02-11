package cmd

import (
	"fmt"
	"valiDTr/db"

	"github.com/spf13/cobra"
)

var addKeyCmd = &cobra.Command{
	Use:   "add-key [key-id]",
	Short: "Add a GPG key to the verification database",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyID := args[0]
		err := db.AddGPGKey(keyID)
		if err != nil {
			fmt.Println("Error adding key:", err)
		} else {
			fmt.Println("Key added successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(addKeyCmd)
}
