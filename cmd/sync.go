package cmd

import (
	"fmt"

	"valiDTr/db"

	"github.com/spf13/cobra"
)

var syncConfigPath string
var syncReset bool
var syncMode string // "reset" | "reconcile"

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync DB from YAML config (reset or reconcile)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if syncConfigPath == "" {
			return fmt.Errorf("--config is required")
		}

		switch syncMode {
		case "reset":
			return db.SyncFromYAMLReset(syncConfigPath, syncReset)
		case "reconcile":
			// reconcile always preserves history (no deletes)
			return db.SyncFromYAMLReconcile(syncConfigPath)
		default:
			return fmt.Errorf("unknown --mode=%q (use reset|reconcile)", syncMode)
		}
	},
}

func init() {
	syncCmd.Flags().StringVar(&syncConfigPath, "config", "", "Path to .validtr/config.yml")
	syncCmd.Flags().BoolVar(&syncReset, "reset", true, "Reset tables before sync (mode=reset only)")
	syncCmd.Flags().StringVar(&syncMode, "mode", "reset", "Sync mode: reset|reconcile")
	rootCmd.AddCommand(syncCmd)
}
