# Make bare worktree usage safe

Blocked by: none
Writes: cmd/bench/main.go, cmd/bench/command_registry_test.go, internal/worktree/worktree.go, internal/worktree/worktree_test.go, tests/canary/package-core-guard/unrouted-subcommand, cmd/bench/command_registry.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LF13

## What to build

Make bare bench worktree print usage and acquire nothing. Reject unknown flags
and subcommands before choosing any human-facing fallback.

## Acceptance

- [ ] Bare worktree exits with usage and creates no assignment.
- [ ] Unknown flags and subcommands exit two before acquisition.
- [ ] Supported leaf commands keep parser-first dispatch.

