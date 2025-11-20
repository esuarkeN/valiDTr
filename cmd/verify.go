package cmd

import (
	"fmt"
	"time"
	"valiDTr/db"
	"valiDTr/git"

	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [commit-hash]",
	Short: "Verify the GPG signature of a Git commit",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		commitHash := args[0]

		if err := git.VerifyCommitSignature(commitHash); err != nil {
			fmt.Println("Signature invalid:", err)
			return
		}

		keyID, err := git.GetCommitKeyID(commitHash)
		if err != nil {
			fmt.Println("Error extracting key ID:", err)
			return
		}

		commitTime, err := git.GetCommitTimestamp(commitHash)
		if err != nil {
			fmt.Println("Error getting commit timestamp:", err)
			return
		}

		trusted, err := db.IsKeyActiveAt(keyID, commitTime)
		if err != nil {
			fmt.Println("Database error:", err)
			return
		}
		if !trusted {
			fmt.Printf("Key %s was not trusted at commit time (%s).\n", keyID, commitTime.Format(time.RFC3339))
			return
		}

		fmt.Println("Commit is verified and key was trusted at commit time.")
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
