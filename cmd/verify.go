package cmd

import (
	"fmt"
	"valiDTr/git"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [commit-hash]",
	Short: "Verify the GPG signature of a Git commit",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitHash := args[0]
		err := git.VerifyCommitSignature(commitHash)
		if err != nil {
			fmt.Println("Verification failed:", err)
		} else {
			fmt.Println("Commit verified successfully!")
		}
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
