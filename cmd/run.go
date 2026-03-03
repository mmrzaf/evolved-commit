package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs evolved-commit checks",
	Long: `The run command executes all configured evolved-commit checks on your
commit message, branch name, and staged diffs. This command is typically
called by Git hooks but can also be run manually or in CI/CD pipelines.`,
	Run: func(cmd *cobra.Command, args []string) {
		runCommandLogic(cmd, args, os.Exit)
	},
}

// runCommandLogic encapsulates the core logic of the run command, allowing os.Exit to be mocked for testing.
func runCommandLogic(cmd *cobra.Command, args []string, exit func(code int)) {
	fmt.Println("evolved-commit run: Running checks (not yet implemented).")
	// TODO: Implement actual checks here.
	// For now, exit with success to not block commits.
	// In the future, this will return an error and os.Exit(1) if checks fail.
	exit(0)
}

func init() {
	rootCmd.AddCommand(runCmd)
}
