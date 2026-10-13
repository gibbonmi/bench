# Pin both clocks in the handoff route test

Blocked by: none
Writes: internal/handoff/render_test.go
Covers: none

## What to build

Pin the commit timestamp and the modification time of the handoff file in `TestCommandRoutesOnlyAnEmptyNextCommand`.
Set `GIT_AUTHOR_DATE` and `GIT_COMMITTER_DATE` before `benchRepo` creates the base commit.
Set the modification time of the handoff file with `os.Chtimes` to a time after the base commit.

## Acceptance

- [ ] `TestCommandRoutesOnlyAnEmptyNextCommand` pins the base commit timestamp and the modification time of the handoff file, and routes to `bench roadmap`.
