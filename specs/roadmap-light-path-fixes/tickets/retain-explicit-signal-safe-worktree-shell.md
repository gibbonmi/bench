# Retain an explicit signal-safe worktree shell

Blocked by: make-bare-worktree-usage-safe.md
Writes: cmd/bench/main.go, cmd/bench/main_test.go, cmd/bench/command_registry_test.go, internal/worktree/subshell.go, internal/worktree/subshell_test.go, internal/worktree/worktree.go
Covers: LF14

## What to build

Move the human shell journey to bench worktree shell. Release normally and
preserve enough lease state for reclamation after interrupt or termination.

## Acceptance

- [ ] bench worktree shell opens the existing human shell journey.
- [ ] Normal exit releases its worktree.
- [ ] Interrupt and termination leave no unreclaimable worktree.
- [ ] Machine-facing path output remains literal and absolute.
