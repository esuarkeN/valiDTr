package cmd 
import (
	"valiDTr/cmd"
	"valiDTr/db"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "valiDTr",
	Short: "valiDTr is a CLI tool for verifying Git commit chains using GPG signatures",
}

func main() {
	db.InitDB()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
