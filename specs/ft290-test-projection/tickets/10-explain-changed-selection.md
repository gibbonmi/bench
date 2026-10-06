# 10. Explain each changed selection

Blocked by: 2-count-tests-run.md
Writes: internal/testreport/
Covers: TP41, TP42, TP43, TP44, TP45, TP46, TP48

## What to build

Chunk: TP-C4.

Make the changed-package selector return one cause for each selected package.
The precedence is `go-metadata`, then `changed`, then `embed`, then `imports <package>`.

For `imports`, name the first selected direct dependency in sorted import-path order, across the imports, the test imports, and the external test imports.
A `--changed` result prints `packages[N]{package,status,elapsed_ms,tests_run,selected_by}`, and an empty selection prints that header with zero rows.

Each other form keeps the header of ticket 2.

The selector returns a bare package list today, and its callers change with its return type. The empty-selection path builds a new report, so it also needs the cause input.
`TestChangedNonGoSubjectRendersExplicitEmpty` checks only the `packages[0]` prefix today. Add an exact header assertion with the `selected_by` cell.
The empty header names the empty selection, so print no separate marker.

## Acceptance

- [ ] One changed Go file prints the cause `changed`.
- [ ] A package that changed and that imports a changed package prints `changed`.
- [ ] A package with only a changed embed file prints `embed`.
- [ ] A package that imports two selected packages names the first one in sorted order.
- [ ] A package two import steps from the change names its direct dependency.
- [ ] A changed `go.mod` prints `go-metadata` on each row.
- [ ] An empty selection prints the `selected_by` header with zero rows at exit 0.
