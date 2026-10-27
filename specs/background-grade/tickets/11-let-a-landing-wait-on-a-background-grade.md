# Let a landing wait on a background grade

Blocked by: 09-let-a-chunk-checkpoint-wait-on-a-background-grade.md, 10-grade-a-landing-subject-in-the-background.md
Writes: internal/worktree/land.go, internal/worktree/land_background.go (new), internal/worktree/land_background_test.go (new), internal/worktree/parallel_census_test.go, internal/gate/background_wait.go (new), internal/gate/backgroundstate (new), CHANGELOG.md
Covers: BG2, BG3, BG75, BG76, BG77

## What to build

Today a landing takes the destination lease before it grades (`internal/worktree/land.go:201-204`).

Before the lease, a landing looks for a live claim of its own worktree's landing grade.
If it finds one, it waits until the claim ends, with the wait rules of ticket 09.
Under the lease, it computes its key and looks it up, as it does today. So the landing composes no extra tree before the lease.

A landing whose key differs from the key of the running grade of its worktree stops that grade and drops its queued grade before it grades.

`land.go` has 8 lines of headroom, so put the new code in `land_background.go`.
Each new test in `internal/worktree` runs in parallel, or it raises `worktreeSerialCeiling` in this ticket.

## Acceptance

- [ ] The landing of a spec-less kit worktree after its background green runs no oracle.
- [ ] The landing of a tickets-only kit worktree after its background green runs no oracle.
- [ ] A landing that finds a live claim of its worktree's landing grade waits before it takes the destination lease.
- [ ] After the wait, the landing takes its lease and reuses the green with no oracle run.
- [ ] A landing of another key stops the grade of its worktree and drops its queued grade before it grades.
