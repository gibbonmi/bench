# Repair control-bearing action values

Blocked by: repair-worktree-compatibility-closure.md
Writes: `internal/axi/`, `internal/worktree/list_actions_test.go`, `internal/coverage/coverage_test.go`

## What to build

Let known action arguments carrying TOON-supported tab, newline, or return bytes survive through the existing shell-quoting and TOON escaping pipeline while unsupported control bytes continue to refuse structurally.

## Acceptance

- [ ] [CV1] (covers QD1, QD6) `KnownArgument` renders tab, newline, and return values as executable shell arguments whose final help cell is TOON-safe and round-trippable; unsupported controls remain rejected.
- [ ] [CV2] (covers QD1, QD6) public coverage and orphaned-worktree queries with supported control-bearing paths preserve their primary result and exit while appending the exact escaped action instead of replacing the response with an error.

