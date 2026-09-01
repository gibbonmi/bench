# Add an elapsed_ms column to the bench test packages table

Blocked by: none
Writes: internal/testreport/testreport.go, internal/testreport/testreport_test.go, internal/testreport/check_test.go, CHANGELOG.md
Covers: none

## What to build

`bench test` prints no wall time, so a caller cannot tell a real run from a
skipped one. The caller must put a `date` call on each side of the run to get
that time.

The `packages` table gets one more column, `elapsed_ms`. The table then reads
`packages[N]{package,status,elapsed_ms}`. Each row gives the milliseconds that
`go test` reports for that one package. The decoder reads the per-package
`Elapsed` field of the `go test -json` stream. That field gives seconds as a
float, and the report keeps it as an integer count of milliseconds.

The name `elapsed_ms` is the same name the gate's `phases` table uses, so the
two tables read alike. A stream that gives no `Elapsed` for a package makes a
`0` cell. The `failures`, `skips`, and `root_conformance` tables do not change.

## Acceptance

- [ ] `bench test --package ./internal/sanitize` prints
      `packages[1]{package,status,elapsed_ms}` with a positive integer cell.
- [ ] A package whose `go test -json` stream carries no `Elapsed` field prints a
      `0` cell.
- [ ] Every test in `internal/testreport` stays green.
