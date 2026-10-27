# Derive the run-binary source once

Blocked by: 01-measure-the-background-census.md
Writes: internal/gate/gate.go, internal/gate/subject.go, internal/gate/run_transaction.go, internal/gate/run_binary_source.go (new), internal/gate/run_binary_source_test.go (new), internal/gate/grade_identity_test.go, CHANGELOG.md
Covers: BG12, BG13, BG14, BG15

## What to build

Today the grade identity picks the kind of its run-binary source from the presence of `.bench/gate-prospective.sh` (`internal/gate/subject.go:93-100`).
The run picks the kind from the resolved kind of its plan (`internal/gate/run_transaction.go:161-162`).
Both then call `runBinarySource` (`internal/gate/gate.go:311-321`).
The two answers differ when `BENCH_KIT` names another checkout.

Make one derivation answer the source root of the run binary for a run.
The run selects from that answer, and the identity frames the build inputs of that same root.
A run that selects no run binary, because it has no owner or no phase-table gate, frames one stable value.
The author picks that value.

`subject.go` has 512 lines, which is over its budget. Keep its line count, or move code into `run_binary_source.go` in this ticket.

## Acceptance

- [ ] An ordinary kit run whose `BENCH_KIT` names another checkout frames the build inputs of that checkout.
- [ ] No grade identity code reads the presence of `.bench/gate-prospective.sh` to pick a kind. `rg -n "sourceKind" internal/gate` finds no hit.
- [ ] A run with no run-binary selection frames one stable identity value.
- [ ] One clean kit tree gives one run-binary member to its ordinary run and to its prospective run.
- [ ] `bench structure --growth` reports no growth of `internal/gate/subject.go`.
