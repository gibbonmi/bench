# 12. Record a span per hook plumbing verb

Blocked by: keep-the-started-phase-line-after-a-kill.md, check-that-each-registered-seam-starts-a-span.md, record-the-commit-and-landing-spans.md
Line: opus / medium
Rows: OT21
Writes: cmd/bench/main.go, internal/otelrecord/registry.go, internal/systemtest/otel_hook_test.go (new)

## What to build

The hook plumbing verbs share one dispatch seam in `cmd/bench`. This ticket
instruments that one seam, and it instruments no adapter package. Each hook
verb then writes a span that names the verb and its exit, so the harness events
join the record.

This ticket adds the hook dispatch seam to the registry from ticket 8. Ticket
10 writes the registry before it, so the coordinator runs this ticket after
ticket 10.

The system test pipes a hook envelope to the built binary with a temporary
`BENCH_HOME`, then reads the record lines back. The gate runs the `system`
phase only when the graded root is the kit checkout. The ticket-time
observation is therefore the focused hand-run `go test -tags=system
./internal/systemtest`, with `BENCH_KIT` and `BENCH_RUN_BINARY` set.

## Acceptance

- [ ] OT21: a hook plumbing verb writes a span that names the verb and its exit.
- [ ] the dispatch seam holds the one instrumented call, and no adapter package starts a span.
- [ ] the registry names the hook dispatch seam, and the conformance check passes.
