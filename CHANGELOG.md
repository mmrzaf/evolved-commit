# Changelog

All notable changes to this project are documented here.

## Unreleased

- Introduce `run` subcommand with a basic stub functionality and test.
- Implement actual `pre-commit` Git hook installation, writing a script that calls `evolved-commit run`.
- Implement actual `pre-commit` Git hook uninstallation, safely removing only evolved-commit-managed hooks.
- Introduce `pkg/githooks` package with `Install` and `Uninstall` stubs and Git directory discovery logic.
- Integrate `pkg/githooks` into `install` and `uninstall` commands.
- Implement `findGitDir` utility function to locate `.git` directory.
- Add `uninstall` subcommand to remove Git hooks (basic structure).
- Add `install` subcommand to set up Git hooks (basic structure).
- Introduce `cobra` for CLI framework and basic root command structure.
- Initialize Go module and basic CLI structure (`main.go`).
- Repository bootstrapped.

- Introduce `run` subcommand with a basic stub functionality and test.

- Introduce `pkg/checks` package and initial `CheckCommitMessageSubjectNotEmpty` rule.
- Update `run` subcommand to execute commit message subject checks when a commit message file is provided.
