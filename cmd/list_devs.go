package cmd

import (
	"fmt"

	"github.com/esuarkeN/valiDTr/db"

	"github.com/spf13/cobra"
)

var listDevsCmd = &cobra.Command{
	Use:     "list-devs",
	Short:   "List developers",
	PreRunE: initDB,
	RunE: func(cmd *cobra.Command, args []string) error {
		devs, err := db.ListDevelopers()
		if err != nil {
			return err
		}
		for _, d := range devs {
			name := ""
			if d.Name.Valid {
				name = d.Name.String
			}
			status := "active"
			if d.RemovedAt.Valid {
				status = "removed@" + d.RemovedAt.Time.UTC().Format("2006-01-02T15:04:05Z")
			}
			fmt.Printf("%s\t%s\t%s\tadded@%s\n",
				d.Email, name, status, d.AddedAt.UTC().Format("2006-01-02T15:04:05Z"))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listDevsCmd)
}
