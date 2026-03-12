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
  evolved-commit explain commit-message-subject-length
  evolved-commit explain commit-message-subject-no-trailing-period
  evolved-commit explain commit-message-subject-starts-with-uppercase
  evolved-commit explain commit-message-subject-imperative`,
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
		case "commit-message-subject-no-trailing-period":
			fmt.Println(`Rule: Commit Message Subject No Trailing Period

Purpose:
To keep commit subjects clean, concise, and consistent, they should not
end with a period.

Why it's important:
Periods at the end of subject lines are often considered superfluous in
this context, as the subject is a title, not a sentence. Omitting them
contributes to a consistent and scan-friendly commit history.

How to fix:
Remove any trailing period from the first line of your commit message.
Example:
  git commit -m "fix: Corrected a critical bug in user authentication"
  (This is incorrect because of the trailing period)
  git commit -m "fix: Correct a critical bug in user authentication"
  (This is correct)`) 
			exit(0)
		case "commit-message-subject-starts-with-uppercase":
			fmt.Println(`Rule: Commit Message Subject Starts With Uppercase

Purpose:
A commit message subject should start with an uppercase letter to ensure
consistency and readability across the commit history.

Why it's important:
Following a consistent capitalization style for commit subjects makes the
commit history easier to scan and understand. It's a common convention that
improves the professional appearance and navigability of version control logs.

How to fix:
Ensure the first significant character of your commit message subject line is
an uppercase letter. This applies even if your subject starts with a type prefix
like 'feat:' or 'fix:'.
Example:
  git commit -m "feat: add new user registration flow"
  (This is incorrect because 'add' starts with a lowercase 'a')
  git commit -m "Feat: Add new user registration flow"
  (This is correct)
  git commit -m "Fix: Correct critical bug"
  (This is correct)`) 
			exit(0)
		case "commit-message-subject-imperative":
			fmt.Println(`Rule: Commit Message Subject Imperative Mood

Purpose:
A commit message subject should describe the change in the imperative mood,
as if giving a command or instruction. This means it should start with an
action verb in its base form.

Why it's important:
Using the imperative mood ("Fix bug", "Add feature") makes commit history
consistent, easier to read, and aligns with how Git itself phrases merge
commits ("Merge branch 'feature'"). It clarifies what the commit *does*,
rather than what it *did* or *does*.

How to fix:
Rephrase the first word of your commit message subject to be an imperative
verb. Avoid past tense (-ed), third-person singular (-s), or present participle (-ing).

Example:
  git commit -m "feat: adds user authentication" (Incorrect: "adds")
  git commit -m "feat: Added user authentication" (Incorrect: "Added")
  git commit -m "feat: Adding user authentication" (Incorrect: "Adding")
  git commit -m "feat: Add user authentication" (Correct)

  git commit -m "fix: corrected a critical bug" (Incorrect: "corrected")
  git commit -m "fix: Correct a critical bug" (Correct)`) 
			exit(0)
		default:
			fmt.Fprintf(os.Stderr, "Error: Unknown rule '%s'. Use 'evolved-commit explain --help' for available rules.\n", ruleName)
			exit(1)
	}
}

func init() {
	rootCmd.AddCommand(explainCmd)
}
