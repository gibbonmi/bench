# Skip the build phase on its attested seal

Blocked by: Attest the build the gate itself ran; Skip evidence-covered
components on their own slots
Ownership fence: the build branch of `internal/gate/component_decision.go`, the
build-skip wiring in `internal/gate/gate.go`,
`internal/gate/build_attestation_test.go`
Assumptions: the `Needs` edges in `BenchkitPhases` make every reader of
`dist/bench` depend on the build phase, and a need that ends skipped is
satisfied trivially with no writer and no race; the gate-entry untrusted-runner
check in `gate.sh` is a different surface and does not change. Re-derive from
the tree at pickup.

## What to build

The last component to scope, and the only one that skips through artifact reuse
rather than an ancestor slot. The build phase executes when `dist/bench` is
absent, when `freshness.Check` refuses for any reason, or when the seal's
executable digest is not attested by a prior gate-run build. A valid attested
seal skips the phase, and the binary readers exec the sealed binary unchanged.
mtime plays no part in any of it.

A green build republishes the seal and re-authors the attestation. `--fresh`
executes the build like everything else.

**Evidence authorship.** `bench gate` authors the attestation when it runs the
build green; a run that skips the build leaves the attestation and the seal
untouched.

## Acceptance

- [ ] PC3 — a valid attested seal skips the build and leaves `dist/bench` byte-identical, and the phases that exec it exec the sealed binary.
- [ ] PC4 — each of absent binary, absent seal, source-digest mismatch, executable-digest mismatch, symlinked or irregular sidecar, missing attestation, and attestation/seal mismatch runs the build, one subtest per case.
- [ ] PC5b — a planted binary with a recomputed self-consistent seal runs the build, which overwrites both and re-authors the attestation.
- [ ] PC20 — a first run, a pruned evidence store, and `--fresh` each execute every component including the build, and author every slot and the attestation.
- [ ] PS30 — the verdict's evidence entry for a skipped build names the seal's source digest.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PC3 | rebuild whenever the binary's mtime predates any source | `TestAttestedSealSkipsTheBuild` | seed a green kit-shaped fixture, record `dist/bench` bytes, touch a source's mtime only, execute, byte-compare |
| PC4 | drop the attestation conjunct, keeping only `freshness.Check` | `TestBuildRunsOnEveryUnsoundArtifact` | apply each of the seven faults in its own subtest, execute, assert the build phase ran |
| PC5b | compare the attestation against the seal instead of the binary | `TestPlantedArtifactRunsTheBuild` | plant a foreign binary and recompute its seal, execute, assert the build ran and the attestation moved |
| PC20 | treat an absent attestation as attested | `TestFirstRunAndFreshBuildEverything` | run against a fresh clone, a pruned store, and `--fresh` in three subtests, assert the build executed and the attestation exists |
| PS30 | record the executable digest as the build's evidence | `TestSkippedBuildEvidenceNamesTheSourceDigest` | execute a run that skips the build, load the verdict, compare the entry against `freshness.SealDigests`' source digest |
