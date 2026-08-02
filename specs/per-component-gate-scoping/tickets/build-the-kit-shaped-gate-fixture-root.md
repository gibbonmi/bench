# Build the kit-shaped gate fixture root

Blocked by: none
Ownership fence: `internal/gate/kitshaped_fixture_test.go`
Assumptions: `reducedRunFixture` (in `reduced_run_test.go`) is the existing
synthetic root — two phases, no Go module, no build phase — and stays untouched;
`writeGateTestFile`, `gitRun`, and `gitOutput` are the package's existing test
helpers. Re-derive from the tree at pickup.

## What to build

The existing fixture cannot exercise build, seal, or canary scoping, so every
later ticket in this build needs a second one. Add a kit-shaped temp root
alongside it: a real Go module with a `./cmd/bench` main, at least one package
outside that binary's closure whose test files belong to the module-wide
`go list -deps -test ./...` closure, `scripts/go-build.sh` with its
`scripts/go-build.inputs` manifest, `bin/bench.sh`, an `internal/canary/`
package, a `tests/canary/` fixture directory, and a phase manifest whose table
carries build, the toolchain phases, conformance, contract, shellcheck, and
canary. Each phase appends its own name to `.git/phase-runs`, exactly as the
existing fixture does, so executed-set assertions read a durable marker rather
than a return value. The root claims kit identity through `BENCH_KIT` and seeds
a published `dist/bench` with a valid seal via `freshness.Publish`.

This ticket lands the fixture and its self-tests only. It changes no production
behavior: the fixture gates green today, under the whole-tree path.

## Acceptance

- [ ] [PS4] the fixture gates green end to end and every phase in its resolved table leaves exactly one marker line.
- [ ] [PS5] the fixture's resolved phase table carries the build phase and the canary phase by name, asserted before any test relies on either.
- [ ] [PS6] the fixture carries a package outside `./cmd/bench`'s build closure whose `_test.go` files are inside the module-wide test closure, asserted by comparing the two `go list` closures.
- [ ] [PS7] the fixture's `dist/bench` passes `freshness.Verify` immediately after construction.
- [ ] [PS8] a second gate run over an unedited fixture reuses the whole-tree green and leaves no new phase markers.

## Red mutations

| criterion | mutation | owner | operation sequence |
|---|---|---|---|
| PS4 | drop `-count=1`-equivalent marker append from one phase script | `TestKitShapedFixtureGatesGreen` | construct the fixture, execute the gate, read `.git/phase-runs`, compare against the resolved table |
| PS5 | remove `scripts/go-build.sh` from the fixture so `BenchkitPhases` materializes no build phase | `TestKitShapedFixtureCarriesBuildAndCanary` | construct the fixture, resolve its phase table, assert both names present |
| PS6 | move the outside-closure package under `cmd/bench`'s import graph | `TestKitShapedFixtureHasAPackageOutsideTheBinaryClosure` | construct the fixture, run both `go list` closures, assert the module-wide set strictly contains a package the binary set omits |
| PS7 | seed `dist/bench` by copying bytes without calling `Publish` | `TestKitShapedFixtureBinaryIsSealed` | construct the fixture, call `freshness.Verify(root, dist/bench)`, expect nil |
| PS8 | make a phase script mutate the tree during the run | `TestKitShapedFixtureReusesItsGreen` | execute twice with no edit, assert the second run adds no marker line |
