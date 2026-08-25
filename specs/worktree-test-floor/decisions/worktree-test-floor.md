# The worktree test floor

Status: ready

## Destination

The `internal/worktree` test package no longer sets the gate's wall. The gate's
`test` phase runs packages in parallel, so its wall is the slowest package plus
the build. `internal/worktree` at 51–87 s is that package. The destination is
a package wall near 20 s, so the floor moves to the next package at about 14 s.
What the gate proves does not change. This is the second spec that
`specs/one-change-one-grade/decisions/one-change-one-grade.md` #6 splits out.

## #1: What sets the `internal/worktree` test floor?

Blocked by: none
Type: Research

### Question

Profile the package with per-test timing. Name the top costs (real Git journeys,
executable builds, waits) and price one cut. FT246 is adjacent.

### Answer

Resolved 2026-08-25: count, not one slow test. The suite is fully serial: 334
top-level tests at about 218 ms each, every one a chain of `git` spawns and, for
the journeys, one real `bench` child. The executable is already built once per
run, so FT246's cut does not apply here. The priced cut is `t.Parallel()` on the
fixture subtests, which needs the journey harness to stop calling `t.Setenv` and
`t.Chdir`. See `specs/worktree-test-floor/decisions/assets/ft215-worktree-floor.md`.

## Not yet specified

## Spec-writer discretion

- Which subtests run in parallel and which stay serial. The choice is reversible
  and changes nothing the gate proves.
- How the journey harness passes paths and environment to the spawned binary
  instead of `t.Setenv` and `t.Chdir`.
- The order of files in the refactor, and whether `TestSerialJourneyHarnessCensus`
  learns the new shape or a successor census replaces it.

## Out of scope

- Removing, merging, or weakening any test to buy wall-clock.
- A second executable-build cut in this package. It already selects one
  executable per run.
- The `internal/freshness` and `internal/runbinary` floors. FT246 owns them.
- Any change to the gate's phase table or to what green means.

## Sources

- Path: `specs/worktree-test-floor/decisions/assets/ft215-worktree-floor.md`
  Supports: #1's answer, the destination's numbers, and the discretion items. Produced 2026-08-25 by one read-only research delegate from a `go test -count=1 -json` run in the shaping worktree.
  Drift: re-measure if `internal/worktree`'s test count or `internal/gittest.RepoOnBranch` changes before the spec reads this map.
