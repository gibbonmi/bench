# 10. Record the commit and landing spans with their measures

Blocked by: keep-the-started-phase-line-after-a-kill.md, check-that-each-registered-seam-starts-a-span.md, record-the-lane-span-with-its-failing-check.md
Line: opus / medium
Rows: OT18, OT19
Writes: internal/commit/commit.go, internal/landing/landing.go, internal/landing/composition.go, internal/worktree/land.go, internal/otelrecord/registry.go, internal/systemtest/otel_verbs_test.go (new)

## What to build

`bench commit` writes a commit span at its verb boundary. The span carries the
subject digest, the outcome, and the composed path count. The span carries the
subject digest and never the subject text, because a commit subject holds
objective text by design. The record must not copy that text into a third
durable place, and `DATA_HANDLING.md` owns the reason.

`bench worktree land` writes a landing span that covers the composition and the
publication. The landing span carries the composed path count and the census
raw-call count, so the count survives the release.

Each seam derives its own count where it already holds the list. The commit
span counts the named attributed paths at the commit seam. The landing span
counts the reviewed name-only diff at the landing seam. No third counter enters
the tree.

The landing span reads the census raw-call count before the release deletes the
census record. This ticket adds the commit seam and the landing seam to the
registry from ticket 8, after ticket 9 writes it.

The system tests drive the built binary with a temporary `BENCH_HOME`. The gate
runs the `system` phase only when the graded root is the kit checkout. The
ticket-time observation is therefore the focused hand-run `go test -tags=system
./internal/systemtest`, with `BENCH_KIT` and `BENCH_RUN_BINARY` set.

## Acceptance

- [ ] OT18: `bench commit` writes a commit span that carries the subject digest, the outcome, and the composed path count.
- [ ] OT19: `bench worktree land` writes a landing span that carries the census raw-call count.
- [ ] the commit span carries no subject text.
- [ ] the landing span also carries the composed path count, and it counts the reviewed name-only diff.
- [ ] the registry names the commit seam and the landing seam, and the conformance check passes.
