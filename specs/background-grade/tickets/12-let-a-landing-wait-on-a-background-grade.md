# Let a landing wait on a background grade

Blocked by: 10-let-a-chunk-checkpoint-wait-on-a-background-grade.md, 11-grade-a-landing-subject-in-the-background.md
Writes: internal/worktree/land.go, internal/worktree/land_background.go (new), internal/worktree/land_background_test.go (new), internal/worktree/parallel_census_test.go, internal/gate/background_wait.go (new), internal/gate/background_landing.go (new), internal/gate/backgroundstate (new), CHANGELOG.md
Covers: BG2, BG3, BG75, BG76, BG77, BG110, BG111

## What to build

Today a landing takes the destination lease before it grades (`internal/worktree/land.go:201-204`).

Before the lease, a landing computes its grade key: the published tree composed on the current destination, under the broker.
If it finds a live claim on that key, it waits until the claim ends, with the wait rules of ticket 10.
Under the lease, it computes its key again and looks it up, as it does today. A destination that moved during the wait gives a new key, so the landing runs.

A landing whose key differs from the key of the running grade of its worktree does not wait.
It stops that grade, waits for its exit, and drops the queued grade of its worktree before it grades.
A queued grade of the landing's own key leaves the queue when the landing starts its own run.

`land.go` has 8 lines of headroom, so put the new code in `land_background.go`.
Each new test in `internal/worktree` runs in parallel, or it raises `worktreeSerialCeiling` in this ticket.

## Acceptance

- [ ] The landing of a spec-less kit worktree after its background green runs no oracle.
- [ ] The landing of a tickets-only kit worktree after its background green runs no oracle.
- [ ] A landing that finds a live claim on its own key waits before it takes the destination lease.
- [ ] After the wait, the landing takes its lease and reuses the green with no oracle run.
- [ ] A landing whose key differs from the running grade of its worktree does not wait.
- [ ] A landing of another key stops the grade of its worktree and drops its queued grade before it grades.
- [ ] A queued grade of the landing's key leaves the queue when the landing starts its own run.
