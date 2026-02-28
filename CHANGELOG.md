# Changelog

All notable changes to this project are documented here.

## Unreleased

- Introduce `pkg/githooks` package with `Install` and `Uninstall` stubs and Git directory discovery logic.
- Integrate `pkg/githooks` into `install` and `uninstall` commands.
- Implement `findGitDir` utility function to locate `.git` directory.
- Add `uninstall` subcommand to remove Git hooks (basic structure).
- Add `install` subcommand to set up Git hooks (basic structure).
- Introduce `cobra` for CLI framework and basic root command structure.
- Initialize Go module and basic CLI structure (`main.go`).
- Repository bootstrapped.

- Introduce `pkg/githooks` package with `Install` and `Uninstall` stubs and Git directory discovery logic.
- Integrate `pkg/githooks` into `install` and `uninstall` commands.
- Implement `findGitDir` utility function to locate `.git` directory.
