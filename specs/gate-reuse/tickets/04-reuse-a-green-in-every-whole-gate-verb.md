# Reuse another run's green in every whole-gate verb

Blocked by: 02-name-evidence-by-the-grade-key.md, 03-record-a-reused-green.md
Writes: internal/worktree/land_reuse_test.go (new), internal/worktree/land_fixtures_test.go, internal/gate/verb_reuse_test.go (new), internal/gate/reader_set_test.go, internal/commit/commit_reuse_test.go (new), internal/worktree/parallel_census_test.go
Covers: GR25, GR26, GR27, GR28, GR66, GR58, GR48, GR49, GR50, GR51, GR52

## What to build

Each whole-gate verb reaches the store through `executeSubjectOwning`: `bench gate` in each form directly, and the landing, a merge to a no-lane target, and a no-lane commit through `authorization.AuthorizeWithWriters` and `gate.ExecuteTree`.
The reuse checks at its two sites run before the delta plan and the reader plan, so the lookup already precedes every reduced schedule.
This ticket proves each verb on the grade key with real fixtures that count the whole-gate runs. A count that shows a second run is a plan expansion, which the operating guide owns.

A landing that reuses a green publishes as a green landing and advances the green marker.
`--fresh` skips the lookup. A commit or a merge with a declared lane runs its lane and reads no store green.

The running executable seam of ticket 01 lives in `internal/gate`, so the row for another running executable drives the complete checkpoint and `gate.ExecuteTree` from `internal/gate/verb_reuse_test.go`.
The fixtures in `internal/worktree` and `internal/commit` run every verb in one test process, so their running executables stay equal.

## Acceptance

- [ ] In a fixture with no lane, a spec-less landing of a merge tip that the merge graded green adds no whole-gate run and advances the green marker.
- [ ] In a fixture with no lane, a merge whose composed tree another checkout graded green adds no whole-gate run.
- [ ] A spec landing after a green complete checkpoint in the source worktree adds no whole-gate run, also when the source changes a Bench build input.
- [ ] A prospective run after a complete checkpoint under another running executable runs the whole gate.
- [ ] A plain `bench gate` in the primary checkout after a green prospective run of the same tree in a linked worktree prints the reuse line.
- [ ] In a repository with no lane, a commit of a composition that another checkout graded green adds no whole-gate run.
- [ ] A commit with a declared lane runs its lane when a whole green of the composed tree exists.
- [ ] A kit landing that is eligible for the reader set reuses a whole green of its tree and writes no reader record.
