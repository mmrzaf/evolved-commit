package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// uninstallCmd represents the uninstall command
var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Removes Git hooks installed by evolved-commit",
	Long: `The uninstall command removes any Git hooks previously installed by
evolved-commit from your repository, disabling its automatic checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("uninstall called: This will remove Git hooks.")
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}
