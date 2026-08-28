# Switch the native-proof matrix to the proven targets

Blocked by: add-proven-target-field.md
Writes: .github/workflows/native-runtime.yml, internal/conformance

## What to build

The `native-proof` job stops starting a macOS runner.

The `preflight` job publishes a second matrix output built from `proof-matrix-json`. The `native-proof` job reads that output, and the `smoke` job keeps reading the shipped matrix, so the shipped macOS binary still runs on macOS.

`checkNativeRuntimeWorkflow` grades both bindings. A `native-proof` job reading the shipped matrix is a diagnostic, and a `smoke` job reading the proven matrix is a diagnostic.

## Acceptance

- [ ] The `native-proof` job reads the proven matrix output (row B4).
- [ ] The conformance check names a diagnostic when `native-proof` reads the shipped matrix (row B4).
- [ ] The `smoke` job reads the shipped matrix output (row B5).
- [ ] The conformance check names a diagnostic when `smoke` reads the proven matrix (row B5).
