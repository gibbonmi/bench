# Carry and drop the census across a landing

Blocked by: 05-show-the-census-row-on-the-board.md
Writes: internal/census/census.go, internal/census/census_test.go, internal/worktree/land.go, internal/worktree/land_resume.go, internal/worktree/land_refusal.go, internal/worktree/lifecycle.go, internal/worktree/ownership.go, internal/worktree/resume.go, internal/worktree/clean_landed.go, internal/worktree/worktree.go, internal/worktree/land_journey_test.go, internal/worktree/land_resume_test.go, internal/worktree/worktree_test.go

## What to build

The landing carries the final count out, and one function drops the records
when the assignment retires.

This contract crosses into ticket 07:

- `census.Drop(home, root, assignment string) error` removes one assignment's file.

`Drop` is the one owner. An absent file is not an error.

Both landed forms gain `census=<n>` as the last key of `landed{...}`.
`landedComplete` and `landedIncomplete` render it, so a complete landing and an
incomplete landing both state the count. The landing reads the count before its
release step, because that step drops the records. An incomplete landing keeps
the file, so a resume reads the same count and prints it again. An assignment
with no file prints `census=0`. A landing that refuses before its gate prints no
landed record and drops nothing.

The drop belongs to the one shared retirement path in
`internal/worktree/lifecycle.go`, which `executeCleanup` owns. `bench worktree
release`, `bench worktree clean`, and the landing's own release step all reach
that path, so no second copy exists.

The landing already holds the home. `ReleaseCommand` and `CleanCommand`
receive a home and discard it today. Thread it through `releaseAssignment` in
`ownership.go` and `applyCleanupTransaction` in `resume.go` to the drop.
The `clean --landed` sweep in `clean_landed.go` is the fourth caller. Do
not read the environment below the command boundary; pass the value. Call
`t.Parallel()` in each eligible new test, because the package census grades
that call.

## Acceptance

- [ ] A landing of an assignment with three records prints `census=3` as the last key of `landed{...}`. (EC20)
- [ ] A landing of an assignment with no file prints `census=0`. (EC21)
- [ ] An incomplete landing keeps the file, and its resume prints the same count again. (EC22)
- [ ] `bench worktree release` of an assignment with records leaves no file for it. (EC23)
- [ ] `bench worktree clean` of an assignment and a complete landing each leave no file. (EC24)
- [ ] A landing that refuses before its gate prints no landed record and keeps the file.
- [ ] One `census.Drop` call site exists in `internal/worktree`, on the shared retirement path.
