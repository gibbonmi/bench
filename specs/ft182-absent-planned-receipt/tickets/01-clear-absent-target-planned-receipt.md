# Clear a Planned-phase receipt over an absent target on resume

Blocked by: none
Writes: internal/worktree/resume.go, internal/worktree/resume_test.go, ROADMAP.md, roadmap/FT182.md (deleted)

## What to build

A cleanup that crashed between its receipt write and its first checkpoint
leaves an in-flight receipt at phase `planned`. When the target path is absent
by the time the retry runs, the retry reaches `finishInterruptedExplicit` in
`internal/worktree/resume.go`. That function accepts only the `removing`,
`removed`, `branch-removed`, and `terminal` phases, and it returns
`errStaleFingerprint` for `planned`. The retry then wedges: every later attempt
reads the same receipt and returns the same error.

Make the resume path treat a `planned` receipt over an absent target as a
cleanup that never started. The seam is the `found && ReceiptInFlight` branch
in `resume.go` before `interruptedCleanupIsPastReplanning` runs. When the
receipt's phase is `planned` and `ClassifyPathShape(target)` is `ShapeAbsent`,
the branch supersedes the receipt. It treats the receipt as not found, so the
retry does not enter `finishInterruptedExplicit`. The ordinary absent-target
plan then runs. Its fresh receipt overwrites the superseded one through
`intent.PutCleanupReceipt`, and the terminal receipt records the outcome.

Do not add `planned` to the phase switch in `finishInterruptedExplicit`. That
function verifies an existing recovery ref and never creates one. Do not add a
new lifecycle step.

A `planned` phase means no side effect ran. An absent target has nothing left
to preserve, so a superseded plan that would have preserved loses nothing.
A `planned` receipt over a target that still exists keeps its current behavior.
The `preserved` phase keeps its current behavior too, because a preserved
target has already had a side effect.

The focused test selector below matches about 55 test functions in the
package; that breadth is expected.

The landing retires the FT182 row: the index line in `ROADMAP.md` and the
detail file `roadmap/FT182.md` leave together. `ROADMAP.md` names FT182 only
on its index line and in no `## Dependencies` edge, so no other row changes.

## Acceptance

- [ ] A resume over an absent target with an in-flight `planned` receipt completes and returns no error. The receipt on disk is then `complete` at phase `terminal`.
- [ ] The same resume, run a second time, returns the completed plan from the receipt with no error.
- [ ] A resume over an absent target with an in-flight `planned` receipt of a preserving plan completes the same way. The superseded receipt is gone from disk.
- [ ] A resume over an absent target with an in-flight receipt at phase `preserved` keeps its current result.
- [ ] A resume over a present target with an in-flight `planned` receipt whose fingerprint does not match keeps returning `cleanup fingerprint is stale`.
- [ ] The new test in `internal/worktree/resume_test.go` reds on the tree before the fix with `cleanup fingerprint is stale`, and the delegate records that red.
- [ ] `go test ./internal/worktree/ -run 'Resume|Cleanup|Interrupted' -parallel 2` passes.
- [ ] `ROADMAP.md` carries no FT182 line and `roadmap/FT182.md` is absent.
