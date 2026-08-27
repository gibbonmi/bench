# Resolve the test roots without runtime.Caller

Blocked by: none

Writes: internal/runbinary/bootstrap_test.go, internal/conformance/harness_test.go, internal/preprelease/preprelease_test.go, internal/bounds/bounds_test.go

## What to build

Three test helpers derive the repository root from `runtime.Caller`. Go strips
the absolute file name under `-trimpath`, so each helper builds a path from the
module path and its package reds.
On 2026-08-27 the command `go test -trimpath ./internal/runbinary` failed seven
tests with `lstat github.com: no such file or directory`.

Move each helper to a working-directory resolution, or to the git top level.
The conformance harness already holds a git-based resolver beside the broken
one; mirror that resolver. Keep every existing assertion unchanged, because
this ticket repairs the root lookup only.

This ticket is the prefactor for `03-add-trimpath-to-every-bench-go-argv.md`.
That ticket puts `-trimpath` on the argvs, and these three packages must
already survive the flag.

## Acceptance

- [ ] T06 — `go test -trimpath` on `internal/runbinary`, `internal/conformance`, `internal/preprelease`, and `internal/bounds` is green.

Delivered outcome: the three packages that read their own source tree stay
green when Go strips the file paths.
