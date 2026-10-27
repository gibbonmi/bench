# Start a background grade after a lane pass

Blocked by: 06-grade-a-chunk-subject-in-the-background.md
Writes: internal/commit/commit.go, internal/commit/background.go (new), internal/commit/background_test.go (new), internal/worktree/merge.go, internal/worktree/merge_background_test.go (new), internal/worktree/parallel_census_test.go, internal/gate/background_launch.go (new), internal/gate/background_launch_test.go (new), internal/gate/backgroundstate (new), internal/subprocess, internal/conformance/process_lifetime_test.go, CHANGELOG.md
Covers: BG5, BG6, BG7, BG8, BG9, BG10, BG11, BG43, BG44, BG45

## What to build

After `bench commit` publishes a commit, and after `bench worktree merge` publishes a merge, start a background grade of the new tip.
Start one only in a Bench assignment worktree for which `KitForRoot` answers the worktree itself (`internal/gate/kit_root.go:12-19`).
A linked repo, a dry run, a refused commit, and a merge of kind `current` start none.

The launch goes through the lifetime owner of `internal/subprocess`.
The grade runs as a detached session, outside the process group of the verb and of the harness, so it outlives `bench commit`.
The verb returns at once and does not wait for the grade.
Add the launch caller to `processLifetimeCensus` (`internal/conformance/process_lifetime_test.go:40`).

This ticket keeps one machine slot. While one grade runs, a second lane pass starts no second oracle. Ticket 08 adds the queue.

`commit.go` has 400 lines, at its budget. Put the new code in `internal/commit/background.go`.
Each new test in `internal/worktree` runs in parallel, or it raises `worktreeSerialCeiling` in this ticket.

## Acceptance

- [ ] A `bench commit` that publishes in a kit worktree starts one grade of its tip.
- [ ] A `bench worktree merge` of kind `merge` starts one grade of the new tip.
- [ ] A commit in a linked repo, a `--dry-run`, a refused commit, and a merge of kind `current` each start no grade.
- [ ] `bench commit` exits 0 while its grade still holds its claim.
- [ ] The grade runs in a session other than the session of the launching verb.
- [ ] `bench test --check process-lifetime-closure` passes with the new census row.
- [ ] A second lane pass while one grade runs starts no second oracle.
