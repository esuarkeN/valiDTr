package cmd

import (
	"fmt"
	"valiDTr/db"

	"github.com/spf13/cobra"
)

var listKeysCmd = &cobra.Command{
	Use:   "list-keys",
	Short: "List all stored GPG keys in the database",
	Run: func(cmd *cobra.Command, args []string) {
		keys, err := db.ListGPGKeys()
		if err != nil {
			fmt.Println("Error retrieving keys:", err)
			return
		}

		fmt.Println("Stored GPG Keys:")
		fmt.Println("----------------------------")
		for _, key := range keys {
			fmt.Printf("Key ID: %s | Status: %s\n", key.KeyID, key.Status)
		}
	},
}

func init() {
	rootCmd.AddCommand(listKeysCmd)
}
