package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "valiDTr",
	Short: "A CLI tool to verify Git commit chains using GPG",
}

func Execute() error {
	return rootCmd.Execute()
}
