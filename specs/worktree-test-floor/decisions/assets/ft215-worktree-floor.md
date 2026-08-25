# The `internal/worktree` test floor

Measured on 2026-08-25 in the `ft215-shape` worktree with
`go test -count=1 -json`. This asset supports `specs/worktree-test-floor/decisions/worktree-test-floor.md`.

## Findings

- The package ran 72.9 s alone and 87.1 s inside the full `./...` run. The sum of
  the top-level test times equals the wall, so the suite is fully serial. No test
  file in `internal/worktree` calls `t.Parallel()`.
- The cost is count, not one slow test: 334 top-level tests (611 with subtests)
  at about 218 ms each. Every fixture repository comes from
  `internal/gittest.RepoOnBranch` (`internal/gittest/gittest.go:52-58`), which
  spawns `git init` and two `git config` calls. Git operations inside tests spawn
  through `descendant()` (`internal/worktree/journey_test.go:104-114`). There are
  239 `gitRun(` call sites and 132 repository-creation call sites.
- The five slowest tests: `TestResumeLandCommandPublicRefusesDestructiveDestinationState`
  4.97 s (`land_resume_refusal_test.go:16-52`),
  `TestResumeLandCommandRefusesWhenTerminalReceiptWasEvicted` 2.80 s,
  `TestLifecycleFaultBoundariesRemainLockedOrAbsent` 2.61 s (`ownership_test.go:233`),
  `TestLandCommandPublicRealGitJourney` 2.28 s (`land_journey_test.go:16`),
  `TestRecoveryPreservesEveryGitVisibleLayerWithoutMovingBranchOrIndex` 2.17 s
  (`lifecycle_acquire_test.go:121`).
- No `time.Sleep` exists in the package. The floor is fork/exec and disk I/O.
- The package already builds one executable per run through a `sync.Once`
  selector (`internal/worktree/test_run_test.go:40-45,78-80`). FT246 names this
  package as the origin of that pattern, so a second-build cut buys nothing here.
- The next packages: `internal/freshness` 14.3 s and `internal/runbinary` 13.5 s,
  both on FT246's list. If `internal/worktree` fell to 20 s, the floor would move
  to `internal/freshness`.

## The priced cut

Run the fixture subtests in parallel with `t.Parallel()`, one file or table at a
time, from the five slowest tests down. Each subtest already owns its own
`t.TempDir()`. On eight or more cores, a fork/exec workload of 72.9 s serial work
lands at 15–25 s wall.

The cost is a harness refactor, not the flag. `bindEnv` and `chdir`
(`internal/worktree/journey_test.go:88-99`) call `t.Setenv` and `t.Chdir`, and Go
panics when a parallel test calls either. The journeys must receive explicit paths
and environment instead, or the tests that need them must stay serial.
`TestSerialJourneyHarnessCensus` enforces the serial harness shape today, so the
census must learn the new shape.
