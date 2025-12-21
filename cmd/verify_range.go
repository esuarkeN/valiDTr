package cmd

import (
	"strings"

	"valiDTr/git"

	"github.com/spf13/cobra"
)

var verifyRangeCmd = &cobra.Command{
	Use:   "verify-range [from] [to]",
	Short: "Verify all commits in a range (from..to), failing on the first rejection",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		from := args[0]
		to := args[1]

		if strings.Trim(from, "0") == "" {
			root, err := git.RootCommit(to)
			if err != nil {
				return err
			}
			from = root
		}

		commits, err := git.CommitsInRange(from, to)
		if err != nil {
			return err
		}

		for _, c := range commits {
			if err := verifyCmd.RunE(cmd, []string{c}); err != nil {
				return err
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyRangeCmd)
}
