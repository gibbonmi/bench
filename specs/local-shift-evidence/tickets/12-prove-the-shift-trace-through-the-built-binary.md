# 12. Prove the shift trace through the built binary

Blocked by: 2-redact-and-version-every-record-line.md, 5-record-each-pass-under-the-shift.md, 6-record-the-resolved-line.md
Writes: internal/systemtest/otel_verbs_test.go
Covers: LE84, LE85

## What to build

Chunk: LE-D.

Add one system journey to the verb span tests. It runs `bench shift` through the built binary in a disposable repository with a private Bench home. The adapter script calls `bench resolve-model` through the wrapper, appends a line to a tracked file, and exits 0. The gate fixture is green.

The objective holds `OBJMARK`, the adapter directory name holds `DIRMARK`, and `OTEL_RESOURCE_ATTRIBUTES` holds `ENVMARK`. The journey reads every record segment below the private home. `internal/systemtest/` has no directory headroom, so the journey joins the existing verb span file. The system suite files are system-tagged, so run the journey with `BENCH_KIT` set, through `bench test --check system`.

## Acceptance

- [ ] The run writes one trace in which the `shift` span is an ancestor of the pass span, the `line.resolve` span, and the `gate` span.
- [ ] No record segment holds `OBJMARK`, `DIRMARK`, or `ENVMARK`.
