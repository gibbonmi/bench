# Repair registry and deep-CWD proof

Blocked by: repair-coverage-disclosure-compatibility.md
Writes: `cmd/bench/command_registry_test.go`, `internal/conformance/axi_query_registry_test.go`

## What to build

Make anchors' deep-CWD case enter the repository-relative path branch, and keep the one sanctioned independently authored AXI-set expectation in conformance rather than duplicating it in package-main tests.

## Acceptance

- [ ] [RR1] (covers QD3) the hermetic anchors fixture materializes `.bench/BENCH.md`; a mutation that bypasses repository-relative normalization makes root and deep-CWD results differ.
- [ ] [RR2] (covers QD3) package-main envelope cases are derived from the production registry declarations and bind every declared member, while `internal/conformance` remains the sole independent approved-set expectation.
- [ ] [RR3] (covers QD3) no fixed registry-entry count is repeated outside the exhaustive disposition derivation; adding a properly exempted production command does not require editing an unrelated count literal.
