# Share the probe clamp golden

Blocked by: none
Writes: internal/probe/golden_test.go (new), internal/probe/probe_test.go, internal/probe/probetest/fixture.go, internal/probe/refusal_test.go, internal/probe/verdict_test.go, cmd/bench/record_operations_test.go, internal/reviewrecord/recordcmd/verification_import_test.go, .bench/commitment.json
Covers: LGP1, LGP2, LGP3, LGP4, LGP5, LGP37, LGP34

## What to build

The probe test package builds the clamp repository and its sealed run binary once for each package run. Each test then takes a private copy.

A `TestMain` owns one temporary parent directory for the golden and removes it after the run, as `internal/adopt/golden_test.go:39-69` does. A golden built inside one test would vanish when that test ends, because `probetest.Repository` takes its directories from `t.TempDir()` (`internal/probe/probetest/fixture.go:75-79`).

`clampFixture` (`internal/probe/probe_test.go:27-36`) copies the golden into private directories for each test. The repository root with its `.git`, the executable, the seal, and the manifest directory move together, because the seal holds digests and no path (`internal/freshness/freshness_publish.go:97-111`). Each test keeps its own environment bindings and its `t.Chdir`.

`probetest.Repository` keeps its signature, because `cmd/bench/record_operations_test.go:131` and `internal/reviewrecord/recordcmd/verification_import_test.go:152` call it. Split its build step from its per-test binding inside `probetest`.

A counting publisher wraps `freshness.Publish`, so that a test observes the number of golden builds. Put the new tests in `internal/probe/golden_test.go`.

`clampFixture` has 101 call sites in 17 files, and `refusalFixture` has 48 in 9 files. Keep both signatures. Run the whole package to prove that no caller depends on a fresh `git init` state.

Every witness keeps its assertions. `TestProbeBitesWhenTheFocusedTestFails` (PB1) keeps its real `go test` child. Remove no test.

Record `bench test --full --package ./internal/probe` before and after the change, for LGP34.

The ticket names `.bench/commitment.json` only because the spec-stage close commits the deliverable binding inside the graded range. The author does not edit that file.

## Acceptance

- [ ] Two `clampFixture` calls in one package run invoke the fixture publisher once (LGP1).
- [ ] A write in one copy leaves a second copy's subject bytes unchanged (LGP2).
- [ ] A copy passes the run-binary freshness check for its own root (LGP3).
- [ ] PB1, PB22, and PS40 to PS43 pass with their assertions unchanged (LGP4, LGP5, LGP37).
- [ ] The probes "build a fresh repository in clampFixture" and "return the golden root instead of a private copy" each bite.
- [ ] The verification record holds the package time before and after the change (LGP34).
