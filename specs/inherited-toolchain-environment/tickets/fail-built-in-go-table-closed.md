# Fail the built-in Go phase table closed

Blocked by: none
Writes: internal/gate/phases.go, internal/gate/manifest.go, internal/gate/branch_native_phases_test.go
Line: `gpt-5.6-terra` / medium — the gate seam is known, but a wrong fail posture weakens the oracle.

## What to build

The built-in phase-table path returns a diagnostic error when a repository has
a regular `go.mod` but PATH cannot resolve `go`; it never returns a schedule
with the Go phases omitted. The error names Go, PATH, and the graded repository.
A repository without `go.mod` keeps its non-Go built-in schedule, and a valid
project-owned phase manifest remains authoritative without a global Go
preflight.

## Acceptance

- [ ] A built-in table for `go.mod` plus Go-less PATH returns red before producing a schedule (covers TE1).
- [ ] The refusal names `go`, `PATH`, and the graded repository path (covers TE2).
- [ ] A repository without `go.mod` keeps its non-Go built-in schedule under the same PATH (covers TE3).
- [ ] A valid project-owned phase manifest returns exactly its declared phases under the same PATH (covers TE4).
