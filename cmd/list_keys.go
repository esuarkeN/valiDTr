package cmd

import (
	"fmt"

	"valiDTr/db"

	"github.com/spf13/cobra"
)

var listKeysEmail string

var listKeysCmd = &cobra.Command{
	Use:   "list-keys",
	Short: "List developer key mappings",
	RunE: func(cmd *cobra.Command, args []string) error {
		keys, err := db.ListKeys(listKeysEmail)
		if err != nil {
			return err
		}
		for _, k := range keys {
			status := "active"
			if k.RevokedAt.Valid {
				status = "revoked@" + k.RevokedAt.Time.UTC().Format("2006-01-02T15:04:05Z")
			}
			fmt.Printf("%s\t%s\t%s\tadded@%s\n",
				k.DeveloperEmail, k.KeyID, status, k.AddedAt.UTC().Format("2006-01-02T15:04:05Z"))
		}
		return nil
	},
}

func init() {
	listKeysCmd.Flags().StringVar(&listKeysEmail, "email", "", "Filter by developer email")
	rootCmd.AddCommand(listKeysCmd)
}
