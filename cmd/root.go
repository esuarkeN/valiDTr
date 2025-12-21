package cmd

import (
	"os"

	"valiDTr/db"

	"github.com/spf13/cobra"
)

var (
	dbPath    string
	emailMode string
	policy    string
)

func envOrDefault(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

var rootCmd = &cobra.Command{
	Use:   "valiDTr",
	Short: "Verify Git commits signed with GPG against a time-aware developer/key allowlist",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return db.InitDB(dbPath)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func EmailMode() string { return emailMode }
func Policy() string    { return policy }
func init() {
	rootCmd.PersistentFlags().StringVar(
		&dbPath,
		"db",
		envOrDefault("VALIDTR_DB", "validtr.db"),
		"Path to SQLite DB (can be temp in CI)",
	)
	rootCmd.PersistentFlags().StringVar(
		&emailMode,
		"email-mode",
		envOrDefault("VALIDTR_EMAIL_MODE", "committer"),
		"Email source: committer|author",
	)
	rootCmd.PersistentFlags().StringVar(
		&policy,
		"policy",
		envOrDefault("VALIDTR_POLICY", "current"),
		"Verification policy: current|historical",
	)
}
