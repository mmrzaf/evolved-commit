package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// explainCmd represents the explain command
var explainCmd = &cobra.Command{
	Use:   "explain <rule>",
	Short: "Explains a specific evolved-commit rule and how to satisfy it",
	Long: `The explain command provides detailed information about a specific
evolved-commit rule, including its purpose, why it's important, and
exact steps on how to fix a violation.

Usage:
  evolved-commit explain commit-message-subject-not-empty
  evolved-commit explain commit-message-subject-length`,
	Args: cobra.ExactArgs(1), // Expect exactly one argument: the rule name
	Run: func(cmd *cobra.Command, args []string) {
		explainCommandLogic(cmd, args, os.Exit)
	},
}

// explainCommandLogic encapsulates the core logic of the explain command, allowing os.Exit to be mocked for testing.
func explainCommandLogic(cmd *cobra.Command, args []string, exit func(code int)) {
	ruleName := strings.ToLower(args[0])

	switch ruleName {
		case "commit-message-subject-not-empty":
			fmt.Println(`Rule: Commit Message Subject Not Empty

Purpose:
A commit message subject provides a concise summary of the change.
It's crucial for quick understanding when browsing history, generating
changelogs, or reviewing pull requests.

Why it's important:
An empty subject makes it difficult to understand the commit's intent
without diving into the full diff or body. It hinders traceability and
code archaeology.

How to fix:
Ensure the first line of your commit message is not blank.
Example:
  git commit -m "feat: Add user authentication"
  git commit
  (Then type your subject line in the editor)`) 
			exit(0) // Explicitly exit with 0 on success
		case "commit-message-subject-length":
			fmt.Println(`Rule: Commit Message Subject Length

Purpose:
A commit message subject should be concise and ideally fit on a single line
when viewing Git logs or reviewing changes. A common convention is to keep
it under 50 characters.

Why it's important:
Short subject lines are easier to read and scan quickly. They encourage
authors to summarize the most important aspect of their commit, which improves
the readability of commit history. Many Git tools truncate longer subjects.

How to fix:
Condense the first line of your commit message to be 50 characters or less.
If more detail is needed, use the commit message body (second paragraph onwards).
Example:
  git commit -m "feat: Implement user authentication and permissions"
  (This is 50 characters)
  git commit -m "refactor: Optimize data fetching"
  (This is 30 characters)`) 
			exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown rule '%s'. Use 'evolved-commit explain --help' for available rules.\n", ruleName)
			exit(1)
	}
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
