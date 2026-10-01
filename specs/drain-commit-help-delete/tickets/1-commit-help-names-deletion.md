# Name the deletion route in bench commit help

Blocked by: none
Writes: internal/commit/commit.go, internal/commit/landing_test.go
Covers: none

## What to build

A learning of 2026-09-30 records that a delegated author ran a raw `git rm` to delete a tracked file. No help text tells an author how a ticket deletes a tracked file. But `bench commit` already commits a named path that is absent from the worktree as a deletion. The composition step starts a prospective index from the expected base. It stages each named path with `git add -A`, so the index drops a removed file or folder.

This ticket adds one line to the `bench commit --help` text. The line states that a deleted named path commits as a deletion. This applies to a file and to a folder, so no raw `git rm` is necessary. The help text in `internal/commit` stays the one source of this fact.

No current test proves the behavior through the verb. One test in `internal/commit/landing_test.go` removes a tracked file and a tracked folder and names both to the command. The test asserts that the published commit tracks neither path and that the checkout is clean. A second test asserts that the help text has the deletion line.

## Acceptance

- [ ] `bench commit --help` prints one line that names the deletion route and says that no `git rm` is necessary.
- [ ] The help test fails before the help edit and passes after it.
- [ ] The deletion test passes, and a `bench probe` mutation that stops the composition step from staging removals turns it red. The probe restores the file.
- [ ] `go test ./internal/commit/ ./cmd/bench/` passes.
