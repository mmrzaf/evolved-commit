package cmd

import (
	"evolved-commit/pkg/checks"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run [COMMIT_MSG_FILE]",
	Short: "Runs evolved-commit checks",
	Long: `The run command executes all configured evolved-commit checks on your
commit message, branch name, and staged diffs. This command is typically
called by Git hooks but can also be run manually or in CI/CD pipelines.

When run as a Git hook (e.g., pre-commit), it receives the path to the
commit message file as an argument.`,
	Args: cobra.MaximumNArgs(1), // Allow 0 or 1 argument (the commit message file path)
	Run: func(cmd *cobra.Command, args []string) {
		runCommandLogic(cmd, args, os.Exit)
	},
}

// runCommandLogic encapsulates the core logic of the run command, allowing os.Exit to be mocked for testing.
func runCommandLogic(cmd *cobra.Command, args []string, exit func(code int)) {
	if len(args) > 0 {
		commitMessageFilePath := args[0]
		content, err := os.ReadFile(commitMessageFilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to read commit message file %s: %v\n", commitMessageFilePath, err)
			exit(1)
			return
		}

		commitMessage := string(content)
		subjectLine := strings.SplitN(commitMessage, "\n", 2)[0] // Get the first line

		var failures []error // Collect all failures

		// Run commit message subject checks
		if err := checks.CheckCommitMessageSubjectNotEmpty(subjectLine); err != nil {
			failures = append(failures, err)
		}
		if err := checks.CheckCommitMessageSubjectLength(subjectLine); err != nil {
			failures = append(failures, err)
		}
		if err := checks.CheckCommitMessageSubjectNoTrailingPeriod(subjectLine); err != nil {
			failures = append(failures, err)
		}
		if err := checks.CheckCommitMessageSubjectStartsWithUppercase(subjectLine); err != nil {
			failures = append(failures, err)
		}
		if err := checks.CheckCommitMessageSubjectImperative(subjectLine); err != nil {
			failures = append(failures, err)
		}

		// Report all collected failures or exit successfully
		if len(failures) > 0 {
			fmt.Fprintln(os.Stderr, "Commit message checks failed:")
			for _, failure := range failures {
				fmt.Fprintf(os.Stderr, "- %v\n", failure)
			}
			exit(1)
			return
		}

		exit(0) // All checks passed
		return
	}

	fmt.Println("evolved-commit run: No commit message file provided. Running general checks (not yet implemented).")
	// TODO: Implement actual general checks here or adapt to read from stdin/config if no file is provided.
	// For now, exit with success to not block commits when called manually without arguments.
	exit(0)
}

func init() {
	rootCmd.AddCommand(runCmd)
}
