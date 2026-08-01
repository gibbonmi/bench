# Declare shellcheck's and canary's input sets

Blocked by: Derive the toolchain and contract input sets
Ownership fence: the `shellcheck` and `canary` entries of
`internal/gate/component_inputs.go`, their cases in
`internal/gate/component_inputs_test.go`
Assumptions: `shellcheckArgv(root)` in `phases.go` already enumerates the exact
file list shellcheck lints, and is read — not edited — by this ticket. Re-derive
from the tree at pickup.

## What to build

The two remaining scoped components' declarations, against the registry seam the
previous ticket landed.

`shellcheck` derives from `shellcheckArgv`'s own enumeration, so the component
that lints a file and the declaration that says it reads that file cannot
disagree — adding a script to the argv adds it to the declaration with no second
edit.

`canary` is the one hand-declared entry: `internal/canary/`, `tests/canary/`,
and the wrapper scripts its phase execs (`bin/bench.sh` and the canary dispatch
behind it). The binary digest is deliberately excluded, by the reviewer ruling
of 2026-08-01 — a recorded tripwire narrowing, not an oversight. State that
exclusion where the declaration is written, so a later reader cannot mistake it
for a missing input.

## Acceptance

- [ ] PS12 — `shellcheck`'s input set is derived from `shellcheckArgv`'s file list, and a script added to that argv appears in the declaration with no edit here.
- [ ] PS13 — `canary`'s input set contains `internal/canary/`, `tests/canary/`, and the wrapper scripts its phase execs, and contains no binary digest.
- [ ] PS14 — `canary`'s registry entry reports itself hand-declared while every other entry reports a named derivation.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PS12 | replace the argv-sourced set with a literal list of the same paths | `TestShellcheckInputsFollowItsArgv` | add a `.bench/hooks/new.sh` to the kit-shaped fixture, resolve `shellcheck`'s set, assert the new path is present |
| PS13 | add the seal's executable digest to `canary`'s declaration | `TestCanaryInputsExcludeTheBinary` | publish a different binary into the fixture with sources unchanged, resolve `canary`'s declaration twice, assert it did not move |
| PS14 | mark `canary` as derived | `TestOnlyCanaryIsHandDeclared` | enumerate the registry, assert exactly one hand-declared entry and that it is `canary` |
