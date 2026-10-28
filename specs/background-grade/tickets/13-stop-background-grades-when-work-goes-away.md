# Stop background grades when work goes away

Blocked by: 09-queue-and-supersede-background-grades.md
Writes: internal/worktree/worktree.go, internal/worktree/ownership.go, internal/worktree/resume.go, internal/worktree/clean.go, internal/worktree/clean_set.go, internal/worktree/pool_reclaim.go, internal/worktree/land_stack.go, internal/worktree/land_effects.go, internal/worktree/background_stop.go (new), internal/worktree/background_stop_test.go (new), internal/worktree/parallel_census_test.go, internal/gocache/clean.go, internal/gocache/background_clean_test.go (new), cmd/bench/main.go, cmd/bench/help_inventory_test.go, internal/gate/backgroundstate (new), CHANGELOG.md, internal/gate/background_launch.go, internal/gate/background_launch_test.go, internal/gate/background_grade.go, internal/gate/background_outcome_test.go, internal/gate/background_queue_test.go, internal/status/background.go, internal/status/status_background_test.go, internal/gate/complete_checkpoint_test.go, .agents/commands/bench-implement-spec.md, .bench/BENCH-reference.md, internal/gate/checkout_guard.go, internal/gate/checkout_guard_test.go
Covers: BG78, BG79, BG80, BG81, BG82, BG117

## What to build

Stop the grades of each worktree that a verb removes or lands as a stack member, before the verb acts.
Each stop goes through the stop operation of ticket 05, so the verb waits until the grade has exited.

- `bench worktree release` stops the running grade and drops the queued grade of its worktree before it removes the worktree.
- `bench worktree clean` and `bench worktree reclaim` do the same for each worktree that they remove.
- A landing stack does the same for each member worktree before it grades.
- `bench cache clean` stops every grade and drops every queued grade, then takes its lock as it does today.

Release dispatches through `worktree.go:424` to `ownership.go:216`, removal runs through `resume.go:104`, and clean selects its set in `clean_set.go:359`.
The author picks the one hook site that covers each verb.
The cache clean help row ends with "refuses under a live run". Keep or change it, and keep its pin in `TestHelpInventoryIsComplete` true.
`cmd/bench/main.go` has 415 lines, which is over its budget, so it keeps its line count.
Each new test in `internal/worktree` runs in parallel, or it raises `worktreeSerialCeiling` in this ticket.

## Acceptance

- [ ] `bench worktree release` stops the grade of its worktree before it removes the worktree.
- [ ] `bench worktree clean` stops the grade of each worktree that it removes.
- [ ] `bench worktree reclaim` stops the grade of each worktree that it removes.
- [ ] A landing stack stops the grade of each member worktree before it grades.
- [ ] `bench cache clean` stops every grade and drops every queued grade, then cleans.
- [ ] `bench structure --growth` reports no growth of `cmd/bench/main.go`.
