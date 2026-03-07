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
  evolved-commit explain commit-message-subject-not-empty`,
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
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown rule '%s'. Use 'evolved-commit explain --help' for available rules.\n", ruleName)
			exit(1)
	}
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
