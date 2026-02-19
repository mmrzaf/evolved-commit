# evolved-commit

Make every Git commit review-ready by default.

`evolved-commit` installs safe Git hooks and runs fast, opinionated checks on
commit messages, branch names, and staged diffs. When something fails, it tells
you exactly why and how to fix it.

## Usage

- `evolved-commit install` – install Git hooks
- `evolved-commit uninstall` – remove hooks
- `evolved-commit run` – run checks manually or in CI
- `evolved-commit explain <rule>` – explain a rule and how to satisfy it

## Principles

- Zero configuration by default
- Deterministic, offline, fast
- Clear failures with exact fix instructions

See `ROADMAP.md` for planned rules and improvements.
