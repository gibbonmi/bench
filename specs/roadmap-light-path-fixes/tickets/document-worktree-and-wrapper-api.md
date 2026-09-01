# Document worktree and wrapper contracts

Blocked by: add-retrospective-writer.md, retain-explicit-signal-safe-worktree-shell.md
Writes: internal/worktree/lifecycle.go, internal/worktree/ownership.go, internal/worktree/worktree.go, internal/worktree/subshell.go, cmd/bench/main.go, tests/canary/package-core-guard/unrouted-subcommand, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LF19

## What to build

Document ProbeLease, Acquire, Release, Create, Pool, LeaseFile,
ClassifyRegisteredWorktrees, CleanCommand, ReleaseCommand, CreateCommand,
Subshell, and top-level worktree dispatch. State lease, ownership, cleanup,
and parser boundaries without port narration.

## Acceptance

- [ ] All eleven named functions state their caller-facing contracts.
- [ ] Cleanup comments distinguish planning from mutation.
- [ ] Wrapper dispatch identifies which leaf owns each grammar.
- [ ] Edited comments contain no provenance or reviewer argument.
