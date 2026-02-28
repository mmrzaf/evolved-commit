package cmd

import (
	"evolved-commit/pkg/githooks"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Removes Git hooks installed by evolved-commit",
	Long: `The uninstall command removes any Git hooks previously installed by
evolved-commit from your repository, disabling its automatic checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Attempting to uninstall Git hooks...")
		err := githooks.Uninstall("") // Pass empty string to let githooks find the repo
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error uninstalling Git hooks: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Git hooks uninstalled successfully.")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
