# ROADMAP

## North Star

Every commit is review-ready by default.

## Core Capabilities

- CLI: install / uninstall / run / explain
- Git hook management (safe, idempotent)
- CI-compatible execution
- Clear, actionable output

## Baseline Rules

- Commit subject format and length
- Branch naming conventions
- Forbidden diff patterns (debug prints, TODOs without context)
- Commented-out code detection

## Rule System

- Rules are isolated and order-independent
- Each rule has:
  - unique id
  - explanation
  - fix guidance
  - test cases

## UX Backlog (never empty)

- Warn vs enforce modes
- Safe autofix for mechanical rules
- Language-aware rule packs
- Ignore and allowlist support
- Improved error messaging

## Quality Backlog

- Reduce false positives
- Improve performance on large diffs
- Sharpen rule explanations
