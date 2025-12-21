package cmd

import (
	"fmt"

	"valiDTr/db"

	"github.com/spf13/cobra"
)

var syncConfigPath string
var syncReset bool

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync DB from a YAML config (best for CI/pipelines)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if syncConfigPath == "" {
			return fmt.Errorf("--config is required")
		}
		return db.SyncFromYAML(syncConfigPath, syncReset)
	},
}

func init() {
	syncCmd.Flags().StringVar(&syncConfigPath, "config", "", "Path to .validtr/config.yml")
	syncCmd.Flags().BoolVar(&syncReset, "reset", true, "Reset tables before sync")
	rootCmd.AddCommand(syncCmd)
}
