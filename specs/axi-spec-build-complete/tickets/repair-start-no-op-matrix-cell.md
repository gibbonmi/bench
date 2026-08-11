# Classify the real start no-op in the AXI matrix

Blocked by: none
Ownership fence: `internal/specbuild/assign.go`, `internal/specbuild/render.go`, `internal/specbuild/disclosure.go`, `internal/specbuild/disclosure_observation.go`, `internal/specbuild/disclosure_test.go`, `internal/specbuild/start_test.go`, `internal/specbuild/testdata/axi-ledger.md`, `internal/specbuild/testdata/axi-cases.jsonl`
Integration surfaces: existing-run Start branch→one internal typed start outcome owned beside `Service.Start`; typed start outcome→`ExecuteObserved` disclosure class; `start/no-op` applicability→`DisclosureCells`; real repeated-start observation→checked ledger and `matrix/start/no-op` fixture
Contracts: the Start owner supplies one status plus its typed success/no-op outcome to the execution boundary without changing the public `Service.Start` signature, asserted by NS1; the newly applicable cell crosses `DisclosureCells`→`ObserveDisclosureCell`→the checked ledger and historical fixture, asserted by NS1

## What to build

Close accepted Terra Spec finding P1-start-no-op-matrix-exclusion. Preserve the
public `Service.Start` API while making its internal owner report whether it
created/resumed work or returned an already-active run unchanged. Consume that
single outcome in `ExecuteObserved`; do not infer no-op again from rendered bytes,
state strings, or test-only flags. Mark only the service's real `start/no-op` cell
applicable, prepare it by performing a successful start before the one public
observation, and check in its exact old/new payload and corrected ledger counts.
All other no-op dispositions remain explicit not-applicable entries.

## Acceptance

- [ ] [NS1] (covers SB2) (P1-start-no-op-matrix-exclusion) a repeated public start that returns the existing active run unchanged is typed `no-op`, exits 0, appears as applicable `start/no-op`, and is exercised through the real observation service against `matrix/start/no-op`; fresh start remains typed `success` and every other no-op cell retains its exact not-applicable disposition.
