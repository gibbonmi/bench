# Rename the native-proof upload

Blocked by: none
Writes: .github/workflows/native-runtime.yml

## What to build

The native-proof job uploads its proof under the artifact name
`<prefix>-native-proof-<os>-<arch>`. The evidence job downloads that proof with
the pattern `<prefix>-native-proof-*`. The `first` segment is a vestige of the
retired two-generation shape, and no second upload exists. The producer name
and the download pattern move together. No other behavior changes.

## Acceptance

- [ ] The `native-proof` job's upload step names the artifact `${{ inputs.artifact-prefix || 'bench-runtime' }}-native-proof-${{ matrix.os }}-${{ matrix.arch }}`.
- [ ] The `evidence` job's download step uses the pattern `${{ inputs.artifact-prefix || 'bench-runtime' }}-native-proof-*`.
- [ ] No file in the tree outside this ticket folder, dot-directories included, still contains the bytes `native-proof-first`.
- [ ] The conformance workflow checks in `internal/conformance/native_workflow_test.go` stay green.
