# 10. Record the worktree shell session

Blocked by: 4-record-the-shift-boundaries.md
Writes: internal/worktree/subshell.go, internal/worktree/clean.go, internal/worktree/verb_span.go, internal/worktree/subshell_test.go, internal/worktree/otel_seams_test.go, cmd/bench/command_registry.go, cmd/bench/command_registry_test.go, cmd/bench/help_inventory_test.go, internal/conformance/axi_query_registry_test.go, internal/conformance/subcommand_routing_table_test.go, internal/otelrecord/registry.go
Covers: LE80, LE81, LE82, LE83

## What to build

Chunk: LE-D.

Open one `worktree.shell` span before `bench worktree shell` creates its assignment, and end it at every return. The end line carries the assignment id, the work state, the cleanup decision, and the outcome. Use the work-state and cleanup vocabularies from ticket 4.

A normal shell exit is `completed`, with cleanup `released` when the release exits 0 and `retained` otherwise. A signal is `interrupted` with `retained`, because the lease stays for a reclaim. A shell that cannot start is `failed`.

Add the seam to the registry and to the package's seam constants, so the seam equality test holds both. `internal/worktree/subshell.go` is over its line budget, so move the ignored-inventory function into `internal/worktree/clean.go` and leave `subshell.go` shorter. The command registry and its conformance tests join the Writes line through the binding closure only, and the build expects no edit there.

## Acceptance

- [ ] A normal shell exit writes a `worktree.shell` span whose end line carries the assignment id, work state `completed`, and cleanup `released`.
- [ ] A signalled session writes an end line with work state `interrupted` and cleanup `retained`.
- [ ] A session whose shell path does not exist writes an end line with work state `failed`.
- [ ] `TestWorktreeSeamsMatchTheRegistry` holds `worktree.shell` on both sides.
