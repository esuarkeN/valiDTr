package cmd

import (
	"fmt"
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

		err := git.VerifyCommitSignature(commitHash)
		if err != nil {
			fmt.Println("Signature invalid:", err)
			return
		}

		keyID, err := git.GetCommitKeyID(commitHash)
		if err != nil {
			fmt.Println("Error extracting key ID:", err)
			return
		}

		isValid, err := db.IsKeyActive(keyID)
		if err != nil {
			fmt.Println("Database error:", err)
			return
		}
		if !isValid {
			fmt.Println("Key is not trusted or has been revoked.")
			return
		}

		fmt.Println("Commit is verified and key is trusted.")
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
