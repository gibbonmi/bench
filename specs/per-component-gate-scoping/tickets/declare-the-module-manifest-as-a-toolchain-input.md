# Declare the module manifest as a toolchain input

Blocked by: Declare shellcheck's and canary's input sets
Ownership fence: `internal/gate/component_inputs.go`, `internal/gate/component_inputs_test.go`
Assumptions: `moduleTestClosure` derives from `go list -deps -test ./...` plus
the `testdata/` directories of listed packages, and `go list` reports neither
`go.mod` nor `go.sum`; `freshness.BuildInputs` already adds both explicitly, so
`build` is unaffected. Re-derive all three from the tree at pickup.

## What to build

A dependency version bump edits `go.mod` and `go.sum` while leaving the module
closure's file set byte-identical. Under the derivation as it stands, `gofmt`,
`vet`, `test`, `race`, and `contract` would each declare themselves unmoved by a
change that can red all five — a component skipping on evidence graded against a
different dependency set, which is the exact failure this feature exists to
prevent.

The toolchain and contract components declare the module manifest alongside
their derived closure. `build` already sets the precedent inside this same
feature: `freshness.BuildInputs` derives its closure and then adds `go.mod`,
`go.sum`, and the auxiliary manifest explicitly. This entry follows that shape,
and its recorded source says so — derivation plus a named, bounded addition,
rather than a pure derivation. The reviewer approved widening it on 2026-08-01;
the alternative was recording the gap as a profile residual, which was declined.

The addition is bounded and named, not a licence to hand-extend the set: only
the module manifest files join, `go.sum` only when it exists, and the derivation
still supplies every source path.

## Acceptance

- [ ] [PC8c] the toolchain and contract input sets contain `go.mod`, and contain `go.sum` when the module has one.
- [ ] [PC8d] a `go.mod` edit that leaves the closure's file set unchanged moves those components' declared inputs.
- [ ] [PS37] the toolchain entries report a source naming both halves — the module-test closure and the manifest addition — so the derivation-source check can tell this shape from a hand-copied list.
- [ ] [PS38] a module with no `go.sum` resolves without error and without a phantom entry.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC8c | drop the manifest addition from `moduleTestClosure` | `TestToolchainInputsCoverTheModuleManifest` | build the kit-shaped fixture, resolve `gofmt`'s set, assert `go.mod` and `go.sum` present |
| PC8d | add the manifest paths to the recorded source label only, not to the path set | `TestModuleManifestEditMovesToolchainInputs` | resolve, bump a require line in the fixture's `go.mod`, resolve again, assert the sets differ |
| PS37 | keep the bare module-test-closure source label | `TestToolchainSourceNamesTheManifestAddition` | enumerate the registry, assert the toolchain entries' source names the addition |
| PS38 | add `go.sum` unconditionally | `TestModuleWithoutGoSumResolves` | remove `go.sum` from the fixture, resolve, assert no error and no `go.sum` entry |

## Note for the derivation-source check

The conformance check in "Red the gate on a copied list instead of a derivation"
must accept this two-part shape for the toolchain and contract entries — a named
derivation plus a bounded manifest addition — while still redding a literal path
list. `canary` remains the only fully hand-declared entry.
