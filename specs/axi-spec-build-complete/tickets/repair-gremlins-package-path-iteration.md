# Preserve changed package paths through the mutation referee

Blocked by: none
Ownership fence: `scripts/gremlins-diff.sh`, `internal/systemtest/gremlins_diff_test.go`
Integration surfaces: Git changed-Go-file names→stable lexical de-duplicated newline-delimited `dirs` in `scripts/gremlins-diff.sh`; `dirs` entries outside `testdata` whose directories still exist and contain live `*.go` files→eligible package entries in the same-shell Gremlins invocation loop; script exit/output→`internal/systemtest/gremlins_diff_test.go` subprocess oracle
Contracts: each eligible entry in `dirs` crosses `scripts/gremlins-diff.sh` into exactly one quoted filesystem check and Gremlins argv element without shell word splitting, asserted by GR1; no eligible entries after the `testdata`, existing-directory, and live-Go-file filters crosses to the explicit no-changed-package success, asserted by GR1; the real script's per-package and aggregate output crosses to the existing TestMain-owned repository oracle in `internal/systemtest/gremlins_diff_test.go`, asserted by GR1

## What to build

Close accepted Terra Standards finding S1-gremlins-package-path-splitting.
Iterate the selector's newline-delimited package paths with `IFS=` and `read -r`
in the current POSIX shell, fed from `dirs` by a here-document rather than a
pipeline, so spaces and glob characters remain literal and the aggregate
counters remain owned by the parent shell. Add a system test that reuses and
restores one of `internal/systemtest`'s three TestMain-owned repositories,
creates two changed Go package directories whose names contain spaces and
literal glob metacharacters, with at least two changed Go files in one directory,
injects a deterministic fake `gremlins`, and invokes the real script. The fake
records the complete operand sequence; the test requires each expected
`./<literal-dir>` exactly once in lexical order, exactly one per-package row for
each operand, and exact killed/lived aggregate totals covering both. Keep
changed-package granularity; do not restore the broken Gremlins `-D` filter or
allocate a fourth repository. Add a testdata-only changed-package control that
proves the preserved exclusion invokes no fake Gremlins process and takes the
explicit no-eligible-package success.
Add a deleted/no-live-Go changed-package control with the same no-invocation and
zero-eligible assertions.

## Acceptance

- [ ] [GR1] (covers local) (S1-gremlins-package-path-splitting) stable lexical de-duplicated `dirs` passes every eligible changed Go package path byte-for-byte as one Gremlins package operand despite spaces or glob characters, the per-package and exact aggregate verdicts include it exactly once, and no entry surviving the `testdata`, existing-directory, and live-Go-file filters is the only path to the explicit no-changed-package success.
