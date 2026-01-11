package cmd

import (
	"github.com/spf13/cobra"
)


var checkCmdAgainstDb = &cobra.Command{
	Use: "check commit against db",
	Short: "check commit signature against db",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
	},
}