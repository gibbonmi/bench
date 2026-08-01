# Derive the toolchain and contract input sets

Blocked by: Expose the build inputs and seal digests the gate reads
Ownership fence: new `internal/gate/component_inputs.go`, new
`internal/gate/component_inputs_test.go`
Assumptions: `ReducedScope()` in `scope.go` stays as the capture-surface
declaration and is not edited by this ticket; `canary.Phase*` constants name the
phases. Re-derive from the tree at pickup.

## What to build

The single source for what each scoped component reads. This ticket lands the
registry and the components whose input sets have a derivable source in Go:
`build`, `gofmt`, `vet`, `test`, `race`, and `contract`.

`build` takes its set from `freshness.BuildInputs` — the `go list -deps
./cmd/bench` closure. The toolchain and contract components take theirs from the
module-wide `go list -deps -test ./...` closure plus the `testdata/` directories
of listed packages, because the binary's closure excludes the conformance and
contract packages those components grade. `contract` additionally carries the
seal's source digest, since it execs the CLI and the binary is one of its inputs.

Every entry records how it was sourced — a named derivation or a hand
declaration — so the profile table and the derivation-source conformance check
have one thing to read. A component with no derivable source is not added here;
`shellcheck` and `canary` arrive in the next ticket against this same seam.

The registry is the enumerated family this build adds. Trace an existing sibling
before adding a member: `ReducedScope()`'s accessors, the profile's rendered
table, and the scope-binding conformance check are the registries a phase-name
family already appears in, and a component added here must appear in each of
them or the binding ticket reds.

## Acceptance

- [ ] PC8a — the toolchain and contract input sets contain every file of every package in the module-wide `-test` closure, including `_test.go` files of packages outside `./cmd/bench`'s closure.
- [ ] PC8b — the toolchain and contract input sets contain the `testdata/` contents of listed packages.
- [ ] PS9 — `build`'s input set is exactly `freshness.BuildInputs`' return, with no restatement of that path list in this package.
- [ ] PS10 — `contract`'s declaration carries the seal's source digest as an input alongside its file set.
- [ ] PS11 — every registry entry reports its source as a named derivation, and a derivation failure returns an error rather than a partial set.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC8a | change the toolchain derivation from `-deps -test ./...` to `-deps ./cmd/bench` | `TestToolchainInputsCoverTestFilesOutsideTheBinaryClosure` | build the kit-shaped fixture, resolve the toolchain input set, assert it contains the outside-closure package's `_test.go` file |
| PC8b | drop the `testdata/` walk from the derivation | `TestToolchainInputsCoverListedTestdata` | add a `testdata/` file to a fixture package, resolve the set, assert the file is present |
| PS9 | restate build's paths as a literal slice in this file | `TestBuildInputsAreNotRestated` | resolve `build`'s set and compare element-wise against `freshness.BuildInputs(root)` |
| PS10 | omit the seal digest from `contract`'s declaration | `TestContractInputsCarryTheSealSourceDigest` | publish a fixture binary, resolve `contract`'s declaration, republish a different binary, assert the declaration's digest input moved |
| PS11 | return the successfully-listed packages when `go list` exits nonzero | `TestComponentInputsErrorOnDerivationFailure` | corrupt the fixture's `go.mod`, resolve each derived component, expect an error and no paths |
