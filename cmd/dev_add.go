package cmd

import (
	"time"

	"valiDTr/db"

	"github.com/spf13/cobra"
)

var devAddName string

var devAddCmd = &cobra.Command{
	Use:   "add-dev [email]",
	Short: "Add a developer (active starting now)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		return db.AddDeveloper(email, devAddName, time.Now().UTC())
	},
}

func init() {
	devAddCmd.Flags().StringVar(&devAddName, "name", "", "Optional display name")
	rootCmd.AddCommand(devAddCmd)
}
