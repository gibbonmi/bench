# Share the kit-shaped staged binary

Follow-up: `reduce-gate-fixture-materialization.md` retains this ticket's process template and real-build boundary but supersedes its ordinary per-root copy with measured hardlink-and-detach materialization.

Blocked by: none
Ownership fence: `internal/gate/kitshaped_fixture_test.go`, `internal/gate/build_attestation_test.go`
Integration surfaces: `internal/gate/kitshaped_fixture_test.go` constructor-only initial seal→process-scoped staged template + SFB1; explicit re-seal helper→real build + SFB3; template→per-root copy→existing freshness publication + SFB5-SFB7; package lifetime→sticky error and cleanup + SFB8-SFB9; `internal/gate/build_attestation_test.go` error-returning build primitive→unchanged testing wrapper + SFB2; changed-source re-seals in `internal/gate/component_inputs_test.go`, `internal/gate/component_identity_test.go`, and `internal/gate/check_slots_test.go`→unchanged real-build helper + SFB3; plant and alternate/authorship builds in `internal/gate/build_skip_test.go` and `internal/gate/build_attestation_test.go`→unchanged real-build helper + SFB4; measured package-wide materialization residual→`reduce-gate-fixture-materialization.md` + GFM1-GFM3
Contracts: one immutable executable path and its construction error cross the process-scoped once in `internal/gate/kitshaped_fixture_test.go`→each mutable fixture root, ordered build once then ordinary copy then per-root publish, with absence represented by the retained error and asserted by SFB1/SFB5/SFB8; the error-returning primitive in `internal/gate/build_attestation_test.go` stays beneath the existing testing wrapper asserted by SFB2; constructor-only reuse excludes every unchanged real-build family asserted by SFB3-SFB4
Closure: SFB1/lazy-template-build, SFB2/error-returning-primitive, SFB2/testing-wrapper, SFB3/component-input-reseal, SFB3/component-identity-reseal, SFB3/check-slot-reseal, SFB4/build-skip-plant, SFB4/attestation-plant, SFB4/alternate-artifact, SFB4/authorship, SFB5/per-root-copy, SFB6/publication-consumption, SFB7/root-isolation, SFB8/sticky-error, SFB9/process-cleanup, SFB13/focused-output, SFB14/race, SFB15/two-core, SFB16/hostile-kit

## What to build

Build the deterministic kit-shaped fixture binary once per `go test` process
behind `sync.Once`. Keep the immutable template in a package-scoped temporary
directory, copy its bytes into each fixture root's `dist/bench.staged`, and pass
only that root-owned copy to the existing freshness publisher. Package cleanup
removes only the template directory after all tests finish. The once retains
both path and construction error so every later fixture sees the same failed
construction instead of observing an empty or consumed template.

Reuse applies only to the byte-identical initial build in
`newKitShapedFixture`, through an initial-seal helper that no changed-source
caller can enter accidentally. The existing re-seal helper always performs a
real build. Add an error-returning build primitive beneath the
existing `testing.T` build wrapper; the template builder consumes the primitive
and every existing direct caller keeps the wrapper interface. Re-seals after
source edits in the component-input, component-identity, and check-slot tests,
plus planted, alternate-artifact, and authorship builds in build-skip and
build-attestation tests, continue through that real build wrapper. Use an
ordinary portable copy; hardlinks would let one mutable fixture corrupt another,
and direct publication would consume the shared staging path.

## Acceptance

- [ ] [SFB1] the first constructor-only kit-shaped fixture lazily authors one process-scoped staged binary and every later constructor reuses that immutable template.
- [ ] [SFB2] the real-build implementation returns errors beneath the existing `testing.T` wrapper and every unchanged direct caller retains that wrapper interface.
- [ ] [SFB3] changed-source re-seals in component-input, component-identity, and check-slot tests continue to perform real builds.
- [ ] [SFB4] planted binaries, alternate artifacts, and authorship assertions in build-skip and build-attestation tests continue to perform real builds.
- [ ] [SFB5] each constructor copies the immutable template into its own `dist/bench.staged` path before publication.
- [ ] [SFB6] each fixture passes only its root-owned staged copy through the existing freshness publisher, so publication never consumes the template.
- [ ] [SFB7] mutating one published fixture binary cannot change any other fixture or the process template.
- [ ] [SFB8] one template-construction failure is retained and returned to every later dependent caller.
- [ ] [SFB9] package teardown removes only the process-scoped template directory and leaves fixture-owned artifacts to their existing cleanup owner.
- [ ] [SFB13] `/usr/bin/time -v env GOMAXPROCS=2 go test -p 1 -parallel 1 -count=1 -run '^(TestKitShapedFixtureCarriesBuildAndCanary|TestKitShapedFixtureBinaryIsSealed|TestKitShapedFixturesPublishIndependentTemplateCopies)$' ./internal/gate` reports at most 60 percent of the same command after the per-root-link mutation, proving reuse across four fixture constructions in one process.
- [ ] [SFB14] `GOMAXPROCS=2 go test -race -p 1 -parallel 2 -count=1 ./internal/gate` is green with immutable template state after once completion.
- [ ] [SFB15] `GOMAXPROCS=2 go test -p 1 -parallel 2 -count=1 ./internal/gate` is green without width-dependent coordination or timeout.
- [ ] [SFB16] `BENCH_KIT=/nonexistent GOMAXPROCS=2 go test -p 1 -parallel 2 -count=1 ./internal/gate` is green without ambient template-source or publication lookup.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| SFB1/lazy-template-build | restore a real `go build -o <root>/dist/bench.staged` inside every constructor seal | the focused filesystem-output measurement | measure the shared baseline, apply the per-root build, rerun the SFB13 command, and expect the output limit to red |
| SFB2/error-returning-primitive | make the primitive terminate through `testing.T` instead of returning its build error | the sticky-error subprocess witness | apply, force template construction to fail twice, and expect the second observable error assertion to become unreachable |
| SFB2/testing-wrapper | bypass the primitive in the existing testing wrapper | the exact build-command source audit | apply, run the audit, and expect two independent `go build -buildvcs=false` constructors to be named |
| SFB3/component-input-reseal | route the edited `cmd/bench/main.go` re-seal through the immutable template | the existing component-input source-change test | apply, run the focused test, and expect the source identity or seal assertion to red |
| SFB3/component-identity-reseal | route a component-identity re-seal through the immutable template | the existing component-identity source-change test | apply, run the focused test, and expect the changed identity assertion to red |
| SFB3/check-slot-reseal | route a check-slot re-seal through the immutable template | the existing check-slot source-change test | apply, run the focused test, and expect the slot reuse or retirement assertion to red |
| SFB4/build-skip-plant | route a build-skip planted binary through the immutable template | the existing planted-build-skip test | apply, run the focused test, and expect the planted digest assertion to red |
| SFB4/attestation-plant | route the attestation plant through the immutable template | `TestPlantedSealFailsAttestation` | apply, run the test, and expect the planted-binary distinction to red |
| SFB4/alternate-artifact | route the alternate package build through the immutable template | the existing alternate-artifact attestation test | apply, run the test, and expect the alternate executable digest assertion to red |
| SFB4/authorship | route the changed-source green build through the immutable template | `TestGreenBuildAttestsItsOwnBinary` | apply, run the test, and expect the authored digest assertion to red |
| SFB5/per-root-copy | pass the process-scoped staged path directly to freshness publication | a sequential two-fixture construction witness | apply, construct two fixtures, and expect the second construction to fail because the first publish consumed the template |
| SFB6/publication-consumption | move template bytes directly into one fixture's final executable without the existing publisher | `TestKitShapedFixtureBinaryIsSealed` | apply, run the test, and expect freshness verification to red |
| SFB7/root-isolation | replace the ordinary copy with a hardlink, then mutate the first published fixture binary | a sequential two-fixture isolation witness | apply, truncate the first fixture's executable, and expect the second fixture's freshness verification to red |
| SFB8/sticky-error | discard the once's construction error after a forced build failure | `TestKitShapedFixtureTemplateFailureIsSticky` error identity | apply, force template construction to fail, call the accessor twice, and expect distinct error values or an unusable empty path |
| SFB9/process-cleanup | remove package teardown or point it away from the process template | `TestKitShapedFixtureTemplateFailureIsSticky` outer-process lifetime assertion | apply, run the subprocess witness, and expect its reported template directory to remain after the child exits |
| SFB13/focused-output | restore per-root linking | the exact paired SFB13 filesystem-output receipts | run the shared command, restore per-root linking, rerun the same command, and expect shared output to exceed 60 percent of the mutation baseline |
| SFB14/race | write the process-scoped template path on every accessor call after once completion | the race detector | apply the unsynchronized write, run the exact two-core race command, and expect a reported read/write race under parallel fixture construction |
| SFB15/two-core | recursively enter the same once while its constructor is active | the bounded two-core package run | apply, run the exact SFB15 command with a 120-second timeout, and expect the once recursion to time out instead of completing |
| SFB16/hostile-kit | resolve the template build source through `kitRoot(root)` | the hostile-ambient package run | apply, run the exact SFB16 command, and expect fixture construction or identity assertions to red |
