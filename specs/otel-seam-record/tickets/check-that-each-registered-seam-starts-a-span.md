# 8. Check that each registered seam starts a span

Blocked by: record-the-gate-run-and-its-phases.md
Line: opus / medium
Rows: OT26
Writes: internal/otelrecord/registry.go (new), internal/conformance/otel_seam_test.go (new)

## What to build

The package `internal/otelrecord` exports the seam registry. Each entry names
one seam and the Go symbol that starts that seam's span. The registry is the
one source for the instrumented set, and no second list exists.

A new `go/ast` check in `internal/conformance` enumerates the registry and
reads each named symbol. The check reds when a registered symbol starts no
span, so no seam stays silently uninstrumented. The loader in
`internal/consumers` is the prior art.

The gate seam from ticket 6 is the registry's first entry. Each later
instrumentation ticket adds its own entry, and the check reds until that seam
starts its span.

## Acceptance

- [ ] OT26: the check returns one diagnostic that names a registered seam whose symbol starts no span.
- [ ] the check passes over the live tree with the gate seam registered.
- [ ] the registry names the gate seam and the symbol that starts its span.
