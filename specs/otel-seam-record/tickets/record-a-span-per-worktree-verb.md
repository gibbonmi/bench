# 11. Record a span per worktree verb

Blocked by: record-the-commit-and-landing-spans.md
Line: opus / medium
Rows: OT20
Writes: internal/worktree/worktree.go, internal/worktree/exec.go, internal/worktree/merge.go, internal/worktree/land.go, internal/otelrecord/registry.go, internal/systemtest/otel_worktree_test.go (new)

## What to build

Each worktree verb writes a span that names the verb and the assignment, so the
worktree traffic lands in the record. An uninstrumented verb leaves no line for
its own invocation.

The spans carry only the attributes that ticket 3 declares. The assignment id
passes through `encoding/json`, which escapes every control rune, so a subject
id cannot forge a second record line.

Ticket 10 edits `internal/worktree/land.go` and the registry, and this ticket
edits both files. The coordinator therefore runs this ticket after ticket 10.

The gate runs the `system` phase only when the graded root is the kit checkout.
The ticket-time observation is therefore the focused hand-run `go test
-tags=system ./internal/systemtest`, with `BENCH_KIT` and `BENCH_RUN_BINARY`
set.

## Acceptance

- [ ] OT20: each worktree verb writes a span that names the verb and the assignment.
- [ ] a span for one verb carries no attribute outside the declared set.
- [ ] the registry names each worktree seam, and the conformance check passes.
