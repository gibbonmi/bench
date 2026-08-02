# Identify a component by its inputs and execution closure

Blocked by: Derive the toolchain and contract input sets
Ownership fence: `internal/gate/component_identity.go`, `internal/gate/component_identity_test.go`, `internal/gate/subject.go`, `internal/gate/tree_snapshot.go`, `internal/gate/tree_snapshot_test.go`
Assumptions: `strippedTreeHash` in `subject.go` reads a `git add -A` snapshot
through `git ls-tree` and drops declared entries; `strippedPolicyVersion`
extends `policyVersion`. This ticket reads both and edits neither. Re-derive from
the tree at pickup.

## What to build

The content address a component's ancestor slot is stored under. Each component
hashes its declared input contents by positive selection over the same
`git add -A` snapshot the stripped identity reads today — selecting what the
declaration names rather than dropping what the capture allowlist names — plus
its execution closure: argv shape, env contract, and, where the binary is one of
its inputs, the seal's source digest.

Each component hashes under its own policy domain, derived from the component
name, so no slot can answer for another component even if two declarations
happen to cover identical files. A component whose identity cannot be computed —
snapshot unreadable, derivation failed, a declared path outside the repository —
returns an error; the decision function turns that into running the component,
and this ticket never invents a fallback identity that would let it skip.

## Acceptance

- [ ] [PS15] a component's identity moves when any file in its declared input set changes and stays fixed when a file outside it changes.
- [ ] [PS16] two components with identical declared input sets compute different identities, because the policy domain carries the component name.
- [ ] [PS17] a component's identity moves when its argv or declared env changes with every input file unchanged.
- [ ] [PS18] `contract`'s identity moves when the seal's source digest moves with its file set unchanged.
- [ ] [PS19] an unreadable snapshot, a failed derivation, and a declared path outside the repository each return an error and no identity.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PS15 | hash the declared path names without their contents | `TestComponentIdentityTracksDeclaredContent` | edit a declared file in the kit-shaped fixture, recompute, assert the identity moved; edit an undeclared file, assert it did not |
| PS16 | drop the component name from the policy domain | `TestComponentIdentitiesAreDomainSeparated` | declare two components over one file set, compute both, assert inequality |
| PS17 | hash only the input contents, omitting the execution closure | `TestComponentIdentityTracksItsExecutionClosure` | change a phase's argv in the fixture manifest, recompute, assert the identity moved |
| PS18 | omit the seal digest from `contract`'s closure frame | `TestContractIdentityTracksTheSeal` | republish a different binary with sources unchanged, recompute `contract`'s identity, assert it moved |
| PS19 | return a zero identity on a derivation error | `TestComponentIdentityFailsClosed` | corrupt `go.mod`, remove the snapshot, and declare an escaping path in three subtests, assert an error and an empty identity each time |
