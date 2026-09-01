# Single-source the resume test golden

Blocked by: none
Writes: internal/worktree/resume_test.go, internal/worktree/lifecycle_policy_test.go, internal/worktree/resume_summary_test.go (new), cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/main_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_test.go
Covers: LF3

## What to build

Create one test-only expected-format helper for the resume summary. Both unit
and runtime-binary seams consume it. Production rendering remains independent.

## Acceptance

- [ ] Both resume test seams use one expected-format helper.
- [ ] The production summary remains byte-for-byte compatible.
- [ ] No production formatter supplies its own oracle.

